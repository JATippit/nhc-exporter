package main

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
    nhcNodeState *prometheus.GaugeVec
    nhcCheckFailureTotal *prometheus.CounterVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
    m := &metrics{
        nhcNodeState: promauto.With(reg).NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "nhc_node_state",
                Help: "NHC state: 0 indicates a check failed, 1 indicates all checks passed.",
            },
            []string{"node", "check", "reason"},
        ),
        nhcCheckFailureTotal: promauto.With(reg).NewCounterVec(
            prometheus.CounterOpts{
                Name: "nhc_check_failure_total",
                Help: "Per check failure totals",
            },
            []string{"node", "check"},
        ),
    }

    return m
}

