// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package webhookhandler

import (
	"context"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/types/known/structpb"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/yaml"

	serializerYaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
	mondoov1alpha2 "go.mondoo.com/mondoo-operator/api/v1alpha2"
	"go.mondoo.com/mondoo-operator/pkg/client/scanapiclient"
	"go.mondoo.com/mondoo-operator/pkg/constants"
	"go.mondoo.com/mondoo-operator/pkg/feature_flags"
	"go.mondoo.com/mondoo-operator/pkg/utils"
	wutils "go.mondoo.com/mondoo-operator/pkg/webhooks/utils"
)

// Have kubebuilder generate a ValidatingWebhookConfiguration under the path /validate-k8s-mondoo-com that watches Pod/Deployment creation/updates
//+kubebuilder:webhook:path=/validate-k8s-mondoo-com,mutating=false,failurePolicy=ignore,sideEffects=None,groups="";apps;batch,resources=pods;deployments;daemonsets;statefulsets;jobs;cronjobs,verbs=create;update,versions=v1,name=policy.k8s.mondoo.com,admissionReviewVersions=v1

var handlerlog = logf.Log.WithName("webhook-validator")

const (
	// defaultScanPass is our default Allowed result in the event that we never even made it to the scan
	defaultScanPass = "DEFAULT MONDOO PASSED"
	// defaultScanFail is our default Failed result in the event that we never even made it to the scan
	defaultScanFail = "DEFAULT MONDOO FAILED"
	// passedScan is the Allowed result when the sacn came back with a passing result
	passedScan = "PASSED MONDOO SCAN"
	// failedScanPermitted is the Allowed result when in Permissive mode but the scan result was a failing result
	failedScanPermitted = "PERMITTING FAILED SCAN"
	// failedScan is the Denied result when in Enforcing mode and the scan result was a failing result
	failedScan = "FAILED MONDOO SCAN"

	mondooLabelPrefix          = "k8s.mondoo.com/"
	mondooNamespaceLabel       = mondooLabelPrefix + "namespace"
	mondooUIDLabel             = mondooLabelPrefix + "uid"
	mondooResourceVersionLabel = mondooLabelPrefix + "resource-version"
	mondooNameLabel            = mondooLabelPrefix + "name"
	mondooKindLabel            = mondooLabelPrefix + "kind"
	mondooOwnerNameLabel       = mondooLabelPrefix + "owner-name"
	mondooOwnerKindLabel       = mondooLabelPrefix + "owner-kind"
	mondooOwnerUIDLabel        = mondooLabelPrefix + "owner-uid"
	mondooAuthorLabel          = mondooLabelPrefix + "author"
	mondooOperationLabel       = mondooLabelPrefix + "operation"
	mondooClusterIDLabel       = mondooLabelPrefix + "cluster-id"
)

type webhookValidator struct {
	client            client.Client
	decoder           *admission.Decoder
	mode              mondoov1alpha2.AdmissionMode
	scanner           scanapiclient.ScanApiClient
	integrationMRN    string
	clusterID         string
	uniDecoder        runtime.Decoder
	includeNamespaces []string
	excludeNamespaces []string
}

type NewWebhookValidatorOpts struct {
	Client            client.Client
	Mode              string
	ScanUrl           string
	Token             string
	IntegrationMrn    string
	ClusterId         string
	IncludeNamespaces []string
	ExcludeNamespaces []string
}

// NewWebhookValidator will initialize a CoreValidator with the provided k8s Client and
// set it to the provided mode. Returns error if mode is invalid.
func NewWebhookValidator(opts *NewWebhookValidatorOpts) (admission.Handler, error) {
	webhookMode, err := wutils.ModeStringToAdmissionMode(opts.Mode)
	if err != nil {
		return nil, err
	}

	clnt, err := scanapiclient.NewClient(scanapiclient.ScanApiClientOptions{
		ApiEndpoint: opts.ScanUrl,
		Token:       opts.Token,
	})
	if err != nil {
		return nil, err
	}

	return &webhookValidator{
		client:            opts.Client,
		mode:              webhookMode,
		scanner:           clnt,
		integrationMRN:    opts.IntegrationMrn,
		clusterID:         opts.ClusterId,
		uniDecoder:        serializer.NewCodecFactory(opts.Client.Scheme()).UniversalDeserializer(),
		includeNamespaces: opts.IncludeNamespaces,
		excludeNamespaces: opts.ExcludeNamespaces,
	}, nil
}

