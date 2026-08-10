package main

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
    nhcNodeState *prometheus.GaugeVec
    nhcRunTotal *prometheus.CounterVec
    nhcFailureTotal *prometheus.CounterVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
    m := &metrics{
        nhcNodeState: promauto.With(reg).NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "nhc_node_state",
                Help: "NHC node state: 1 indicates active state, 0 indicates inactive",
            },
            []string{"node", "check", "reason"},
        ),
        nhcRunTotal: promauto.With(reg).NewCounterVec(
            prometheus.CounterOpts{
                Name: "nhc_run_total",
                Help: "Number of times NHC has run",
            },
            []string{"node"},
        ),
        nhcFailureTotal: promauto.With(reg).NewCounterVec(
            prometheus.CounterOpts{
                Name: "nhc_failure_total",
                Help: "Per check failure totals",
            },
            []string{"node", "check", "reason"},
        ),
    }

    return m
}

