package main

import (
    "flag"

    "github.com/prometheus/client_golang/prometheus"
)

func main() {
    var httpPort = flag.Int("http-port", 8090, "port for the webserver to listen on.")
    var logPath = flag.String("log-path", "/var/log/nhc.log", "NHC log path")
    var readTime = flag.Int("read-time", 5, "default time in minutes between log reads.")
    flag.Parse()

    reg := prometheus.NewRegistry()
    m := newMetrics(reg)
    recordNHC(m, logPath, readTime)

    nhcexport(reg, httpPort)
}
