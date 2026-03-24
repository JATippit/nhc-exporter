package main

import (
    "errors"
    "strings"
    "testing"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/testutil"
)

func TestParseErrorLine(t *testing.T) {
    lines := [4]string{
        "ERROR:  nhc:  Health check failed:  check_ps_service:  Service chronyd (process chronyd) owned by chrony not running; start in progress",
        "ERROR:  nhc:  Health check failed:  check_gpu_count:  Invalid number of AMD GPUs present.",
        "ERROR:  nhc:  Health check failed:  some reason",
        "ERROR:  nhc:  Health check failed:  check_fake_test:  some reason: some other reason",
    }

    expected := []struct {
        check string
        reason string
        err error
    }{
        {"check_ps_service", "Service chronyd (process chronyd) owned by chrony not running; start in progress", nil},
        {"check_gpu_count", "Invalid number of AMD GPUs present.", nil},
        {"", "", errors.New("unable to parse ERROR line")},
        {"", "", errors.New("unable to parse ERROR line")},
    }

    for i, line := range lines {
        check, reason, err := parseErrorLine(line)

        // verify error conditions
        if (expected[i].err != nil && err == nil) || (expected[i].err == nil && err != nil) {
            t.Errorf("case %d) expected error: %+v, got: %+v", i, expected[i].err, err)
        } else if expected[i].err != nil && err.Error() != expected[i].err.Error() {
            t.Errorf("case %d) expected error: %+v, got: %+v", i, expected[i].err, err)
        }

        // verify check/reason
        if check != expected[i].check || reason != expected[i].reason {
            t.Errorf("case %d) expected failed check: %s with reason %s, got: failed check: %s with reason %s",
                i,
                prettyPrint(expected[i].check),
                prettyPrint(expected[i].reason),
                prettyPrint(check),
                prettyPrint(reason))
        }
    }
}

func TestActOnErrorLine(t *testing.T) {
    reg := prometheus.NewRegistry()
    m := newMetrics(reg)

    lines := []string{
        "Node Health Check starting.",
        "Running check:  \"check_hw_mem_free 1mb\"",
        "Running check:  \"check_fs_mount_rw -f /\"",
        "Running check:  \"check_file_test -r -w -x -d -k /tmp /var/tmp\"",
        "Running check:  \"check_ps_service -u chrony -S chronyd\"",
        "Running check:  \"check_file_test -r -s /etc/passwd /etc/group\"",
        "Running check:  \"check_file_test -c -r -w /dev/null /dev/zero\"",
        "Running check:  \"check_gpu_count 3\"",
        "ERROR:  nhc:  Health check failed:  check_gpu_count:  Invalid number of AMD GPUs present.",
    }
    for _, line := range lines {
        _ = actOnLine(line, m)
    }

    expectedState := `
# HELP nhc_node_state NHC node state: 1 indicates active state, 0 indicates inactive
# TYPE nhc_node_state gauge
nhc_node_state{check="--none--",node="nhc-test",reason="All checks passed"} 0
nhc_node_state{check="check_gpu_count",node="nhc-test",reason="Invalid number of AMD GPUs present."} 1
`

    expectedErrCount := `
# HELP nhc_failure_total Per check failure totals
# TYPE nhc_failure_total counter
nhc_failure_total{check="check_gpu_count",node="nhc-test",reason="Invalid number of AMD GPUs present."} 1
`

    expectedNHCCount := `
# HELP nhc_run_total Number of times NHC has run
# TYPE nhc_run_total counter
nhc_run_total{node="nhc-test"} 1
`

    if err := testutil.CollectAndCompare(m.nhcNodeState, strings.NewReader(expectedState)); err != nil {
        t.Errorf("%s", err)
    }
    if err := testutil.CollectAndCompare(m.nhcRunTotal, strings.NewReader(expectedNHCCount)); err != nil {
        t.Errorf("%s", err)
    }
    if err := testutil.CollectAndCompare(m.nhcFailureTotal, strings.NewReader(expectedErrCount)); err != nil {
        t.Errorf("%s", err)
    }

}

func prettyPrint(val string) string {
    if val == "" {
        return "None"
    }
    return val
}
