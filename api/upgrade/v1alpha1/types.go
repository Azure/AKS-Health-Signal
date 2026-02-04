// Package v1alpha1 contains API Schema definitions for the upgrade.aks.io v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=upgrade.aks.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeReference contains a reference to a Node
type NodeReference struct {
	// Name is the name of the Node
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=unip

// UpgradeNodeInProgress is the Schema for the upgradenodeinprogresses API
type UpgradeNodeInProgress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	NodeRef NodeReference `json:"nodeRef"`
}

// +kubebuilder:object:root=true

// UpgradeNodeInProgressList contains a list of UpgradeNodeInProgress
type UpgradeNodeInProgressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeNodeInProgress `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=uoip

// UpgradeOperationInProgress is the Schema for the upgradeoperationinprogresses API
// The existence of this CR indicates an upgrade operation is in progress.
type UpgradeOperationInProgress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

// UpgradeOperationInProgressList contains a list of UpgradeOperationInProgress
type UpgradeOperationInProgressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeOperationInProgress `json:"items"`
}
