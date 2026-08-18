package main

import (
    "bufio"
    "context"
    "flag"
    "log"
    "os"

    "github.com/prometheus/client_golang/prometheus"
)

func main() {
    var httpPort = flag.Int("http-port", 8090, "port for the webserver to listen on.")
    var logPath = flag.String("log-path", "/var/log/nhc.log", "NHC log path")
    var readTime = flag.Int("read-time", 5, "default time in minutes between log reads.")
    flag.Parse()

    hostname, err := os.Hostname()
    if err != nil {
        log.Printf("err: %v\n", err)
        os.Exit(1)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    nhcLog, err := os.Open(*logPath)
    if err != nil {
        log.Printf("error: %v\n", err)
        os.Exit(1)
    }
    defer nhcLog.Close()
    r := bufio.NewReader(nhcLog)
    
    reg := prometheus.NewRegistry()
    m := newMetrics(reg)
    recordNHC(ctx, hostname, m, r, readTime)

    nhcexport(ctx, reg, httpPort)
}
