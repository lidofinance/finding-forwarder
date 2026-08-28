package notifiler_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lidofinance/onchain-mon/generated/databus"
	"github.com/lidofinance/onchain-mon/internal/connectors/metrics"
	"github.com/lidofinance/onchain-mon/internal/pkg/notifiler"
	"github.com/lidofinance/onchain-mon/internal/utils/registry"
)

func newTestMetrics(t *testing.T) *metrics.Store {
	t.Helper()
	reg := prometheus.NewRegistry()
	return metrics.New(reg, "test", "test", "test")
}

// A nil error used to mean "sent" to the consumer, which then acked the finding
// and set a cooldown for a severity OpsGenie never received. ValidateConfig now
// rejects such a consumer, so this path means the invariant is broken and must
// not look like a successful delivery.
func TestSendFinding_RejectsUnsupportedSeverity(t *testing.T) {
	m := newTestMetrics(t)
	og := notifiler.NewOpsgenie("key", nil, m, "local", "etherscan.io", "test")

	for _, sev := range []databus.Severity{databus.SeverityInfo, databus.SeverityLow, databus.SeverityMedium, databus.SeverityUnknown} {
		alert := &databus.FindingDtoJson{
			Name:        "test",
			Description: "desc",
			Severity:    sev,
			AlertId:     "TEST",
			BotName:     "bot",
			Team:        "team",
			UniqueKey:   "key",
		}
		// httpClient is nil, so a returned error also proves no HTTP call was
		// attempted — a send would have panicked instead.
		err := og.SendFinding(context.Background(), alert)
		if err == nil {
			t.Fatalf("SendFinding(%s) returned nil, want an error", sev)
		}
		if !strings.Contains(err.Error(), "cannot deliver severity") {
			t.Errorf("SendFinding(%s) error is %q, want it to mention the severity", sev, err)
		}
		// The consumer terminates on this type instead of retrying it.
		if _, ok := errors.AsType[*notifiler.UndeliverableError](err); !ok {
			t.Errorf("SendFinding(%s) error is %T, want *notifiler.UndeliverableError", sev, err)
		}
		if !errors.Is(err, notifiler.ErrUndeliverable) {
			t.Errorf("SendFinding(%s) error does not match ErrUndeliverable", sev)
		}
	}
}

func TestOpsGeniePriority_CoversOnlyPageableSeverities(t *testing.T) {
	want := map[databus.Severity]string{
		databus.SeverityCritical: "P1",
		databus.SeverityHigh:     "P2",
		databus.SeverityMedium:   "",
		databus.SeverityLow:      "",
		databus.SeverityInfo:     "",
		databus.SeverityUnknown:  "",
	}

	for _, severity := range registry.CanonicalSeverities {
		if got := notifiler.OpsGeniePriority(severity); got != want[severity] {
			t.Errorf("OpsGeniePriority(%s) = %q, want %q", severity, got, want[severity])
		}
	}

	if got := notifiler.OpsGenieSeverityList(); got != "High, Critical" {
		t.Errorf("OpsGenieSeverityList() = %q, want %q", got, "High, Critical")
	}
}

func TestAlertPayload_DetailsContainsForwarderAttributes(t *testing.T) {
	payload := notifiler.AlertPayload{
		Message:  "test alert",
		Priority: "P1",
		Alias:    "mainnet-TEST",
		Details: map[string]string{
			"env":     "mainnet",
			"source":  "cluster1",
			"team":    "probe",
			"botName": "unusual-activity",
			"alertId": "UNUSUAL-ACTIVITY-LOW-BALANCE",
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	details, ok := decoded["details"].(map[string]any)
	if !ok {
		t.Fatal("details field missing or wrong type in serialized payload")
	}

	expected := map[string]string{
		"env":     "mainnet",
		"source":  "cluster1",
		"team":    "probe",
		"botName": "unusual-activity",
		"alertId": "UNUSUAL-ACTIVITY-LOW-BALANCE",
	}
	for k, want := range expected {
		if got := details[k]; got != want {
			t.Errorf("details[%s] = %v, want %s", k, got, want)
		}
	}
}

func TestAlertPayload_NilDetailsOmitted(t *testing.T) {
	payload := notifiler.AlertPayload{
		Message:  "test alert",
		Priority: "P2",
		Alias:    "test-TEST",
		Details:  nil,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := decoded["details"]; ok {
		t.Error("details field should be omitted when nil")
	}
}
