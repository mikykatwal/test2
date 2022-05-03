package utils

import (
	mondoov2 "go.mondoo.com/mondoo-operator/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MondooClientSecret = "mondoo-client"

// DefaultAuditConfig returns a new Mondoo audit config with some default settings to
// make sure a tests passes (e.g. setting the correct secret name).
func DefaultAuditConfig(ns string, workloads, nodes, admission bool) mondoov2.MondooAuditConfig {
	return mondoov2.MondooAuditConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "mondoo-client",
			Namespace: ns,
		},
		Spec: mondoov2.MondooAuditConfigSpec{
			MondooCredsSecretRef: corev1.LocalObjectReference{Name: MondooClientSecret},
			KubernetesResources:  mondoov2.KubernetesResources{Enable: workloads},
			Nodes:                mondoov2.Nodes{Enable: nodes},
			Admission:            mondoov2.Admission{Enable: admission},
		},
	}
}
