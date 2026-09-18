package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Duration is a string representing a span of time.
//
// The accepted syntax is GEP-2257: an unsigned integer with a unit of h, m, s
// or ms, repeated up to four times ("5m", "1m30s", "500ms"). This is a strict
// subset of what Go's time.ParseDuration accepts — fractional values ("1.5h"),
// sub-millisecond units ("100ns", "2us") and negative values ("-5s") are all
// rejected. Negatives in particular cannot mean anything for a timeout, yet
// time.ParseDuration accepts them.
//
// This is a named string type rather than metav1.Duration because a pattern
// cannot be attached to metav1.Duration: it is a struct with custom JSON
// marshalling, so controller-gen emits a bare "type: string" with no constraint
// and rejects +kubebuilder:validation:Pattern on it. Without a pattern the API
// server admits any string whatsoever, and a malformed value is not discovered
// until the RP parses it, part-way through a customer's upgrade. The wire
// format is identical either way, and sigs.k8s.io/gateway-api makes the same
// trade for the same reason (apis/v1.Duration).
//
// Consumers parse with time.ParseDuration(string(d)); every value admitted by
// the pattern is accepted by it.
//
// +kubebuilder:validation:Pattern=`^([0-9]{1,5}(h|m|s|ms)){1,4}$`
type Duration string

// OnFailureAction defines what the AKS Resource Provider does when a provider
// reports an unhealthy verdict, or fails to report one before NodeTimeout.
// Only Abort is supported.
// +kubebuilder:validation:Enum=Abort
type OnFailureAction string

const (
	// OnFailureAbort halts the upgrade operation and returns an error.
	OnFailureAbort OnFailureAction = "Abort"
)

// UpgradeGateTarget selects which nodes a rule applies to.
// +kubebuilder:validation:XValidation:rule="self.scope != 'NodePool' || (has(self.name) && size(self.name) > 0)",message="name is required when scope is NodePool"
type UpgradeGateTarget struct {
	// Scope indicates whether this rule targets the whole cluster or a single
	// node pool.
	// +kubebuilder:validation:Required
	Scope UpgradeType `json:"scope"`

	// Name is the node pool name this rule applies to.
	// Required when Scope is NodePool. Ignored when Scope is Cluster.
	// +optional
	Name string `json:"name,omitempty"`
}

// UpgradeGateRule declares that this policy's provider gates one target, and on
// what terms.
//
// After each node is upgraded the RP opens a window for this provider and watches
// its HealthSignals at every scope. NodeTimeout is the ceiling on that window.
type UpgradeGateRule struct {
	// Target selects the nodes this rule applies to.
	// +kubebuilder:validation:Required
	Target UpgradeGateTarget `json:"target"`

	// NodeTimeout is the longest the AKS Resource Provider waits for this
	// provider's verdict after each node is upgraded. The timer starts when the
	// node's health check begins, so it also bounds a provider that never creates
	// a HealthSignal at all.
	// Expressed in GEP-2257 duration syntax (e.g., "30s", "5m", "1m30s"):
	// whole units of h, m, s or ms only. Fractional ("1.5h") and negative values
	// are rejected.
	// Capped at 5 minutes; larger values are clamped. Omitted means the cap.
	//
	// The window closes as soon as the node reports healthy.
	//
	// Each provider's NodeTimeout is independent. A target that is gated by
	// several providers waits for each against its own deadline.
	//
	// +optional
	NodeTimeout *Duration `json:"nodeTimeout,omitempty"`

	// OnFailure is the action the RP takes when this provider reports unhealthy
	// or fails to report before NodeTimeout.
	// +optional
	// +kubebuilder:default=Abort
	OnFailure OnFailureAction `json:"onFailure,omitempty"`
}

// UpgradeGatePolicySpec defines the desired state of an UpgradeGatePolicy.
type UpgradeGatePolicySpec struct {
	// Provider is the identity of the monitoring app this policy registers.
	// It must match the "health.aks.io/provider" label the app sets on every
	// HealthSignal it creates.
	//
	// Constrained to the Kubernetes label value rules: at most 63 characters,
	// alphanumeric at each end, and only '-', '_' or '.' in between. '/' is not
	// permitted. Use reverse-DNS to avoid collisions between vendors, for
	// example "node-monitor.example.com".
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9]([-A-Za-z0-9_.]*[A-Za-z0-9])?$`
	Provider string `json:"provider"`

	// Rules are the targets this provider gates. A provider gating two node
	// pools on different terms lists two rules.
	//
	// When more than one rule matches a target, the most specific rule applies:
	// a NodePool rule overrides a Cluster rule for nodes in that pool.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	Rules []UpgradeGateRule `json:"rules"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=ugp

// UpgradeGatePolicy is the Schema for the upgradegatepolicies API.
//
// Installed by a monitoring app, typically from its Helm chart, to register
// that it gates upgrades for a set of nodes. It is install-time configuration:
// created once and unchanged across upgrade operations.
//
// During an upgrade the RP matches each node against the installed policies. A
// node that matches no policy is not gated and receives no HealthCheckRequest;
// a node that matches waits for every matched provider, each bounded by its own
// NodeTimeout. HealthSignals from a provider with no policy are ignored.
//
// One object per provider. Written only by the monitoring app; the RP reads it.
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
