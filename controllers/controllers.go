package controllers

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	pcontext "github.com/ricoleabricot/rudder/context"
)

var (
	reconciliationCountSuccess = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rudder",
		Subsystem: "reconciliation",
		Name:      "success",
	}, []string{"namespace", "name"})

	reconciliationCountError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rudder",
		Subsystem: "reconciliation",
		Name:      "errors",
	}, []string{"namespace", "name"})

	reconciliationDuration = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rudder",
		Subsystem: "reconciliation",
		Name:      "duration",
	}, []string{"namespace", "name"})
)

func init() {
	metrics.Registry.MustRegister(
		reconciliationCountSuccess,
		reconciliationCountError,
		reconciliationDuration,
	)
}

// Reconciler defines the interface which all reconciler should implement
type Reconciler struct {
	client.Client

	Scheme *runtime.Scheme
	Logger *logrus.Logger

	Env string
}

// InitContext set all clients and configurations into the context
func (r *Reconciler) InitContext(ctx context.Context, req ctrl.Request) context.Context {
	// Forge background context and logger
	log := r.Logger.
		WithField("namespace", req.Namespace).
		WithField("name", req.Name)

	// Set all useful information to context
	ctx = pcontext.WithKubernetesClient(ctx, r.Client)
	ctx = pcontext.WithLogger(ctx, log)
	ctx = pcontext.WithRequest(ctx, req)

	return ctx
}

// AddReconciliationCount increase metrics counter for reconciliations
func (r *Reconciler) AddReconciliationCount(namespace, name string, start *time.Time) {
	reconciliationCountSuccess.WithLabelValues(namespace, name).Inc()
	if start != nil {
		duration := time.Since(*start)
		reconciliationDuration.WithLabelValues(namespace, name).Add(float64(duration.Milliseconds()))
	}
}

// AddReconciliationError increase metrics counter for reconciliation errors
func (r *Reconciler) AddReconciliationError(namespace, name string) {
	reconciliationCountError.WithLabelValues(namespace, name).Inc()
}
