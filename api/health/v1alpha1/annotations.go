package v1alpha1

const (
	// AnnotationUpgradeCorrelationID is a UUID set by the AKS Resource Provider (RP)
	// that links HealthCheckRequest and HealthSignal CRs to a specific upgrade operation.
	// Format: UUID string (e.g., "6e8ef28e-bb8a-42cb-aa0b-d05a05b1ba0a")
	AnnotationUpgradeCorrelationID = "kubernetes.azure.com/upgradeCorrelationID"

	// AnnotationTargetKubernetesVersion is the target Kubernetes version for the upgrade.
	// Format: semver string (e.g., "1.33.5")
	AnnotationTargetKubernetesVersion = "kubernetes.azure.com/targetKubernetesVersion"
)
