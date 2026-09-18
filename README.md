# AKS Health Signal

Kubernetes Custom Resource Definitions (CRDs) that let a monitoring app gate AKS
upgrades on its own health verdict.

## How gating works

1. A monitoring app installs an **UpgradeGatePolicy** declaring which node pools
   or the whole cluster it gates, and how long AKS should wait for its verdict.
   Installing this policy is what arms the gate.
2. An upgrade starts. AKS creates an **UpgradeOperation** describing what is
   being upgraded.
3. AKS creates **HealthCheckRequest** objects for the targets: 
   one per node being upgraded, 
   plus one for the node pool or cluster currently being upgraded.
4. The health signal provider which registered in **UpgradeGatePolicy** watches 
   for HealthCheckRequests and, for each one it handles, creates a **HealthSignal** 
   that owns its verdict, labelled with `health.aks.io/provider` 
   so AKS knows which policy it satisfies.
5. The app reports the verdict in `status.conditions`.
6. AKS proceeds on `"True"`, aborts on `"False"`, and keeps waiting on
   `"Unknown"` until the policy's `nodeTimeout` elapses — at which point it aborts.

Steps 3–6 repeat for each node in turn.

Two rules follow from this, and both matter:

- **Registration is installation.** A target no installed policy matches is not gated.
- **Every HealthSignal must carry `health.aks.io/provider` to work with upgrade gate.** 
  AKS matches that label against `UpgradeGatePolicy.spec.provider`. 

## Resources

| Resource | Created by | Purpose |
|----------|-----------|---------|
| `UpgradeGatePolicy` | Monitoring app (install time) | Registers that this app gates some targets, and on what terms |
| `UpgradeOperation` | AKS | Describes an in-progress upgrade |
| `HealthCheckRequest` | AKS | Asks for a health verdict on a node, node pool, or cluster |
| `HealthSignal` | Monitoring app | Carries the verdict |

All four are cluster-scoped.

---

## UpgradeGatePolicy

Durable, install-time configuration — typically shipped in the app's Helm chart,
the same way an application ships its own PodDisruptionBudget. Created once and
survives every upgrade. AKS only ever reads it.

**One object per provider.** An app gating several pools on different terms lists
several rules in one object. Two different apps use two separate objects, so
uninstalling one cannot disturb the other.

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeGatePolicy
metadata:
  name: node-monitor
spec:
  # Must match the health.aks.io/provider label this app sets on its HealthSignals.
  provider: node-monitor.example.com
  rules:
    # Gate userpool, waiting up to 5 minutes after each node for a verdict.
    - target:
        scope: NodePool
        name: userpool
      nodeTimeout: 5m
      onFailure: Abort
    # The same provider may gate another pool on different terms.
    - target:
        scope: NodePool
        name: agentpool
      nodeTimeout: 2m
      onFailure: Abort
```

Gating the whole cluster — every node, including pools added later:

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeGatePolicy
metadata:
  name: cluster-monitor
spec:
  provider: cluster-monitor.example.com
  rules:
    # Cluster scope gates every node; name is omitted.
    - target:
        scope: Cluster
      nodeTimeout: 1m
      onFailure: Abort
```

### Matching semantics

- A target may be gated by several providers at once. It waits for **all** of
  them.
- **Within one policy, the most specific rule wins.** A `NodePool` rule overrides
  that provider's `Cluster` rule for nodes in that pool, so a policy can set a
  cluster-wide default and then make exceptions per pool.
- **Each provider's `nodeTimeout` is independent.** A target waits for every
  provider that gates it, each against its own deadline. If one provider has no
  verdict by its deadline, the upgrade fails then. A provider that has already
  reported healthy is done, and the others keep running against their own
  deadlines.
- `nodeTimeout` starts when the health check begins, **not** when the HealthSignal
  appears. A provider that never creates one is still bounded by it.
- **`nodeTimeout` is capped at 5 minutes in AKS RP.** Larger values are clamped, and
  omitting it means the cap. The ceiling is per node, not per operation, so a
  50-node pool can still take hours — it bounds how long any single node may
  stall the rollout.
