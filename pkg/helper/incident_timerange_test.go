package helper

import (
	"testing"
	"time"

	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
)

func TestTimeRange(t *testing.T) {
	start := time.Date(2020, 1, 2, 3, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 2, 3, 0, 0, 0, time.UTC)

	timePast1 := time.Date(2018, 2, 2, 3, 0, 0, 0, time.UTC)
	timePast2 := time.Date(2019, 2, 2, 3, 0, 0, 0, time.UTC)
	timeWithin1 := time.Date(2020, 2, 2, 3, 0, 0, 0, time.UTC)
	timeWithin2 := time.Date(2021, 2, 2, 3, 0, 0, 0, time.UTC)
	timeLater1 := time.Date(2022, 2, 2, 3, 0, 0, 0, time.UTC)
	timeLater2 := time.Date(2024, 2, 2, 3, 0, 0, 0, time.UTC)

	testCases := []struct {
		Incident            api.Incident
		ExpectedToBeInRange bool
	}{
		{
			Incident:            api.Incident{},
			ExpectedToBeInRange: true,
		},
		{
			Incident: api.Incident{
				BeganAt: &timePast1, // Begins in the past, does not end
			},
			ExpectedToBeInRange: true,
		},
		{
			Incident: api.Incident{
				BeganAt: &timeLater2, // Begins in the future, does not end
			},
			ExpectedToBeInRange: false,
		},
		{
			Incident: api.Incident{ // Begins and ends in the future
				BeganAt: &timeLater1,
				EndedAt: &timeLater2,
			},
			ExpectedToBeInRange: false,
		},
		{
			Incident: api.Incident{ // Begins within, ends in the future
				BeganAt: &timeWithin2,
				EndedAt: &timeLater1,
			},
			ExpectedToBeInRange: true,
		},
		{
			Incident: api.Incident{ // Begins in the past, ends within
				BeganAt: &timePast1,
				EndedAt: &timeWithin1,
			},
			ExpectedToBeInRange: true,
		},
		{
			Incident: api.Incident{ // Begins in the past, ends in the past
				BeganAt: &timePast1,
				EndedAt: &timePast2,
			},
			ExpectedToBeInRange: false,
		},
		{
			Incident: api.Incident{ // Exactly within range
				BeganAt: &timeWithin1,
				EndedAt: &timeWithin2,
			},
			ExpectedToBeInRange: true,
		},
	}

	for i, testCase := range testCases {
		result := OverlapsTimeRange(&testCase.Incident, start, end)
		if result != testCase.ExpectedToBeInRange {
			t.Errorf(
				"incident %d: should be in range: %v, is in range: %v",
				i, testCase.ExpectedToBeInRange, result,
			)
		}
	}
}
