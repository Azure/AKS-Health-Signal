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

	// LabelHealthSignalProvider identifies which monitoring app produced a
	// HealthSignal. Every app MUST set it on every HealthSignal it creates.
	//
	// The AKS Resource Provider uses it to decide whether a signal comes from a
	// provider it is waiting for: the value is matched against
	// UpgradeGatePolicy.spec.provider (upgrade.aks.io/v1alpha1). A signal whose
	// provider the RP is not waiting for is ignored, and a provider the RP is
	// waiting for that never produces a matching signal fails at its configured
	// timeout. An unlabelled HealthSignal therefore cannot satisfy any gate.
	//
	// This is a label rather than an annotation so the value is constrained by
	// the API server to the label value rules (see the matching validation on
	// UpgradeGatePolicy.spec.provider), and so signals can be selected by
	// provider, e.g. kubectl get healthsignals -l health.aks.io/provider=<name>.
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
