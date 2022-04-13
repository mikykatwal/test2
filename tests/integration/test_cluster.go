package integration

import (
	"testing"

	"go.mondoo.com/mondoo-operator/tests/framework/installer"
	"go.uber.org/zap"
)

type TestCluster struct {
	*installer.MondooInstaller
	T func() *testing.T
}

func StartTestCluster(settings installer.Settings, t func() *testing.T) *TestCluster {
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)

	cluster := &TestCluster{
		MondooInstaller: installer.NewMondooInstaller(settings, t),
		T:               t,
	}

	if err := cluster.InstallOperator(); err != nil {
		zap.S().Errorf("Mondoo operator was not installed successfully: %v", err)
		if !t().Failed() {
			cluster.GatherAllMondooLogs(t().Name(), installer.MondooNamespace)
		}
		t().Fail()
		if err := cluster.UninstallOperator(); err != nil {
			zap.S().Errorf("Failed to uninstall Mondoo operator. %v", err)
		}
		t().FailNow()
	}

	return cluster
}

func HandlePanics(r interface{}, uninstaller func(), t func() *testing.T) {
	if r != nil {
		zap.S().Infof("unexpected panic occurred during test %s, --> %v", t().Name(), r)
		t().Fail()
		uninstaller()
		t().FailNow()
	}
}
