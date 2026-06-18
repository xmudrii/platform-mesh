package subroutines

import (
	"context"
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type contextKey struct{}

// WithClient stores a client.Client in the context.
func WithClient(ctx context.Context, cl client.Client) context.Context {
	return context.WithValue(ctx, contextKey{}, cl)
}

// ClientFromContext retrieves the client.Client from the context.
func ClientFromContext(ctx context.Context) (client.Client, error) {
	cl, ok := ctx.Value(contextKey{}).(client.Client)
	if !ok || cl == nil {
		return nil, errors.New("no client in context")
	}
	return cl, nil
}

// MustClientFromContext retrieves the client.Client from the context.
// It panics if no client is stored in the context.
func MustClientFromContext(ctx context.Context) client.Client {
	cl, err := ClientFromContext(ctx)
	if err != nil {
		panic(err)
	}
	return cl
}
