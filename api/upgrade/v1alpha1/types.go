// Package v1alpha1 contains API Schema definitions for the upgrade.aks.io v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=upgrade.aks.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeReference is a reference to a node
type NodeReference struct {
	// Name is the name of the node
	// +required
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`
}

// UpgradeNodeInProgressSpec defines the desired state of UpgradeNodeInProgress
type UpgradeNodeInProgressSpec struct {
	// NodeRef is a reference to the node being upgraded
	// +optional
	NodeRef *NodeReference `json:"nodeRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=unip

// UpgradeNodeInProgress is the Schema for the upgradenodeinprogresses API
type UpgradeNodeInProgress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec UpgradeNodeInProgressSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// UpgradeNodeInProgressList contains a list of UpgradeNodeInProgress
type UpgradeNodeInProgressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeNodeInProgress `json:"items"`
}

// HealthSignalSpec defines the desired state of HealthSignal
type HealthSignalSpec struct {
	// Source identifies the checker that created this health signal
	// For example: "ClusterHealthMonitor", "PrometheusMetrics"
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Source string `json:"source"`
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