var _ admission.Handler = &webhookValidator{}

func (a *webhookValidator) Handle(ctx context.Context, req admission.Request) (response admission.Response) {
	resource := fmt.Sprintf("%s/%s", req.Namespace, req.Name)
	handlerlog.Info("Webhook triggered", "kind", req.Kind.Kind, "resource", resource)

	// the default/safe response
	response = admission.Allowed(defaultScanPass)
	if a.mode == mondoov1alpha2.Enforcing {
		response = admission.Denied(defaultScanFail)
	}

	obj, err := a.objFromRaw(req.Object)
	if err == nil {
		if !shouldScanObject(obj) {
			handlerlog.Info("skipping because the resource has a parent", "resource", resource)
			return
		}
	}

	skip, err := a.skipNamespace(obj)
	if err != nil {
		handlerlog.Error(err, "error while checking whether to skip resource based on namespace")
		return
	}
	if skip {
		handlerlog.Info("skipping based on namespace filtering", "resource", resource)
		return
	}

	// updates with server-side apply (SSA) might only change the resourceVersion
	// this is a known bug: https://github.com/kubernetes/kubernetes/issues/95460
	// it is fixed in k8s 1.25, but we need to support older versions
	// not skipping these updates floods our CI/CD with GKE, AKS, ... addon updates
	// SSA also changes the managedFields metadata, so we drop them also to get rid of the noise
	if req.AdmissionRequest.Operation == admissionv1.Update && req.AdmissionRequest.OldObject.Raw != nil {
		skip, err := objectsOnlyDifferInSSAFields(req.AdmissionRequest)
		if err != nil {
			handlerlog.Error(err, "failed to get difference between objects")
			return
		}
		if skip {
			handlerlog.V(9).Info("skipping because the old and new object only differ in resourceVersion; happens with server-side apply")
			return
		}
	}

	k8sLabels, err := a.generateLabels(req, obj)
	if err != nil {
		handlerlog.Error(err, "failed to set labels for incoming request")
		return
	}

	// Call into Mondoo Scan Service to scan the resource
	reqData, err := yaml.Marshal(admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{Kind: "AdmissionReview", APIVersion: "admission.k8s.io/v1"},
		Request:  &req.AdmissionRequest,
	})
	if err != nil {
		handlerlog.Error(err, "failed to marshal incoming request")
		return
	}

	mapData := make(map[string]interface{})
	if err := yaml.Unmarshal(reqData, &mapData); err != nil {
		handlerlog.Error(err, "failed to unmarshal object to map")
		return
	}

	data, err := structpb.NewStruct(mapData)
	if err != nil {
		handlerlog.Error(err, "failed to create proto struct from admission request")
		return
	}
	scanJob := &scanapiclient.AdmissionReviewJob{
		Data:       data,
		Labels:     k8sLabels,
		ReportType: scanapiclient.ReportType_ERROR,
	}

	scanJob.Discovery = &inventory.Discovery{}
	scanJob.Options = map[string]string{"all-namespaces": "true"}
	// do not use auto discovery here, because we do not want to scan the cluster
	scanJob.Discovery.Targets = []string{"pods", "deployments", "daemonsets", "statefulsets", "replicasets", "jobs", "cronjobs"}

	if feature_flags.GetAdmissionReviewDiscovery() {
		scanJob.Discovery.Targets = append(scanJob.Discovery.Targets, "admissionreviews")
	}

	result, err := a.scanner.RunAdmissionReview(ctx, scanJob)
	if err != nil {
		handlerlog.Error(err, "error returned from scan request")
		return
	}

	passed := false
	if result.WorstScore != nil && result.WorstScore.Type == scanapiclient.ValidScanResult && result.WorstScore.Value == 100 {
		passed = true
	}

	handlerlog.Info("Scan result", "shouldAdmit", passed, "kind", req.Kind.Kind, "resource", resource, "worstscore", result.WorstScore)

	// Depending on the mode, we either just allow the resource through no matter the scan result
	// or allow/deny based on the scan result
	switch a.mode {
	case mondoov1alpha2.Permissive:
		if passed {
			response = admission.Allowed(passedScan)
		} else {
			response = admission.Allowed(failedScanPermitted)
		}
	case mondoov1alpha2.Enforcing:
		if passed {
			response = admission.Allowed(passedScan)
		} else {
			response = admission.Denied(failedScan)
		}
	default:
		err := fmt.Errorf("neither permissive nor enforcing modes defined")
		handlerlog.Error(err, "unexpected runtime environment, allowing the resource through")
	}
	return
}

