// Package v1alpha1 contains API Schema definitions for the upgrade.aks.io v1alpha1 API group.
//
// This package defines one custom resource:
//
//   - UpgradeOperation: Created by the AKS Resource Provider (RP) to represent
//     an in-progress upgrade of a cluster or agent pool. Cluster-scoped.
//
// Constraints:
//   - Multiple node-pool UpgradeOperations may coexist on the same cluster,
//     but each target (cluster or node pool) may have at most one active
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
type UpgradeOperationSpec struct {
	// Type indicates the upgrade type: Cluster or NodePool.
	// +kubebuilder:validation:Required
	Type UpgradeType `json:"type"`

	// TargetName is the name of the upgrade target.
	// For Cluster upgrades this is the cluster name; for NodePool upgrades
	// this is the node pool name.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	TargetName string `json:"targetName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=uo

// UpgradeOperation is the Schema for the upgradeoperations API.
//
// Created by the AKS Resource Provider to represent an in-progress upgrade of
// a cluster or node pool. Multiple node-pool UpgradeOperations may exist
// concurrently, but each target may have at most one active UpgradeOperation.
type UpgradeOperation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec UpgradeOperationSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// UpgradeOperationList contains a list of UpgradeOperation.
type UpgradeOperationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpgradeOperation `json:"items"`
}