- **`nodeTimeout` is a window, not just a node deadline.** It opens after each
  node is upgraded and closes early once that node's check succeeds. While it is
  open the RP watches this provider's signals at *every* scope:
  - a `"False"` from **any** scope — node, node pool or cluster — aborts at once,
    without waiting out the remaining time;
  - only the **node** verdict is awaited. A NodePool- or Cluster-scoped check that
    never reports does not block the upgrade, so a slow or absent cluster-level
    monitor cannot stall a rollout.

  So failure is immediate at any scope, while success is only ever required of
  the node. This is the same asymmetry the kubelet applies to a `startupProbe`:
  the budget is spent waiting for success, but a crash aborts instantly.
- `onFailure` currently accepts only `Abort`; A policy
  asking for anything else is rejected at admission rather than silently
  downgraded.
- `spec.provider` is compared against a label value, so it must be a valid one:
  at most 63 characters, alphanumeric at each end, `-`, `_` or `.` in between,
  and **no `/`**. Use reverse-DNS to avoid collisions between vendors.

A cluster-wide default with a per-pool exception:

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeGatePolicy
metadata:
  name: mixed-monitor
spec:
  provider: mixed-monitor.example.com
  rules:
    # Default for every node in the cluster.
    - target:
        scope: Cluster
      nodeTimeout: 2m
      onFailure: Abort
    # userpool is slower to settle, so it overrides the default above.
    - target:
        scope: NodePool
        name: userpool
      nodeTimeout: 5m
      onFailure: Abort
```

---

## UpgradeOperation

Created by AKS to represent one or more in-progress upgrade targets. Cluster and
node-pool targets may coexist in the same UpgradeOperation, but each target may
have at most one active UpgradeOperation at a time.

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeOperation
metadata:
  name: cluster-upgrade
  annotations:
    kubernetes.azure.com/upgradeOperationId: "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a"
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
spec:
  - type: Cluster              # Cluster | NodePool
    targetNames:
    - my-aks-cluster
  - type: NodePool
    targetNames:
    - systempool
    - userpool
```

---

## HealthCheckRequest

Created by AKS to ask for a verdict. **Scope determines what is being asked
about**, and an app must branch on it:

| `scope` | Created when | `spec.targetRef.name` | Use it to check |
|---------|--------------|------------------|-----------------|
| `Node` | Each node is replaced | the node name | that specific node |
| `NodePool` | A node pool is upgraded | the node pool name | workloads, services, network paths across the pool |
| `Cluster` | The cluster is upgraded | the cluster name | cluster-wide app, service and network health |

The NodePool and Cluster scopes exist so an app can report on **its own
workloads** — deployment replicas, service endpoints, network reachability —
rather than only on per-node health.

`spec.targetRef` is always set. `scope` says what kind of object the name refers
to, so an app must branch on `scope` before interpreting `targetRef.name`.

Node scope:

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthCheckRequest
metadata:
  name: aks-userpool-31207608-vmss000000-6e8ef28e
  annotations:
    kubernetes.azure.com/upgradeCorrelationID: "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a"
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
spec:
  scope: Node
  targetRef:
    name: aks-userpool-31207608-vmss000000
```

Node pool scope:

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthCheckRequest
metadata:
  name: userpool-6e8ef28e
  annotations:
    kubernetes.azure.com/upgradeCorrelationID: "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a"
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
spec:
  scope: NodePool
  targetRef:
    name: userpool
```

Cluster scope:

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthCheckRequest
metadata:
  name: upgrade-6e8ef28e
  annotations:
    kubernetes.azure.com/upgradeCorrelationID: "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a"
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
spec:
  scope: Cluster
  targetRef:
    name: my-aks-cluster
```

---

## HealthSignal

Created and updated **entirely by monitoring apps** in response to a
HealthCheckRequest. AKS only reads it.

Every HealthSignal must:

- carry the `health.aks.io/provider` label, matching its UpgradeGatePolicy;
- set an `ownerReference` to the HealthCheckRequest it answers, which links the
  two and lets Kubernetes garbage-collect the signal with the request.

`spec.type` mirrors the request's scope: `Node` → `NodeHealth`, `NodePool` →
`NodePoolHealth`, `Cluster` → `ClusterHealth`.

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthSignal
metadata:
  name: aks-userpool-31207608-vmss000000-nodemonitor
  labels:
    health.aks.io/provider: node-monitor.example.com
  ownerReferences:
  - apiVersion: health.aks.io/v1alpha1
    kind: HealthCheckRequest
    name: aks-userpool-31207608-vmss000000-6e8ef28e
    uid: 0f1b2c3d-4e5f-6071-8293-a4b5c6d7e8f9
spec:
  type: NodeHealth             # NodeHealth | NodePoolHealth | ClusterHealth
  targetRef:
    apiVersion: v1
    kind: Node
    name: aks-userpool-31207608-vmss000000
status:
  conditions:
  - type: Ready
    status: "True"             # True=Healthy, False=Unhealthy, Unknown=No verdict
    reason: Baseline
    message: Node health is healthy over last 5 minutes.
    lastTransitionTime: "2026-02-26T22:15:32Z"
```

