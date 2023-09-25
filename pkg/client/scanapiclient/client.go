package scanapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mondoo.com/cnquery/motor/asset"
	v1 "go.mondoo.com/cnquery/motor/inventory/v1"
	"go.mondoo.com/cnquery/motor/providers"
	"go.mondoo.com/cnspec/policy/scan"
	"go.mondoo.com/mondoo-operator/pkg/client/common"
	"go.mondoo.com/mondoo-operator/pkg/constants"
)

const (
	RunAdmissionReviewEndpoint             = "/Scan/RunAdmissionReview"
	ScanKubernetesResourcesEndpoint        = "/Scan/Run"
	ScheduleKubernetesResourceScanEndpoint = "/Scan/Schedule"
	GarbageCollectAssetsEndpoint           = "/Scan/GarbageCollectAssets"
)

type ScanApiClientOptions struct {
	ApiEndpoint string
	Token       string
}

type scanApiClient struct {
	ApiEndpoint string
	Token       string
	httpClient  http.Client
}

func NewClient(opts ScanApiClientOptions) (ScanApiClient, error) {
	opts.ApiEndpoint = strings.TrimRight(opts.ApiEndpoint, "/")
	client, err := common.DefaultHttpClient(nil)
	if err != nil {
		return nil, err
	}
	return &scanApiClient{
		ApiEndpoint: opts.ApiEndpoint,
		Token:       opts.Token,
		httpClient:  client,
	}, nil
}

func (s *scanApiClient) HealthCheck(ctx context.Context, in *common.HealthCheckRequest) (*common.HealthCheckResponse, error) {
	url := s.ApiEndpoint + common.HealthCheckEndpoint

	reqBodyBytes, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBodyBytes, err := common.Request(ctx, s.httpClient, url, s.Token, reqBodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	out := &common.HealthCheckResponse{}
	if err = json.Unmarshal(respBodyBytes, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto response: %v", err)
	}

	return out, nil
}

func (s *scanApiClient) RunAdmissionReview(ctx context.Context, in *AdmissionReviewJob) (*ScanResult, error) {
	url := s.ApiEndpoint + RunAdmissionReviewEndpoint

	reqBodyBytes, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBodyBytes, err := common.Request(ctx, s.httpClient, url, s.Token, reqBodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	out := &ScanResult{}
	if err = json.Unmarshal(respBodyBytes, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto response: %v", err)
	}

	return out, nil
}

func (s *scanApiClient) ScanKubernetesResources(ctx context.Context, scanOpts *ScanKubernetesResourcesOpts) (*ScanResult, error) {
	url := s.ApiEndpoint + ScanKubernetesResourcesEndpoint
	scanJob := &ScanJob{
		ReportType: ReportType_ERROR,
		Inventory: v1.Inventory{
			Spec: &v1.InventorySpec{
				Assets: []*asset.Asset{
					{
						Connections: []*providers.Config{
							{
								Backend: providers.ProviderType_K8S,
								Options: map[string]string{
									"namespaces":         strings.Join(scanOpts.IncludeNamespaces, ","),
									"namespaces-exclude": strings.Join(scanOpts.ExcludeNamespaces, ","),
								},
								Discover: &providers.Discovery{
									Targets: []string{"auto"},
								},
							},
						},
						ManagedBy: scanOpts.ManagedBy,
					},
				},
			},
		},
	}

	setIntegrationMrn(scanOpts.IntegrationMrn, scanJob)

	if scanOpts.ScanContainerImages {
		scanJob.Inventory.Spec.Assets[0].Connections[0].Discover.Targets = []string{"container-images"}
	}

	reqBodyBytes, err := json.Marshal(scanJob)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBodyBytes, err := common.Request(ctx, s.httpClient, url, s.Token, reqBodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	out := &ScanResult{}
	if err = json.Unmarshal(respBodyBytes, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto response: %v", err)
	}

	return out, nil
}

func (s *scanApiClient) ScheduleKubernetesResourceScan(ctx context.Context, integrationMrn, resourceKey, managedBy string) (*Empty, error) {
	url := s.ApiEndpoint + ScheduleKubernetesResourceScanEndpoint
	scanJob := &ScanJob{
		ReportType: ReportType_ERROR,
		Inventory: v1.Inventory{
			Spec: &v1.InventorySpec{
				Assets: []*asset.Asset{
					{
						Connections: []*providers.Config{
							{
								Backend: providers.ProviderType_K8S,
								Options: map[string]string{
									"k8s-resources": resourceKey,
								},
								Discover: &providers.Discovery{
									Targets: []string{"auto"},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(managedBy) > 0 {
		scanJob.Inventory.Spec.Assets[0].ManagedBy = managedBy
	}

	setIntegrationMrn(integrationMrn, scanJob)

	reqBodyBytes, err := json.Marshal(scanJob)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBodyBytes, err := common.Request(ctx, s.httpClient, url, s.Token, reqBodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	out := &Empty{}
	if err = json.Unmarshal(respBodyBytes, out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto response: %v", err)
	}

	return out, nil
}

func (s *scanApiClient) GarbageCollectAssets(ctx context.Context, in *scan.GarbageCollectOptions) error {
	url := s.ApiEndpoint + GarbageCollectAssetsEndpoint

	reqBodyBytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	_, err = common.Request(ctx, s.httpClient, url, s.Token, reqBodyBytes)
	if err != nil {
		return fmt.Errorf("error calling GarbageCollectAssets: %s", err)
	}

	return nil
}

func setIntegrationMrn(integrationMrn string, scanJob *ScanJob) {
	if integrationMrn != "" {
		if scanJob.Inventory.Spec.Assets[0].Labels == nil {
			scanJob.Inventory.Spec.Assets[0].Labels = make(map[string]string)
		}
		scanJob.Inventory.Spec.Assets[0].Labels[constants.MondooAssetsIntegrationLabel] = integrationMrn
	}
}
