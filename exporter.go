package main

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func nhcexport(reg *prometheus.Registry, httpPort *int) {
    listenPort := fmt.Sprintf(":%d", *httpPort)

    http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
    http.ListenAndServe(listenPort, nil)
}
