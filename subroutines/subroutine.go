package subroutines

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Subroutine is the base interface that all subroutines must implement.
type Subroutine interface {
	GetName() string
}

// Processor handles the main reconciliation logic for a subroutine.
type Processor interface {
	Subroutine
	Process(ctx context.Context, obj client.Object) (Result, error)
}

// Finalizer handles cleanup when an object is being deleted.
type Finalizer interface {
	Subroutine
	Finalize(ctx context.Context, obj client.Object) (Result, error)
	Finalizers(obj client.Object) []string
}

// Initializer handles one-time initialization when an initializer marker is present in status.
type Initializer interface {
	Subroutine
	Initialize(ctx context.Context, obj client.Object) (Result, error)
}

// Terminator handles ordered teardown when a terminator marker is present in status.
type Terminator interface {
	Subroutine
	Terminate(ctx context.Context, obj client.Object) (Result, error)
}
