package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// KindUpgradeOperation is the Kind string for UpgradeOperation resources.
	KindUpgradeOperation = "UpgradeOperation"

	// ResourceUpgradeOperations is the plural resource name for UpgradeOperation.
	ResourceUpgradeOperations = "upgradeoperations"
)

var (
	// UpgradeOperationGVR is the GroupVersionResource for UpgradeOperation.
	UpgradeOperationGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: ResourceUpgradeOperations,
	}
)
