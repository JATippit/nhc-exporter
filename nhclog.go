package main

import (
    "errors"
    "os"
    "strings"
)

//activeFailedCheck := ""
//failedCheckReran := false 
/*
func actOnLine(line string, m *metrics) {
    if strings.HasPrefix(line, "Node Health Check starting.") {
        hostname, _ := os.Hostname()
        m.nhcRunTotal.WithLabelValues(hostname).Inc()
    } else if strings.HasPrefix(line, "Running check:") && activeFailedCheck != "" {
        if !strings.Contains(line, activeFailedCheck) && failedCheckReran {
            // clear the error
        } else if strings.Contains(line, activeFailedCheck) {
            failedCheckReran = true
        }
    } else if strings.HasPrefix(line, "ERROR:") && !strings.Contains(line, activeFailedCheck) {
        // update the error
    } else if strings.Contains("Node Health Check completed successfully") {
        // clear all the errors
    }

}
*/

func parseErrorLine(line string) (string, string, error) {
    parsedLine := strings.Split(line, ":")
    if len(parsedLine) != 5 {
        return "", "", errors.New("unable to parse ERROR line")
    }

    check := strings.TrimSpace(parsedLine[3])
    reason := strings.TrimSpace(parsedLine[4])
    return check, reason, nil
}
