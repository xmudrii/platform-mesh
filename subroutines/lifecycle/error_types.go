package lifecycle

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Action represents the type of operation a subroutine performs.
type Action string

const (
	ActionProcess    Action = "process"
	ActionFinalize   Action = "finalize"
	ActionInitialize Action = "initialize"
	ActionTerminate  Action = "terminate"
)

// String returns the string representation of the action.
func (a Action) String() string {
	return string(a)
}

// IsFinalize returns true if the action is a finalize or terminate operation.
func (a Action) IsFinalize() bool {
	return a == ActionFinalize || a == ActionTerminate
}

// ErrorReporter is called when a subroutine returns an error.
type ErrorReporter interface {
	Report(ctx context.Context, err error, info ErrorInfo)
}

// ErrorInfo provides context about the subroutine that failed.
type ErrorInfo struct {
	Subroutine string
	Object     client.Object
	Action     Action
}
