package v1alpha1

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// KindHealthCheckRequest is the Kind string for HealthCheckRequest resources.
	KindHealthCheckRequest = "HealthCheckRequest"

	// KindHealthSignal is the Kind string for HealthSignal resources.
	KindHealthSignal = "HealthSignal"

	// ResourceHealthCheckRequests is the plural resource name for HealthCheckRequest.
	ResourceHealthCheckRequests = "healthcheckrequests"

	// ResourceHealthSignals is the plural resource name for HealthSignal.
	ResourceHealthSignals = "healthsignals"

	// LabelUpgradeOperation is a label set on HealthCheckRequest and HealthSignal CRs
	// to identify the parent UpgradeOperation. This enables efficient label-selector filtering
	// on watches and lists.
	LabelUpgradeOperation = "upgrade.aks.io/operation"

	// DefaultHealthSignalTimeout is the maximum duration the RP waits for a health verdict
	// per node. If the timeout elapses with no "False" condition, the RP proceeds.
	DefaultHealthSignalTimeout = 5 * time.Minute
)

var (
	// HealthCheckRequestGVR is the GroupVersionResource for HealthCheckRequest.
	HealthCheckRequestGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: ResourceHealthCheckRequests,
	}

	// HealthSignalGVR is the GroupVersionResource for HealthSignal.
	HealthSignalGVR = schema.GroupVersionResource{
		Group:    GroupVersion.Group,
		Version:  GroupVersion.Version,
		Resource: ResourceHealthSignals,
	}
)
