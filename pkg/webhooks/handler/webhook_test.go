package webhookhandler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	mondoov1alpha1 "go.mondoo.com/mondoo-operator/api/v1alpha1"
	"go.mondoo.com/mondoo-operator/pkg/scanner"
	"go.mondoo.com/mondoo-operator/pkg/scanner/fakescanapi"
)

func TestWebhookValidate(t *testing.T) {

	decoder := setupDecoder(t)
	tests := []struct {
		name          string
		mode          mondoov1alpha1.WebhookMode
		expectAllowed bool
		expectReason  string
		object        runtime.RawExtension
	}{
		{
			name:          "example test",
			expectAllowed: true,
			expectReason:  "PASSED MONDOO SCAN",
			object:        testExamplePod(),
		},
		{
			name:          "example Deployment",
			expectAllowed: true,
			expectReason:  "PASSED MONDOO SCAN",
			object:        testExampleDeployment(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			if test.mode == "" {
				test.mode = mondoov1alpha1.Permissive
			}

			testserver := fakescanapi.FakeServer()
			validator := &webhookValidator{
				decoder: decoder,
				mode:    test.mode,
				scanner: &scanner.Scanner{
					Endpoint: testserver.URL,
				},
			}

			request := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Object: test.object,
				},
			}

			// Act
			response := validator.Handle(context.TODO(), request)

			t.Logf("Handle() response: %+v", response)

			// Assert
			assert.Equal(t, test.expectAllowed, response.AdmissionResponse.Allowed)

			if test.expectReason != "" {
				assert.Equal(t, test.expectReason, string(response.AdmissionResponse.Result.Reason))
			}

		})
	}
}

func testExamplePod() runtime.RawExtension {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testPod",
			Namespace: "testNamespace",
		},
	}

	return runtime.RawExtension{
		Object: pod,
	}
}

func testExampleDeployment() runtime.RawExtension {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testDeployment",
			Namespace: "testNamespace",
		},
	}

	return runtime.RawExtension{
		Object: dep,
	}
}

func setupDecoder(t *testing.T) *admission.Decoder {
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	decoder, err := admission.NewDecoder(scheme)
	require.NoError(t, err, "Failed to setup decoder for testing")

	return decoder
}
