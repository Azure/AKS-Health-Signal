package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// KindUpgradeOperation is the Kind string for UpgradeOperation resources.
	KindUpgradeOperation = "UpgradeOperation"

	// KindUpgradeGatePolicy is the Kind string for UpgradeGatePolicy resources.
	KindUpgradeGatePolicy = "UpgradeGatePolicy"

	// ResourceUpgradeOperations is the plural resource name for UpgradeOperation.
	ResourceUpgradeOperations = "upgradeoperations"

	// ResourceUpgradeGatePolicies is the plural resource name for UpgradeGatePolicy.
	ResourceUpgradeGatePolicies = "upgradegatepolicies"
)

var (
	// UpgradeOperationGVR is the GroupVersionResource for UpgradeOperation.
	UpgradeOperationGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: ResourceUpgradeOperations,
	}

	// UpgradeGatePolicyGVR is the GroupVersionResource for UpgradeGatePolicy.
	UpgradeGatePolicyGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: ResourceUpgradeGatePolicies,
	}
)
