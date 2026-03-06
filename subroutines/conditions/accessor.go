package conditions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionAccessor is implemented by objects that expose status conditions.
type ConditionAccessor interface {
	GetConditions() []metav1.Condition
	SetConditions([]metav1.Condition)
}
