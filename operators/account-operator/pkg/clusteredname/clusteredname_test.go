package clusteredname_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"platform-mesh.io/account-operator/pkg/clusteredname"
	corev1alpha1 "platform-mesh.io/apis/core/v1alpha1"
)

func TestGetClusteredName_NoClusterInContext(t *testing.T) {
	obj := &corev1alpha1.Account{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a",
		},
	}

	cn, ok := clusteredname.GetClusteredName(t.Context(), obj)

	require.False(t, ok)
	require.Equal(t, "a", cn.Name)
}
