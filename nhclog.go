package main

import (
    "bufio"
    "context"
    "errors"
    "io"
    "log"
    "os"
    "strings"
    "time"
)

const (
    noneFailed = "--none--"
    passReason = "All checks passed"
)

var activeFailedCheck string = noneFailed
var failedCheckReason string = passReason
var failedCheckReran bool = false

func actOnLine(line string, m *metrics) error {
    var err error
    hostname, _ := os.Hostname()

    if strings.HasPrefix(line, "Node Health Check starting.") {
        m.nhcRunTotal.WithLabelValues(hostname).Inc()
    } else if strings.HasPrefix(line, "Running check:") && activeFailedCheck != noneFailed {
        if !strings.Contains(line, activeFailedCheck) && failedCheckReran {
            // clear the error
            m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(0)
        } else if strings.Contains(line, activeFailedCheck) {
            failedCheckReran = true
        }

    } else if strings.HasPrefix(line, "ERROR:") && !strings.Contains(line, activeFailedCheck) {
        // update the error
        m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(0)
        activeFailedCheck, failedCheckReason, err = parseErrorLine(line)
        if err != nil {
            return err
        }
        m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(1)
        m.nhcFailureTotal.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Inc()
        
    } else if strings.Contains(line, "Node Health Check completed successfully") {
        // clear all the errors
        if activeFailedCheck != noneFailed {
            m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(0)
        }
        activeFailedCheck = noneFailed
        failedCheckReason = passReason
        m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(1)

    }
    return nil
}

func parseErrorLine(line string) (string, string, error) {
    parsedLine := strings.SplitN(line, ":", 5)
    if len(parsedLine) != 5 {
        return "", "", errors.New("unable to parse ERROR line")
    }

    check := strings.TrimSpace(parsedLine[3])
    reason := strings.TrimSpace(parsedLine[4])
    return check, reason, nil
}

func recordNHC(ctx context.Context, m *metrics, r *bufio.Reader, readTime *int) {
    go func () {
        var partialLine string
        var fullLine string
        hitEOF := false

        for {
            line, err := r.ReadBytes('\n')

            select {
            case <-ctx.Done():
                return
            default:
            }

            switch err {
            case nil:
                fullLine = string(line)
                if hitEOF {
                    fullLine = partialLine + fullLine
                    partialLine = ""
                    hitEOF = false
                }
                actOnLine(fullLine, m)

            case io.EOF:
                hitEOF = true
                if len(string(line)) > 0 {
                    partialLine = partialLine + string(line)
                }
                time.Sleep(time.Duration(*readTime) * time.Second)

            default:
                log.Printf("unexpected error: %v", err)
            }
        }
    }()
}
