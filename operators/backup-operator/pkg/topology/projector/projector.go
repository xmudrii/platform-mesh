package projector

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/platform-mesh/backup-operator/pkg/topology"
)

const configMapName = "backup-topology-schemas"

// Projector ensures the topology schema ConfigMap is present and up-to-date.
type Projector struct {
	client    client.Client
	namespace string
}

// New returns a Projector that manages the schema ConfigMap in namespace.
func New(c client.Client, namespace string) *Projector {
	return &Projector{client: c, namespace: namespace}
}

// EnsureConfigMap creates or updates the backup-topology-schemas ConfigMap with
// the current schema(s) keyed by schemaVersion. Idempotent via server-side apply.
func (p *Projector) EnsureConfigMap(ctx context.Context) error {
	schemaData, err := topology.SchemaV1Alpha1()
	if err != nil {
		return fmt.Errorf("reading topology schema: %w", err)
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: p.namespace,
		},
		Data: map[string]string{
			"v1alpha1.json": string(schemaData),
		},
	}

	if err := p.client.Patch(ctx, cm, client.Apply, client.FieldOwner("backup-operator"), client.ForceOwnership); err != nil {
		return fmt.Errorf("applying topology schema ConfigMap: %w", err)
	}
	return nil
}
