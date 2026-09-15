// Package v1alpha1 contains API Schema definitions for the health.aks.io v1alpha1 API group.
//
// This package defines two custom resources:
//
//   - HealthCheckRequest: Created by the AKS RP to request health monitoring for a
//     node, node pool, or cluster during upgrade operations. Cluster-scoped.
//
//   - HealthSignal: Created and updated by monitoring apps (e.g., DaemonSets) in
//     response to a HealthCheckRequest. Contains health conditions that the RP reads
//     to decide whether to proceed or abort an upgrade. Cluster-scoped.
//
// Linkage: Each HealthSignal sets an ownerReference to its HealthCheckRequest,
// ensuring automatic garbage collection when the request is deleted.
//
// +kubebuilder:object:generate=true
// +groupName=health.aks.io
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HealthSignalType defines the type of health signal. It corresponds to the
// scope of the HealthCheckRequest being answered: Node, NodePool or Cluster.
// +kubebuilder:validation:Enum=NodeHealth;NodePoolHealth;ClusterHealth
type HealthSignalType string

const (
	// NodeHealth indicates a per-node health signal.
	NodeHealth HealthSignalType = "NodeHealth"

	// NodePoolHealth indicates a node-pool-wide health signal.
	NodePoolHealth HealthSignalType = "NodePoolHealth"

	// ClusterHealth indicates a cluster-wide health signal.
	ClusterHealth HealthSignalType = "ClusterHealth"
)

const (
	// Healthy means the target is healthy and operating normally.
	Healthy string = "True"

	// Unhealthy means the target is unhealthy.
	Unhealthy string = "False"

	// Ongoing means monitoring has not yet reached a verdict.
	Ongoing string = "Unknown"
)

// HealthSignalSpec defines the desired state of a HealthSignal.
//
// HealthSignal is written entirely by monitoring apps. The AKS RP reads it.
// Decision logic, applied after each node is upgraded:
//   - If any condition status becomes "False", the RP aborts the upgrade. This
//     holds at every scope: an unhealthy NodePool- or Cluster-scoped signal stops
//     the rollout just as a node-scoped one does.
//   - If the node-scoped conditions the RP is waiting on are all "True", the
//     upgrade proceeds to the next node.
//   - "Unknown" on a node-scoped signal means no verdict yet; the RP waits until
//     the nodeTimeout configured on UpgradeGatePolicy (upgrade.aks.io/v1alpha1)
//     elapses, then aborts.
//   - "Unknown" on a NodePool- or Cluster-scoped signal does not block. Those
//     scopes are never awaited, so a monitor that is slow or never reports cannot
//     stall an upgrade. Only an explicit "False" from them has an effect.
type HealthSignalSpec struct {
	// Type is the health signal type (e.g., NodeHealth).
	// +kubebuilder:validation:Required
	Type HealthSignalType `json:"type"`

	// TargetRef optionally identifies the Kubernetes object this health signal
	// targets. The target is already carried by the HealthCheckRequest this
	// signal owner-references, which names the node, node pool or cluster. Set it
	// only as a convenience for readers of the signal alone.
	// +optional
	TargetRef *corev1.ObjectReference `json:"targetRef,omitempty"`
}

// HealthSignalStatus defines the observed state of a HealthSignal.
type HealthSignalStatus struct {
	// Conditions represent health verdicts using standard Kubernetes metav1.Condition.
	//
	// Each condition has:
	//   type              – condition type (e.g., "Ready")
	//   status            – "True" (Healthy), "False" (Unhealthy), "Unknown" (No verdict)
	//   reason            – machine-readable CamelCase reason (e.g., "Baseline", "NotReady")
	//   message           – human-readable description
	//   lastTransitionTime – when the condition last changed
	//
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

// HealthSignal is the Schema for the healthsignals API.
//
// Created and updated by monitoring apps in response to a HealthCheckRequest.
//
// Ownership: each HealthSignal MUST set an ownerReference pointing to its
// corresponding HealthCheckRequest. This ensures:
//   - Automatic garbage collection when the request is deleted.
//   - Clear linkage between request and signal.
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
