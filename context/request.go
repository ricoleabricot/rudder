package context

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

var requestContext = "_context_request"

// Request returns the current operator request from the context
func Request(ctx context.Context) ctrl.Request {
	return ctx.Value(&requestContext).(ctrl.Request)
}

// WithRequest inserts the current operator request into the context
func WithRequest(ctx context.Context, request ctrl.Request) context.Context {
	return context.WithValue(ctx, &requestContext, request)
}
