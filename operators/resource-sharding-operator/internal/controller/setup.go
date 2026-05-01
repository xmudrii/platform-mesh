package controller

import (
	"github.com/platform-mesh/golang-commons/logger"
	"k8s.io/client-go/discovery"
	ctrl "sigs.k8s.io/controller-runtime"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
)

func SetupWithManager(mgr mcmanager.Manager, log *logger.Logger) error {
	localMgr := mgr.GetLocalManager()

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(localMgr.GetConfig())
	if err != nil {
		return err
	}

	registry := NewDynamicControllerRegistry()

	reconciler := &ResourceShardingReconciler{
		Client:    localMgr.GetClient(),
		Discovery: discoveryClient,
		Registry:  registry,
		Manager:   localMgr,
	}

	if err := reconciler.SetupWithManager(localMgr); err != nil {
		return err
	}

	_ = log
	_ = ctrl.Log
	return nil
}
