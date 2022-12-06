package nexus

import (
	"context"

	"go.mondoo.com/mondoo-operator/tests/framework/nexus/api/captain"
	"go.mondoo.com/mondoo-operator/tests/framework/nexus/api/integrations"
	"go.mondoo.com/mondoo-operator/tests/framework/nexus/api/policy"
	"go.mondoo.com/mondoo-operator/tests/framework/nexus/k8s"
)

type Space struct {
	spaceMrn string

	AssetStore   policy.AssetStore
	ReportsStore policy.ReportsStore
	Captain      captain.Captain
	Integrations integrations.IntegrationsManager

	K8s *k8s.Client
}

type AssetWithScore struct {
	Asset *policy.Asset
	Score *policy.Score
}

func NewSpace(spaceMrn string, assetStore policy.AssetStore, reportsStore policy.ReportsStore, captain captain.Captain, integrations integrations.IntegrationsManager) *Space {
	return &Space{
		spaceMrn:     spaceMrn,
		AssetStore:   assetStore,
		ReportsStore: reportsStore,
		Captain:      captain,
		Integrations: integrations,
		K8s:          k8s.NewClient(spaceMrn, integrations),
	}
}

func (s *Space) ListAssetsWithScores(ctx context.Context) ([]AssetWithScore, error) {
	assetsPage, err := s.AssetStore.ListAssets(ctx, &policy.AssetSearchFilter{SpaceMrn: s.spaceMrn})
	if err != nil {
		return nil, err
	}

	mrns := make([]string, len(assetsPage.List))
	for i := range assetsPage.List {
		mrns[i] = assetsPage.List[i].Mrn
	}

	scores, err := s.ReportsStore.ListAssetScores(ctx, &policy.ListAssetScoresReq{AssetMrns: mrns})
	if err != nil {
		return nil, err
	}

	assetScores := make([]AssetWithScore, len(assetsPage.List))
	for i := range assetsPage.List {
		asset := assetsPage.List[i]
		assetScores[i] = AssetWithScore{Asset: asset, Score: scores.Scores[asset.Mrn]}
	}
	return assetScores, nil
}