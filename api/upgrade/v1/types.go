// Package v1 contains API Schema definitions for the upgrade.aks.io v1 API group
// +kubebuilder:object:generate=true
// +groupName=upgrade.aks.io
package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpgradeNodeInProgressStatus defines the observed state of UpgradeNodeInProgress
type UpgradeNodeInProgressStatus struct {
	// Conditions represent the latest available observations of the upgrade state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// Condition types for UpgradeNodeInProgress
const (
	// ConditionInProgress indicates the node upgrade is in progress
	ConditionInProgress string = "InProgress"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=unip

// UpgradeNodeInProgress is the Schema for the upgradenodeinprogresses API
type UpgradeNodeInProgress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	NodeRef corev1.ObjectReference      `json:"nodeRef"`
	Status  UpgradeNodeInProgressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UpgradeNodeInProgressList contains a list of UpgradeNodeInProgress
type UpgradeNodeInProgressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeNodeInProgress `json:"items"`
}
