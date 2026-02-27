// Package v1alpha1 contains API Schema definitions for the health.aks.io v1alpha1 API group.
//
// This package defines two custom resources:
//
//   - HealthCheckRequest: Created by the AKS RP to request health monitoring for a
//     node, node pool, or cluster during upgrade operations. Cluster-scoped.
//
//   - HealthSignal: Created and updated by monitoring apps (e.g., DaemonSets) in
//     response to a HealthCheckRequest. Contains health conditions that the RP reads
//     to decide whether to proceed or abort an upgrade. Namespaced (kube-system).
//
// Linkage: Each HealthSignal sets an ownerReference to its HealthCheckRequest,
// ensuring automatic garbage collection when the request is deleted.
//
// +kubebuilder:object:generate=true
// +groupName=health.aks.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HealthCheckRequestScope defines the level of the health check request.
// +kubebuilder:validation:Enum=Node;NodePool;Cluster
type HealthCheckRequestScope string

const (
	// HealthCheckRequestScopeNode targets a single node.
	HealthCheckRequestScopeNode HealthCheckRequestScope = "Node"

	// HealthCheckRequestScopeNodePool targets an entire node pool.
	HealthCheckRequestScopeNodePool HealthCheckRequestScope = "NodePool"

	// HealthCheckRequestScopeCluster targets the whole cluster.
	HealthCheckRequestScopeCluster HealthCheckRequestScope = "Cluster"
)

// TargetRef identifies a Kubernetes object by name.
type TargetRef struct {
	// Name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// HealthCheckRequestSpec defines the desired state of a HealthCheckRequest.
type HealthCheckRequestSpec struct {
	// Scope indicates the level of this health check: Node, NodePool, or Cluster.
	// +kubebuilder:validation:Required
	Scope HealthCheckRequestScope `json:"scope"`

	// TargetRef references the target object (e.g., a Node or ).
	// +kubebuilder:validation:Required
	TargetRef *TargetRef `json:"targetRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=hcr

// HealthCheckRequest is the Schema for the healthcheckrequests API.
//
// Created to request health monitoring for a specific target
// (node, node pool, or cluster). Monitoring apps watch for these
// and create corresponding HealthSignal resources with ownerReferences
// pointing back to the request.
type HealthCheckRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec HealthCheckRequestSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// HealthCheckRequestList contains a list of HealthCheckRequest
type HealthCheckRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HealthCheckRequest `json:"items"`
}
