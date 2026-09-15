package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OnFailureAction defines what the AKS Resource Provider does when a provider
// reports an unhealthy verdict, or fails to report one before its timeout.
//
// Only Abort is accepted today: the RP has no rollback implementation, so
// accepting Rollback would be a promise the API cannot keep. The enum is
// deliberately shipped with a single value rather than omitting the field, so a
// customer asking for Rollback receives an explicit validation error instead of
// a silent downgrade to Abort. Widening the enum later is additive and
// non-breaking, because every pre-existing object already reads as Abort.
// +kubebuilder:validation:Enum=Abort
type OnFailureAction string

const (
	// OnFailureAbort halts the upgrade operation and returns an error.
	OnFailureAbort OnFailureAction = "Abort"
)

// UpgradeGateTarget selects which nodes a rule applies to.
// +kubebuilder:validation:XValidation:rule="self.type != 'NodePool' || (has(self.name) && size(self.name) > 0)",message="name is required when type is NodePool"
type UpgradeGateTarget struct {
	// Type indicates whether this rule targets the whole cluster or a single
	// node pool.
	// +kubebuilder:validation:Required
	Type UpgradeType `json:"type"`

	// Name is the node pool name this rule applies to.
	// Required when Type is NodePool. Ignored when Type is Cluster.
	// +optional
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`
}

// UpgradeGateRule declares that this policy's provider gates one target, and on
// what terms.
type UpgradeGateRule struct {
	// Target selects the nodes this rule applies to.
	// +kubebuilder:validation:Required
	Target UpgradeGateTarget `json:"target"`

	// Timeout is how long the AKS RP waits for this provider's verdict on one
	// node before failing closed. The timer is armed when the node's health
	// check begins, not when the provider's HealthSignal appears, so a provider
	// that never creates a HealthSignal is still bounded by it.
	// Expressed as a Kubernetes duration (e.g., "5m", "1h30m").
	// Omitted means the platform default.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// OnFailure is the action the RP takes when this provider reports unhealthy
	// or fails to report before Timeout.
	// +optional
	// +kubebuilder:default=Abort
	OnFailure OnFailureAction `json:"onFailure,omitempty"`
}

// UpgradeGatePolicySpec defines the desired state of an UpgradeGatePolicy.
type UpgradeGatePolicySpec struct {
	// Provider is the identity of the monitoring app this policy registers.
	// It MUST match the "health.aks.io/provider" label that the app sets on
	// every HealthSignal it creates; that label is how the RP maps a signal
	// back to this policy.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Provider string `json:"provider"`

	// Rules are the targets this provider gates.
	//
	// listType=atomic is safe here because exactly one owner (the provider's own
	// Helm release) ever writes this object. A provider gating two node pools on
	// different terms lists two rules; two different providers use two separate
	// UpgradeGatePolicy objects, so they can never clobber each other and
	// uninstalling one provider cannot disturb another's gating.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	Rules []UpgradeGateRule `json:"rules"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=ugp

// UpgradeGatePolicy is the Schema for the upgradegatepolicies API.
//
// Installed by a monitoring app (typically shipped in its Helm chart, the same
// way an application ships its own PodDisruptionBudget) to register that it
// gates upgrades for some set of nodes. It is durable install-time
// configuration: it is created once and survives every upgrade operation.
//
// Registration is what arms the upgrade gate. During an upgrade the RP matches
// each node against the installed policies; a node no policy matches is not
// gated at all and gets no HealthCheckRequest, while a node with matches waits
// for exactly the matched providers, each bounded by its own timeout.
// HealthSignals from a provider that has no policy are ignored.
//
// One object per provider. Never written by the RP, which only reads it.
type UpgradeGatePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec UpgradeGatePolicySpec `json:"spec"`
}

// +kubebuilder:object:root=true

// UpgradeGatePolicyList contains a list of UpgradeGatePolicy.
type UpgradeGatePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeGatePolicy `json:"items"`
}
