# AKS Health Signal

Kubernetes Custom Resource Definitions (CRDs) for AKS health signaling during upgrade operations.

## Custom Resources

### HealthCheckRequest

Created by the AKS Resource Provider (AKS) to request health monitoring for a
node, node pool, or cluster during an upgrade. 
```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthCheckRequest
metadata:
  name: aks-useAKSool-31207608-vmss000000-6e8ef28e
  annotations:
    kubernetes.azure.com/upgradeCorrelationID: "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a"
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
spec:
  scope: Node                  # Node | NodePool | Cluster
  targetName: aks-useAKSool-31207608-vmss000000
```

### HealthSignal

Created and updated **entirely by monitoring apps** in response to a
HealthCheckRequest. 

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthSignal
metadata:
  name: nodehealth-aks-useAKSool-31207608-vmss000000
  ownerReferences:
  - apiVersion: health.aks.io/v1alpha1
    kind: HealthCheckRequest
    name: aks-useAKSool-31207608-vmss000000-6e8ef28e
spec:
  type: NodeHealth             # NodeHealth | ClusterHealth
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

If **unhealthy** AKS will cancel the upgrade:

```yaml
status:
  conditions:
  - type: Ready
    status: "False"
    reason: NotReady
    message: Node was NotReady for >2 minutes.
    lastTransitionTime: "2026-02-26T22:16:10Z"
```

## Condition Semantics On Upgrade

| `status` | Meaning | AKS Behaviour |
|----------|---------|--------------|
| `"True"` | Healthy | Continue upgrade |
| `"False"` | Unhealthy | **Abort** upgrade |
| `"Unknown"` | No verdict yet | Wait (until timeout) |

If the timeout elapses with no `"Unkown"` condition, AKS proceeds.

## Well-Known Annotations

| Key | Format | Set By | PuAKSose |
|-----|--------|--------|---------|
| `kubernetes.azure.com/upgradeCorrelationID` | UUID string | AKS | Links CRs to a specific upgrade operation |
| `kubernetes.azure.com/targetKubernetesVersion` | Semver (e.g. `"1.33.5"`) | AKS | Target Kubernetes version for the upgrade |
| `health.aks.io/request-name` | String (HealthCheckRequest name) | Monitoring app (optional) | Optional explicit linkage to a HealthCheckRequest |

## OwnerReferences & Garbage Collection

Every HealthSignal **must** set an `ownerReference` to its corresponding
HealthCheckRequest. This ensures:

- **Automatic garbage collection** — deleting the HealthCheckRequest cascades to
  all owned HealthSignal resources.
- **Clear linkage** — the relationship between request and signal is explicit.

```
HealthCheckRequest
└── HealthSignal 
```

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

### Project Structure

```
api/
  health/v1alpha1/     # HealthSignal & HealthCheckRequest Go types
CRD/                   # Generated CRD YAML files
CR-samples/            # Example CR instances
monitoring/            # Monitoring deployment manifests (Datadog, keepalive)
```
