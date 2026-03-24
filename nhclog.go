package main

import (
    "errors"
    "os"
    "strings"
)

var activeFailedCheck string = ""
var failedCheckReason string = "All checks passed"
var failedCheckReran bool = false

func actOnLine(line string, m *metrics) {
    var err error
    hostname, _ := os.Hostname()

    if strings.HasPrefix(line, "Node Health Check starting.") {
        m.nhcRunTotal.WithLabelValues(hostname).Inc()
    } else if strings.HasPrefix(line, "Running check:") && activeFailedCheck != "" {
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
            m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(1)
            m.nhcFailureTotal.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Inc()
        }
        // error handle here
        
    } else if strings.Contains(line, "Node Health Check completed successfully") {
        // clear all the errors
        if activeFailedCheck != "" {
            m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(0)
        }
        activeFailedCheck = ""
        failedCheckReason = "All checks passed"
        m.nhcNodeState.WithLabelValues(hostname, activeFailedCheck, failedCheckReason).Set(1)

    }

}

func parseErrorLine(line string) (string, string, error) {
    parsedLine := strings.Split(line, ":")
    if len(parsedLine) != 5 {
        return "", "", errors.New("unable to parse ERROR line")
    }

    check := strings.TrimSpace(parsedLine[3])
    reason := strings.TrimSpace(parsedLine[4])
    return check, reason, nil
}