Answering a node-pool-scoped request — the app judges its own workloads across
the pool:

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthSignal
metadata:
  name: userpool-6e8ef28e-nodemonitor
  labels:
    health.aks.io/provider: node-monitor.example.com
  ownerReferences:
  - apiVersion: health.aks.io/v1alpha1
    kind: HealthCheckRequest
    name: userpool-6e8ef28e
    uid: 3c4d5e6f-7a8b-9c0d-1e2f-30415263748a
spec:
  type: NodePoolHealth
status:
  conditions:
  - type: Ready
    status: "True"
    reason: Baseline
    message: All workloads in userpool are healthy.
    lastTransitionTime: "2026-02-26T22:15:32Z"
```

If **unhealthy**, AKS aborts the upgrade:

```yaml
status:
  conditions:
  - type: Ready
    status: "False"
    reason: NotReady
    message: Node was NotReady for >2 minutes.
    lastTransitionTime: "2026-02-26T22:16:10Z"
```

> The `ownerReferences.uid` in these samples is a placeholder. Kubernetes garbage
> collects any object whose owner uid does not resolve, so applying them as-is
> creates the signal and then immediately deletes it. Substitute the real uid
> (`kubectl get healthcheckrequest <name> -o jsonpath='{.metadata.uid}'`). A real
> app sets this from the request it just read.

### Condition semantics

| `status` | Meaning | AKS behaviour |
|----------|---------|--------------|
| `"True"` | Healthy | Continue upgrade |
| `"False"` | Unhealthy | **Abort** upgrade |
| `"Unknown"` | No verdict yet | Wait, until `nodeTimeout` elapses |

A node proceeds only when **every** provider gating it reports `"True"` for that
node. Any `"False"`, at any scope, aborts immediately. If `nodeTimeout` elapses
with no `"True"` from a node-scoped check, AKS treats the verdict as not healthy
and aborts. A NodePool- or Cluster-scoped check left at `"Unknown"` does not
abort — only an explicit `"False"` does.

There is **no timeout field on HealthSignal**. The deadline lives on
`UpgradeGatePolicy.spec.rules[].nodeTimeout` and is armed when the health check
starts, so it also bounds a provider that never creates a HealthSignal — the
failure a signal-side deadline could never catch. An app that wants to fail
sooner should report `"False"`.

---

## Well-known labels and annotations

| Key | On | Set by | Purpose |
|-----|----|--------|---------|
| `health.aks.io/provider` | HealthSignal | Monitoring app | **Required.** Identifies the provider; matched against `UpgradeGatePolicy.spec.provider` |
| `kubernetes.azure.com/upgradeOperationId` | UpgradeOperation | AKS | Uniquely identifies an UpgradeOperation |
| `kubernetes.azure.com/upgradeCorrelationID` | HealthCheckRequest | AKS | Links the request to an upgrade operation |
| `kubernetes.azure.com/targetKubernetesVersion` | UpgradeOperation, HealthCheckRequest | AKS | Target Kubernetes version |

List signals by provider:

```bash
kubectl get healthsignals -l health.aks.io/provider=node-monitor.example.com
```

## OwnerReferences & garbage collection

```
UpgradeOperation
└── HealthCheckRequest
    └── HealthSignal
```

Deleting a HealthCheckRequest cascades to the HealthSignals it owns, so signals
never outlive the request that prompted them.

## Development

### Generate CRDs from Go types

```bash
make manifests
```

### Generate DeepCopy methods

```bash
make generate
```

### Run both

```bash
make
```

### Install the CRDs

```bash
kubectl apply -f CRD/
```

### Project Structure

```
api/
  health/v1alpha1/     # HealthSignal & HealthCheckRequest Go types
  upgrade/v1alpha1/    # UpgradeOperation & UpgradeGatePolicy Go types
CRD/                   # Generated CRD YAML files
CR-samples/            # Example CR instances
monitoring/            # Monitoring deployment manifests (Datadog, keepalive)
```
