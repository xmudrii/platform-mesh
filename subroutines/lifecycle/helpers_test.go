package lifecycle

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
)

// fakeManager implements mcmanager.Manager for tests by embedding the interface
// and only overriding ClusterFromContext.
type fakeManager struct {
	mcmanager.Manager
	cl client.Client
}

func (f *fakeManager) ClusterFromContext(context.Context) (cluster.Cluster, error) {
	return &fakeCluster{cl: f.cl}, nil
}

type fakeCluster struct {
	cluster.Cluster
	cl client.Client
}

func (f *fakeCluster) GetClient() client.Client { return f.cl }

type fakeManagerWithError struct {
	mcmanager.Manager
	err error
}

func (f *fakeManagerWithError) ClusterFromContext(context.Context) (cluster.Cluster, error) {
	return nil, f.err
}
