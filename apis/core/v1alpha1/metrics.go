package v1alpha1

import (
	"github.com/prometheus/client_golang/prometheus"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

// WebhookValidations counts Account admission webhook validation requests. The
// webhook lives in this shared apis module, so its metric is defined here rather
// than in the account-operator (which would create an apis -> operator cycle).
var WebhookValidations = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "account_operator_webhook_validations_total",
		Help: "Total number of Account webhook validation requests.",
	},
	[]string{"operation", "result", "account_type"},
)

func init() {
	ctrlmetrics.Registry.MustRegister(WebhookValidations)
}
