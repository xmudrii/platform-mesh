package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"platform-mesh.io/account-operator/pkg/subroutines/util"
	corev1alpha1 "platform-mesh.io/apis/core/v1alpha1"
)

func TestGetWorkspaceTypeName(t *testing.T) {
	got := util.GetWorkspaceTypeName("test", corev1alpha1.AccountTypeOrg)
	assert.Equal(t, "test-org", got)
}
