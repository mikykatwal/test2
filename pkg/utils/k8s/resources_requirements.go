package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// defaultMondooClientResources for Mondoo Client container
var DefaultMondooClientResources corev1.ResourceRequirements = corev1.ResourceRequirements{
	Limits: corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse("500M"),
		corev1.ResourceCPU:    resource.MustParse("400m"),
	},

	Requests: corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse("180M"),
		corev1.ResourceCPU:    resource.MustParse("150m"),
	},
}

// ResourcesRequirementsWithDefaults will return the resource requirements from the parameter if such
// are specified. If not requirements are specified, default values will be returned.
func ResourcesRequirementsWithDefaults(m corev1.ResourceRequirements) corev1.ResourceRequirements {
	if m.Size() != 0 {
		return m
	}

	// Default values for Mondoo resources requirements.
	return DefaultMondooClientResources
}
