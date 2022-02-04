package context

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var kubeContext = "_context_kubernetes_client"

// KubernetesClient returns a kubernetes client from the context
func KubernetesClient(ctx context.Context) client.Client {
	return ctx.Value(&kubeContext).(client.Client)
}

// WithKubernetesClient inserts a kubernetes client in the context
func WithKubernetesClient(ctx context.Context, client client.Client) context.Context {
	return context.WithValue(ctx, &kubeContext, client)
}
