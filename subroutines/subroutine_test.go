package subroutines

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Compile-time interface assertions.
type testSubroutine struct{}

func (t *testSubroutine) GetName() string                                           { return "test" }
func (t *testSubroutine) Process(context.Context, client.Object) (Result, error)    { return OK(), nil }
func (t *testSubroutine) Finalize(context.Context, client.Object) (Result, error)   { return OK(), nil }
func (t *testSubroutine) Finalizers(client.Object) []string                         { return nil }
func (t *testSubroutine) Initialize(context.Context, client.Object) (Result, error) { return OK(), nil }
func (t *testSubroutine) Terminate(context.Context, client.Object) (Result, error)  { return OK(), nil }

var (
	_ Subroutine  = &testSubroutine{}
	_ Processor   = &testSubroutine{}
	_ Finalizer   = &testSubroutine{}
	_ Initializer = &testSubroutine{}
	_ Terminator  = &testSubroutine{}
)