func (a *webhookValidator) generateLabels(req admission.Request, obj runtime.Object) (map[string]string, error) {
	labels, err := generateLabelsFromAdmissionRequest(req, obj)
	if err != nil {
		return nil, err
	}

	labels[mondooClusterIDLabel] = a.clusterID
	if a.integrationMRN != "" {
		labels[constants.MondooAssetsIntegrationLabel] = a.integrationMRN
	}

	return labels, nil
}

func (a *webhookValidator) objFromRaw(rawObj runtime.RawExtension) (runtime.Object, error) {
	obj, _, err := a.uniDecoder.Decode(rawObj.Raw, nil, nil)
	if err != nil {
		obj, _, err = serializerYaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return obj, err
		}
	}
	return obj, err
}

func (a *webhookValidator) skipNamespace(obj runtime.Object) (bool, error) {
	objmeta, err := meta.Accessor(obj)
	if err != nil {
		handlerlog.Error(err, "error getting metadata from object", "type", reflect.TypeOf(obj))
		return false, nil
	}

	allow, err := utils.AllowNamespace(objmeta.GetNamespace(), a.includeNamespaces, a.excludeNamespaces)
	return !allow, err
}

func generateLabelsFromAdmissionRequest(req admission.Request, obj runtime.Object) (map[string]string, error) {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{
		mondooNamespaceLabel:       objMeta.GetNamespace(),
		mondooUIDLabel:             string(objMeta.GetUID()),
		mondooResourceVersionLabel: objMeta.GetResourceVersion(),
		mondooNameLabel:            objMeta.GetName(),
		mondooKindLabel:            obj.GetObjectKind().GroupVersionKind().Kind,
		mondooAuthorLabel:          req.UserInfo.Username,
		mondooOperationLabel:       string(req.Operation),
	}

	controllerRef := metav1.GetControllerOf(objMeta)
	if controllerRef != nil {
		labels[mondooOwnerNameLabel] = controllerRef.Name
		labels[mondooOwnerKindLabel] = controllerRef.Kind
		labels[mondooOwnerUIDLabel] = string(controllerRef.UID)
	}

	return labels, nil
}

// shouldScanObject determines whether an object should be scanned by Mondoo client.
func shouldScanObject(obj runtime.Object) bool {
	objMeta, err := meta.Accessor(obj)
	if err == nil {
		controller := metav1.GetControllerOf(objMeta)
		if controller != nil {
			// Don't scan objects which parent we have already scanned
			return controller.Kind != "Deployment" &&
				controller.Kind != "ReplicaSet" &&
				controller.Kind != "DaemonSet" &&
				controller.Kind != "StatefulSet" &&
				controller.Kind != "CronJob" &&
				controller.Kind != "Job"
		}
	}

	// In case we couldn't access the meta object or there was no controller owner, then scan
	return true
}

func objectsOnlyDifferInSSAFields(admissionRequest admissionv1.AdmissionRequest) (bool, error) {
	oldObjMapData := make(map[string]interface{})
	if err := yaml.Unmarshal(admissionRequest.OldObject.Raw, &oldObjMapData); err != nil {
		handlerlog.Error(err, "failed to unmarshal old object to map")
		return false, err
	}

	objMapData := make(map[string]interface{})
	if err := yaml.Unmarshal(admissionRequest.Object.Raw, &objMapData); err != nil {
		handlerlog.Error(err, "failed to unmarshal new object to map")
		return false, err
	}

	// we don't care about the actual difference, only whether the resourceVersion is the only difference or not
	oldObjMapData["metadata"].(map[string]interface{})["resourceVersion"] = ""
	objMapData["metadata"].(map[string]interface{})["resourceVersion"] = ""

	// also ignore managedFields which get changed during SSA
	oldObjMapData["metadata"].(map[string]interface{})["managedFields"] = ""
	objMapData["metadata"].(map[string]interface{})["managedFields"] = ""

	if reflect.DeepEqual(oldObjMapData, objMapData) {
		return true, nil
	}

	return false, nil
}
