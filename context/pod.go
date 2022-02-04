package context

import (
	"context"

	"k8s.io/api/core/v1"
)

var podContext = "_context_pod"

// DaemonSet returns a kubernetes pod resource from the context
func Pod(ctx context.Context) *v1.Pod {
	return ctx.Value(&podContext).(*v1.Pod)
}

// WithPod insert a kubernetes pod resource into the context
func WithPod(ctx context.Context, pod *v1.Pod) context.Context {
	return context.WithValue(ctx, &podContext, pod)
}
