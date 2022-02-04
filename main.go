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

package main

import (
	"os"

	"github.com/sirupsen/logrus"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ricoleabricot/rudder/config"
	"github.com/ricoleabricot/rudder/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	// Create the logger
	logger := logrus.StandardLogger()

	// Parse config from command-lines and env vars
	logger.Infof("Setting up the operator configuration")
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Error("Failed to set up the operator configuration")
		os.Exit(1)
	}

	// Setup everything
	logger.Infof("Setting up the operator manager and controllers")
	mgr := SetupManager(cfg, logger)
	SetupControllers(mgr, cfg, logger)

	// Run main process
	logger.Infof("Running the operator manager")
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		logger.WithError(err).Error("Failed to run the operator manager")
		os.Exit(1)
	}
}

func SetupManager(cfg *config.Config, logger *logrus.Logger) manager.Manager {
	// Fetch the kubeconfig from configuration or create in-cluster default configuration
	var err error
	var kubeconfig *rest.Config

	if cfg.Kubeconfig != "" {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch kubeconfig")
			os.Exit(1)
		}
	} else {
		kubeconfig = ctrl.GetConfigOrDie()
	}

	// Create the controller manager
	mgr, err := ctrl.NewManager(kubeconfig, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     cfg.MetricsAddr,
		HealthProbeBindAddress: cfg.ProbeAddr,
		LeaderElection:         cfg.LeaderElectionEnable,
		Port:                   cfg.Port,
		LeaderElectionID:       "c1a01e01.rico.io",
	})
	if err != nil {
		logger.WithError(err).Error("Failed to start the operator manager")
		os.Exit(1)
	}

	err = mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		logger.WithError(err).Error("Failed to start add health check")
		os.Exit(1)
	}

	err = mgr.AddHealthzCheck("readyz", healthz.Ping)
	if err != nil {
		logger.WithError(err).Error("Failed to start add readiness check")
		os.Exit(1)
	}

	return mgr
}

func SetupControllers(mgr manager.Manager, cfg *config.Config, logger *logrus.Logger) {
	// Define the default reconciler interface
	reconciler := &controllers.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Logger: logger,
		Env:    cfg.Env,
	}

	podCtrl := controllers.PodReconciler{Reconciler: reconciler}
	err := podCtrl.SetupWithManager(mgr)
	if err != nil {
		logger.WithError(err).Error("Failed to setup pod controller")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
}
