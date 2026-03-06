package lifecycle

import (
	"time"

	"github.com/platform-mesh/subroutines"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConditionManager manages per-subroutine and aggregate conditions.
type ConditionManager interface {
	InitUnknownConditions(obj client.Object, subroutineNames []string)
	SetSubroutineCondition(obj client.Object, name string, result subroutines.Result, err error, isFinalize bool)
	SetReadyCondition(obj client.Object, reason string)
}

// SpreadManager manages reconciliation spreading.
type SpreadManager interface {
	ReconcileRequired(obj client.Object) bool
	RequeueDelay(obj client.Object) time.Duration
	SetNextReconcileTime(obj client.Object)
	UpdateObservedGeneration(obj client.Object)
	RemoveRefreshLabel(obj client.Object) bool
}
