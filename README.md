# AKS Health Signal

Kubernetes Custom Resource Definitions (CRDs) for AKS health signaling and upgrade tracking.

## CRDs

### HealthSignal

Represents the health state of a cluster or node.

```yaml
apiVersion: health.aks.io/v1
kind: HealthSignal
metadata:
  name: aks-userpool-31207608-vmss000000
spec:
  type: NodeHealth          # NodeHealth or ClusterHealth
  target:                   # Required for NodeHealth
    kind: Node
    name: aks-userpool-31207608-vmss000000
status:
  conditions:
  - type: Ready
    status: "True"          # True (Healthy), False (Unhealthy), Unknown
    reason: Baseline
    message: Node health is healthy.
    lastTransitionTime: "2026-01-29T10:00:00Z"
```

### UpgradeNodeInProgress

Tracks nodes currently undergoing upgrade operations.

```yaml
apiVersion: upgrade.aks.io/v1
kind: UpgradeNodeInProgress
metadata:
  name: aks-userpool-31207608-vmss000000
spec:
  agentPoolName: userpool
  nodeRef:
    kind: Node
    name: aks-userpool-31207608-vmss000000
status:
  conditions:
  - type: InProgress
    status: "True"
    reason: NodeImageUpgrade
    message: Upgrading node image
    lastTransitionTime: "2026-01-29T10:00:00Z"
```

## Development

### Generate CRDs from Go types

```bash
make manifests
```

### Project Structure

```
api/
  health/v1/       # HealthSignal Go types
  upgrade/v1/      # UpgradeNodeInProgress Go types
CRD/               # Generated CRD YAML files
CR-samples/        # Example CR instances
```