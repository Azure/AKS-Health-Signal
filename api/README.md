# Generate DeepCopy Methods for CRD

Use the following commands to generate the required DeepCopy methods for the custom resource:
```
# 1. Install controller-gen
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

# 2. Generate the DeepCopy methods
controller-gen object paths=./api/...
```

# Generate CRD YAML
```
controller-gen crd:crdVersions=v1 paths=./apis/... output:crd:stdout > ./CRD/upgrade.aks.io.yaml
```
