package spread

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SpreadReconcileStatus is implemented by objects that support reconciliation spreading.
type SpreadReconcileStatus interface {
	GetObservedGeneration() int64
	SetObservedGeneration(int64)
	GetNextReconcileTime() metav1.Time
	SetNextReconcileTime(metav1.Time)
}
