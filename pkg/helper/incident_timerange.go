package helper

import (
	"time"

	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
)

// OverlapsTimeRange() returns true unless it is definitely outside of
// the time range.
func OverlapsTimeRange(incident *api.Incident, start, end time.Time) bool {
	if incident.BeganAt != nil {
		if incident.BeganAt.After(end) {
			return false
		}
	}
	if incident.EndedAt != nil {
		if incident.EndedAt.Before(start) {
			return false
		}
	}
	return true
}
