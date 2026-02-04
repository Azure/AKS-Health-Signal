# AKS Health Signal

Kubernetes Custom Resource Definitions (CRDs) for AKS health signaling and upgrade tracking.

## CRDs

### HealthSignal

Represents the health state of a cluster or node.

```yaml
apiVersion: health.aks.io/v1alpha1
kind: HealthSignal
metadata:
  name: aks-userpool-31207608-vmss000000
  ownerReferences:
  - apiVersion: upgrade.aks.io/v1alpha1
    kind: UpgradeNodeInProgress
    name: aks-userpool-31207608-vmss000000
spec:
  type: NodeHealth          # NodeHealth or ClusterHealth
  target:                   # Required for NodeHealth (corev1.ObjectReference)
    kind: Node
    name: aks-userpool-31207608-vmss000000
  source:                   # Optional: component that produced this signal
    kind: DaemonSet
    name: node-health-monitor
    namespace: kube-system
status:
  conditions:
  - type: Ready
    status: "True"          # True (Healthy), False (Unhealthy), Unknown
    reason: Baseline
    message: Node health is healthy.
    lastTransitionTime: "2026-01-29T10:00:00Z"
```

### UpgradeOperationInProgress

Signals that a cluster upgrade operation is in progress. The existence of this CR indicates an upgrade is active.

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeOperationInProgress
metadata:
  name: cluster-upgrade
  annotations:
    kubernetes.azure.com/upgradeOperationId: 6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a
    kubernetes.azure.com/targetKubernetesVersion: "1.33.5"
```

### UpgradeNodeInProgress

Signals that a specific node is currently undergoing upgrade.

```yaml
apiVersion: upgrade.aks.io/v1alpha1
kind: UpgradeNodeInProgress
metadata:
  name: aks-userpool-31207608-vmss000000
  ownerReferences:
  - apiVersion: upgrade.aks.io/v1alpha1
    kind: UpgradeOperationInProgress
    name: cluster-upgrade
  annotations:
    kubernetes.azure.com/agentpool: userpool
nodeRef:
  name: aks-userpool-31207608-vmss000000
```

## Ownership Hierarchy

```
UpgradeOperationInProgress (cluster-upgrade)
├── UpgradeNodeInProgress (per node)
│   └── HealthSignal/NodeHealth (per node)
└── HealthSignal/ClusterHealth
```

Deleting `UpgradeOperationInProgress` cascades to all child resources.

## Development

### Generate CRDs from Go types

```bash
make manifests
```

### Project Structure

```
api/
  health/v1alpha1/     # HealthSignal Go types
  upgrade/v1alpha1/    # UpgradeNodeInProgress, UpgradeOperationInProgress Go types
CRD/                   # Generated CRD YAML files
CR-samples/            # Example CR instances
```