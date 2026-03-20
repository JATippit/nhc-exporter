package main

import (
    "flag"

    "github.com/prometheus/client_golang/prometheus"
)

func main() {
    var httpPort = flag.Int("http-port", 8090, "port for the webserver to listen on.")
    flag.Parse()

    reg := prometheus.NewRegistry()
//    m := newMetrics(reg)

    nhcexport(reg, httpPort)
}
