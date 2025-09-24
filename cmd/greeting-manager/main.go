package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"

	"github.com/marcosartorato/basic-controller/internal/controller"

	metricsServer "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	var metricsAddr, probeAddr string
	var leaderElect bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "metrics addr")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "health addr")
	flag.BoolVar(&leaderElect, "leader-elect", false, "enable leader election")
	flag.Parse()

	logger := zap.New(zap.UseDevMode(true))
	ctrl.SetLogger(logger)

	var scheme = runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsServer.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         leaderElect,
		LeaderElectionID:       "greeting-operator",
	})
	if err != nil {
		logger.Error(err, "manager:")
		os.Exit(1)
	}

	if err := controller.SetupGreeting(mgr); err != nil {
		fmt.Fprintln(os.Stderr, "controller:", err)
		os.Exit(1)
	}

	_ = mgr.AddHealthzCheck("healthz", func(_ *http.Request) error { return nil })
	_ = mgr.AddReadyzCheck("readyz", func(_ *http.Request) error { return nil })

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Fprintln(os.Stderr, "start:", err)
		os.Exit(1)
	}
}
