package v1alpha1

import (
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

	// LabelHealthSignalProvider identifies the monitoring app that produced a
	// HealthSignal. Required on every HealthSignal.
	//
	// The value is matched against UpgradeGatePolicy.spec.provider
	// (upgrade.aks.io/v1alpha1). A signal that is unlabelled, or labelled with a
	// provider the RP is not waiting for, does not satisfy any gate.
	LabelHealthSignalProvider = "health.aks.io/provider"
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
