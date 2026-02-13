// Package v1alpha1 contains API Schema definitions for the health.aks.io v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=health.aks.io
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HealthSignalType defines the type of health signal
// +kubebuilder:validation:Enum=NodeHealth;ClusterHealth
type HealthSignalType string

const (
	NodeHealth    HealthSignalType = "NodeHealth"
	ClusterHealth HealthSignalType = "ClusterHealth"
)

// HealthSignalSpec defines the desired state of HealthSignal
type HealthSignalSpec struct {
	// Type is the health signal type
	// +kubebuilder:validation:Required
	Type HealthSignalType `json:"type"`

	// Target of the health signal. Required when type=NodeHealth
	// +optional
	Target *corev1.ObjectReference `json:"target,omitempty"`

	// Source is the controller/component that produced this health signal
	// +kubebuilder:validation:Required
	Source corev1.ObjectReference `json:"source"`

	// TimeoutSeconds is the maximum duration in seconds that RP should wait for
	// the health signal to reach a verdict. Defaults to 300 (5 minutes) if not set.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=300
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
}

// HealthSignalStatus defines the observed state of HealthSignal
type HealthSignalStatus struct {
	// Conditions represent the latest available observations of the HealthSignal's state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=hs
// +kubebuilder:validation:XValidation:rule="self.spec.type != 'NodeHealth' || has(self.spec.target) && has(self.spec.target.name) && self.spec.target.name != ''",message="spec.target.name is required when spec.type is NodeHealth."
// +kubebuilder:validation:XValidation:rule="self.spec.type != 'ClusterHealth' || !has(self.spec.target) || !has(self.spec.target.name)",message="spec.target must not be set when spec.type is ClusterHealth."

// HealthSignal is the Schema for the healthsignals API
type HealthSignal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   HealthSignalSpec   `json:"spec"`
	Status HealthSignalStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HealthSignalList contains a list of HealthSignal
type HealthSignalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HealthSignal `json:"items"`
}
