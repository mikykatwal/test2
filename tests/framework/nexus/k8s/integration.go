package k8s

import (
	"context"

	"go.mondoo.com/mondoo-operator/tests/framework/nexus/api/integrations"
)

type IntegrationBuilder struct {
	integrations integrations.IntegrationsManager

	spaceMrn            string
	name                string
	scanNodes           bool
	scanWorkloads       bool
	scanContainerImages bool
	admissionController bool
}

func (i *IntegrationBuilder) EnableNodesScan() *IntegrationBuilder {
	i.scanNodes = true
	return i
}

func (i *IntegrationBuilder) EnableWorkloadsScan() *IntegrationBuilder {
	i.scanWorkloads = true
	return i
}

func (i *IntegrationBuilder) EnableContainerImagesScan() *IntegrationBuilder {
	i.scanContainerImages = true
	return i
}

func (i *IntegrationBuilder) EnableAdmissionController() *IntegrationBuilder {
	i.admissionController = true
	return i
}

func (b *IntegrationBuilder) Run(ctx context.Context) (*Integration, error) {
	resp, err := b.integrations.Create(ctx, &integrations.CreateIntegrationRequest{
		Name:     b.name,
		SpaceMrn: b.spaceMrn,
		Type:     integrations.Type_K8S,
		ConfigurationInput: &integrations.ConfigurationInput{
			ConfigOptions: &integrations.ConfigurationInput_K8SOptions{
				K8SOptions: &integrations.K8SConfigurationOptionsInput{
					ScanNodes:        b.scanNodes,
					ScanWorkloads:    b.scanWorkloads,
					ScanPublicImages: b.scanContainerImages,
					ScanDeploys:      b.admissionController,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &Integration{integrations: b.integrations, mrn: resp.Integration.Mrn}, nil
}

type Integration struct {
	integrations integrations.IntegrationsManager

	mrn string
}

func (i *Integration) GetLongLivedToken(ctx context.Context) (string, error) {
	resp, err := i.integrations.GetTokenForIntegration(
		ctx, &integrations.GetTokenForIntegrationRequest{Mrn: i.mrn, LongLivedToken: true})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (i *Integration) Delete(ctx context.Context) error {
	_, err := i.integrations.Delete(ctx, &integrations.DeleteIntegrationRequest{Mrn: i.mrn})
	return err
}