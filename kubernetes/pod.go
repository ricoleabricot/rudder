package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	pcontext "github.com/ricoleabricot/rudder/context"
)

// GetPod fetch a pod resource in a Kubernetes namespace
func GetPod(ctx context.Context, namespace string, name string) (*v1.Pod, error) {
	client := pcontext.KubernetesClient(ctx)
	pod := &v1.Pod{}

	err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pod)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod resource: %w", err)
	}

	// If UID is not set, returns an empty pointer
	if pod.UID == "" {
		return nil, nil
	}

	return pod, nil
}
