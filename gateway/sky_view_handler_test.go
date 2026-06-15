package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSkyView_InvalidParameters(t *testing.T) {
	tests := []struct {
		name           string
		target         string
		expectedMsg    string
		expectedParam  string
		expectedValue  string
		expectedStatus int
	}{
		{
			name:           "invalid latitude",
			target:         "/api/v1/sky-view?lat=bad&lng=77.5946",
			expectedMsg:    "Invalid latitude. Must be between -90 and 90",
			expectedParam:  "lat",
			expectedValue:  "bad",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid longitude",
			target:         "/api/v1/sky-view?lat=12.9716&lng=bad",
			expectedMsg:    "Invalid longitude. Must be between -180 and 180",
			expectedParam:  "lng",
			expectedValue:  "bad",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid altitude",
			target:         "/api/v1/sky-view?lat=12.9716&lng=77.5946&alt=bad",
			expectedMsg:    "Invalid altitude format",
			expectedParam:  "alt",
			expectedValue:  "bad",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &GatewayServer{ephemerisProvider: skyViewValidationProvider{}}
			handler := server.handleSkyView()

			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var errResp APIError
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &errResp))
			assert.Equal(t, "INVALID_PARAMETER", errResp.Error.Code)
			assert.Equal(t, tt.expectedMsg, errResp.Error.Message)
			assert.Equal(t, tt.expectedParam, errResp.Error.Details["parameter"])
			assert.Equal(t, tt.expectedValue, errResp.Error.Details["value"])
		})
	}
}

func TestParseSkyViewObservationTime(t *testing.T) {
	fallback := time.Date(2026, time.June, 5, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		date      string
		timeValue string
		want      time.Time
	}{
		{
			name: "uses fallback when date is empty",
			want: fallback,
		},
		{
			name: "parses date and trims whitespace",
			date: " 2024-06-21 ",
			want: time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "parses date time and trims whitespace",
			date:      " 2024-06-21 ",
			timeValue: " 18:30:45 ",
			want:      time.Date(2024, time.June, 21, 18, 30, 45, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSkyViewObservationTime(tt.date, tt.timeValue, fallback)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseSkyViewObservationTimeRejectsInvalidInput(t *testing.T) {
	fallback := time.Date(2026, time.June, 5, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		date      string
		timeValue string
	}{
		{
			name: "invalid date",
			date: "2024-99-99",
		},
		{
			name:      "invalid date time",
			date:      "2024-06-21",
			timeValue: "bad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSkyViewObservationTime(tt.date, tt.timeValue, fallback)
			require.Error(t, err)
		})
	}
}

type skyViewValidationProvider struct{}

func (skyViewValidationProvider) GetPlanetaryPositions(context.Context, ephemeris.JulianDay) (*ephemeris.PlanetaryPositions, error) {
	return nil, nil
}

func (skyViewValidationProvider) GetSunPosition(context.Context, ephemeris.JulianDay) (*ephemeris.SolarPosition, error) {
	return nil, nil
}

func (skyViewValidationProvider) GetMoonPosition(context.Context, ephemeris.JulianDay) (*ephemeris.LunarPosition, error) {
	return nil, nil
}

func (skyViewValidationProvider) IsAvailable(context.Context) bool {
	return true
}

func (skyViewValidationProvider) GetDataRange() (ephemeris.JulianDay, ephemeris.JulianDay) {
	return 0, 0
}

func (skyViewValidationProvider) GetHealthStatus(context.Context) (*ephemeris.HealthStatus, error) {
	return &ephemeris.HealthStatus{}, nil
}

func (skyViewValidationProvider) GetProviderName() string {
	return "sky-view-validation"
}

func (skyViewValidationProvider) GetVersion() string {
	return "test"
}

func (skyViewValidationProvider) Close() error {
	return nil
}
