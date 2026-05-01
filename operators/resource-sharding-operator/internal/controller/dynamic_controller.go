package controller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/platform-mesh/resource-sharding-operator/api/v1alpha1"
)

type AssignmentReconciler struct {
	Client   client.Client
	Assigner *ShardAssigner
	LabelKey string
	GVR      schema.GroupVersionResource
}

func (r *AssignmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	obj := &metav1.PartialObjectMetadata{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   r.GVR.Group,
		Version: r.GVR.Version,
		Kind:    r.GVR.Resource, // Will be resolved by the cache
	})

	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if _, exists := obj.Labels[r.LabelKey]; exists {
		return ctrl.Result{}, nil
	}

	shard := r.Assigner.Next()
	patch := client.MergeFrom(obj.DeepCopy())
	if obj.Labels == nil {
		obj.Labels = make(map[string]string)
	}
	obj.Labels[r.LabelKey] = shard

	if err := r.Client.Patch(ctx, obj, patch); err != nil {
		return ctrl.Result{}, err
	}

	logger.V(1).Info("assigned shard", "resource", req.NamespacedName, "shard", shard)
	return ctrl.Result{}, nil
}

func StartDynamicController(ctx context.Context, mgr ctrl.Manager, rs *v1alpha1.ResourceSharding, gvr schema.GroupVersionResource) (*RunningController, error) {
	labelKey := rs.Spec.ShardLabelKey
	if labelKey == "" {
		labelKey = "sharding.platform-mesh.io/shard"
	}

	selector, err := labels.Parse("!" + labelKey)
	if err != nil {
		return nil, fmt.Errorf("parsing label selector: %w", err)
	}

	gvk := schema.GroupVersionKind{
		Group:   gvr.Group,
		Version: gvr.Version,
		Kind:    gvr.Resource,
	}

	obj := &metav1.PartialObjectMetadata{}
	obj.SetGroupVersionKind(gvk)

	informerCache, err := cache.New(mgr.GetConfig(), cache.Options{
		Scheme: mgr.GetScheme(),
		ByObject: map[client.Object]cache.ByObject{
			obj: {Label: selector},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating cache: %w", err)
	}

	ctrlCtx, cancel := context.WithCancel(ctx)

	go func() {
		_ = informerCache.Start(ctrlCtx)
	}()

	if !informerCache.WaitForCacheSync(ctrlCtx) {
		cancel()
		return nil, fmt.Errorf("cache sync failed for %s", gvr.String())
	}

	assigner := NewShardAssigner(shardNames(rs.Spec.Shards))

	c, err := controller.NewUnmanaged("shard-assign-"+rs.Name, controller.Options{
		Reconciler: &AssignmentReconciler{
			Client:   mgr.GetClient(),
			Assigner: assigner,
			LabelKey: labelKey,
			GVR:      gvr,
		},
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating controller: %w", err)
	}

	if err := c.Watch(source.Kind(informerCache, obj, &handler.TypedEnqueueRequestForObject[*metav1.PartialObjectMetadata]{})); err != nil {
		cancel()
		return nil, fmt.Errorf("setting up watch: %w", err)
	}

	go func() {
		_ = c.Start(ctrlCtx)
	}()

	return &RunningController{
		Cancel:   cancel,
		GVR:      gvr,
		Assigner: assigner,
	}, nil
}
