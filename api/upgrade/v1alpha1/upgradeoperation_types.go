// Package v1alpha1 contains API Schema definitions for the upgrade.aks.io v1alpha1 API group.
//
// This package defines one custom resource:
//
//   - UpgradeOperation: Created by the AKS Resource Provider (RP) to represent
//     one or more in-progress upgrade targets on a cluster. Cluster-scoped.
//
// Constraints:
//   - Cluster and node-pool targets may coexist in the same UpgradeOperation.
//   - Each target (cluster or node pool) may have at most one active
//     UpgradeOperation at a time.
//
// +kubebuilder:object:generate=true
// +groupName=upgrade.aks.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpgradeType defines the type of upgrade operation.
// +kubebuilder:validation:Enum=Cluster;NodePool
type UpgradeType string

const (
	// UpgradeTypeCluster represents a cluster-wide upgrade.
	UpgradeTypeCluster UpgradeType = "Cluster"

	// UpgradeTypeNodePool represents a node-pool-level upgrade.
	UpgradeTypeNodePool UpgradeType = "NodePool"
)

// UpgradeOperationSpec defines the desired state of an UpgradeOperation.
// +kubebuilder:validation:XValidation:rule="self.type != 'Cluster' || size(self.targetNames) == 1",message="cluster upgrades must target exactly one cluster name"
// +kubebuilder:validation:XValidation:rule="self.targetNames.all(target, size(target) > 0)",message="targetNames entries must be non-empty"
type UpgradeOperationSpec struct {
	// Type indicates the upgrade type: Cluster or NodePool.
	// +kubebuilder:validation:Required
	Type UpgradeType `json:"type"`

	// TargetNames are the names of the upgrade targets.
	// For Cluster upgrades this list contains the cluster name; for NodePool
	// upgrades this list contains one or more node pool names.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	TargetNames []string `json:"targetNames"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=uo

// UpgradeOperation is the Schema for the upgradeoperations API.
//
// Created by the AKS Resource Provider to represent one or more in-progress
// upgrade targets on a cluster. Cluster and node-pool targets may coexist in
// the same UpgradeOperation, but each target may have at most one active
// UpgradeOperation at a time.
type UpgradeOperation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the list of upgrade targets carried by this operation.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:XValidation:rule="self.filter(item, item.type == 'Cluster').size() <= 1",message="at most one cluster entry is allowed"
	// +listType=atomic
	Spec []UpgradeOperationSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// UpgradeOperationList contains a list of UpgradeOperation.
type UpgradeOperationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeOperation `json:"items"`
}
