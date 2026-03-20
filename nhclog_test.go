package main

import (
    "errors"
    "testing"
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

func prettyPrint(val string) string {
    if val == "" {
        return "None"
    }
    return val
}
