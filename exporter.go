package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func nhcexport(ctx context.Context, reg *prometheus.Registry, httpPort int) error {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", httpPort),
        Handler: mux,
    }

    errChannel := make(chan error, 1)
    go func() {
        errChannel <- server.ListenAndServe()
    }()

    select {
    case err := <- errChannel:
        return err
    case <-ctx.Done():
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        return server.Shutdown(shutdownCtx)
    }
}
