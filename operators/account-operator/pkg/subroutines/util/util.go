package util

import (
	"fmt"

	corev1alpha1 "platform-mesh.io/apis/core/v1alpha1"
)

func GetWorkspaceTypeName(accountName string, accountType corev1alpha1.AccountType) string {
	return fmt.Sprintf("%s-%s", accountName, accountType)
}
