package helper

import (
	"time"

	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
)

func IsWithinTimeRange(incident *api.Incident, start, end time.Time) bool {
	return true
}
