// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resource_monitor

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mondoo.com/mondoo-operator/controllers/resource_monitor/debouncer/mock"
	"go.mondoo.com/mondoo-operator/controllers/resource_monitor/scan_api_store"
	scanapistoremock "go.mondoo.com/mondoo-operator/controllers/resource_monitor/scan_api_store/mock"
	"go.mondoo.com/mondoo-operator/tests/framework/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type ResourceMonitorControllerSuite struct {
	suite.Suite
	mockCtrl          *gomock.Controller
	debouncerMock     *mock.MockDebouncer
	fakeClientBuilder *fake.ClientBuilder
}

func (s *ResourceMonitorControllerSuite) BeforeTest(suiteName, testName string) {
	s.mockCtrl = gomock.NewController(s.T())
	s.debouncerMock = mock.NewMockDebouncer(s.mockCtrl)
	s.fakeClientBuilder = fake.NewClientBuilder().WithObjects(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}})
}

func (s *ResourceMonitorControllerSuite) AfterTest(suiteName, testName string) {
	s.mockCtrl.Finish()
}

func (s *ResourceMonitorControllerSuite) TestReconcile_Pod() {
	ctx := context.Background()
	scanApiStore := scanapistoremock.NewMockScanApiStore(s.mockCtrl)
	r, err := NewResourceMonitorController(
		s.fakeClientBuilder.Build(),
		func() client.Object { return &corev1.Pod{} },
		scanApiStore)
	s.Require().NoError(err)
	r.debouncer = s.debouncerMock

	ns := utils.RandString(10)
	name := utils.RandString(10)
	scanApiStore.EXPECT().GetAll().Return([]scan_api_store.ClientConfiguration{{}}).Times(1)
	s.debouncerMock.EXPECT().Add(fmt.Sprintf("pod:%s:%s", ns, name)).Times(1)

	res, err := r.Reconcile(ctx, controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      name,
		},
	})
	s.True(res.IsZero())
	s.NoError(err)
}

func (s *ResourceMonitorControllerSuite) TestReconcile_Pod_NoScanApi() {
	ctx := context.Background()
	scanApiStore := scanapistoremock.NewMockScanApiStore(s.mockCtrl)
	r, err := NewResourceMonitorController(
		s.fakeClientBuilder.Build(),
		func() client.Object { return &corev1.Pod{} },
		scanApiStore)
	s.Require().NoError(err)
	r.debouncer = s.debouncerMock

	ns := utils.RandString(10)
	name := utils.RandString(10)
	scanApiStore.EXPECT().GetAll().Return(nil)

	res, err := r.Reconcile(ctx, controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      name,
		},
	})
	s.True(res.IsZero())
	s.NoError(err)
}

func TestResourceMonitorControllerSuite(t *testing.T) {
	suite.Run(t, new(ResourceMonitorControllerSuite))
}
