/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"

	pcontext "github.com/ricoleabricot/rudder/context"
	"github.com/ricoleabricot/rudder/kubernetes"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	*Reconciler
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()

	// Init the current reconcile context
	ctx = r.InitContext(ctx, req)
	logger := pcontext.Logger(ctx)

	// Get the deploy resource from request
	pod, err := kubernetes.GetPod(ctx, req.Namespace, req.Name)
	if err != nil {
		// If the resource has already been deleted, stop watching it
		if kubernetes.IsNotFoundErr(err) {
			logger.Info("The resource does not exist anymore, do not watch it anymore")
			return ctrl.Result{}, nil
		}

		logger.WithError(err).Error("Failed to watch resource")
		r.AddReconciliationError(req.Namespace, req.Name)

		return ctrl.Result{}, fmt.Errorf("failed to get stop watching resource: %w", err)
	}

	ctx = pcontext.WithPod(ctx, pod)

	// TODO: implement logic here.

	r.AddReconciliationCount(pod.Namespace, pod.Name, &start)

	return ctrl.Result{}, nil
}
