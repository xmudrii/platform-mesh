package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	shardDistribution = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_sharding_distribution",
			Help: "Number of resources assigned to each shard",
		},
		[]string{"resourcesharding", "shard"},
	)

	shardImbalanceRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "resource_sharding_imbalance_ratio",
			Help: "Maximum deviation from ideal distribution (0.0 = perfect balance)",
		},
		[]string{"resourcesharding"},
	)

	assignmentsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resource_sharding_assignments_total",
			Help: "Total number of shard assignments made",
		},
		[]string{"resourcesharding", "shard"},
	)

	rebalanceMovesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resource_sharding_rebalance_moves_total",
			Help: "Total number of resources moved during rebalancing",
		},
		[]string{"resourcesharding"},
	)
)

func init() {
	metrics.Registry.MustRegister(
		shardDistribution,
		shardImbalanceRatio,
		assignmentsTotal,
		rebalanceMovesTotal,
	)
}
