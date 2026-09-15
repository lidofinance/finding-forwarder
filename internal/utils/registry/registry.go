package registry

import (
	"slices"
	"strings"

	"github.com/lidofinance/onchain-mon/generated/databus"
)

type FindingMapping = map[databus.Severity]bool
type FindingFilterMap = map[string]bool

type NotificationChannel string

const (
	Telegram NotificationChannel = `Telegram`
	Discord  NotificationChannel = `Discord`
	OpsGenie NotificationChannel = `OpsGenie`
	Slack    NotificationChannel = `Slack`
)

var CanonicalSeverities = []databus.Severity{
	databus.SeverityUnknown,
	databus.SeverityInfo,
	databus.SeverityLow,
	databus.SeverityMedium,
	databus.SeverityHigh,
	databus.SeverityCritical,
}

func IsCanonicalSeverity(severity databus.Severity) bool {
	return slices.Contains(CanonicalSeverities, severity)
}

func CanonicalSeverityList() string {
	out := make([]string, 0, len(CanonicalSeverities))
	for _, severity := range CanonicalSeverities {
		out = append(out, string(severity))
	}

	return strings.Join(out, ", ")
}
