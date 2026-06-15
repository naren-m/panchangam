package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleCurrentTithi_Success(t *testing.T) {
	mockClient := new(MockPanchangamClient)

	expectedResponse := &ppb.GetPanchangamResponse{
		PanchangamData: &ppb.PanchangamData{
			Date:        "2024-01-15",
			Tithi:       "Thuthiya - Krishna Paksha Day 3 (Purnimanta)",
			Nakshatra:   "Shravana (22)",
			Yoga:        "Vaidhriti (27)",
			Karana:      "Bava (2)",
			SunriseTime: "06:45:32",
			SunsetTime:  "18:21:47",
			Events: []*ppb.PanchangamEvent{
				{
					Name:      "Vara: Somavara",
					Time:      "06:45:32",
					EventType: "VARA",
				},
				{
					Name:      "Abhijit Muhurta",
					Time:      "12:03:00",
					EventType: "ABHIJIT_MUHURTA",
				},
			},
		},
	}

	mockClient.On("Get", mock.Anything, mock.MatchedBy(func(req *ppb.GetPanchangamRequest) bool {
		return req.Date == "2024-01-15T18:00:00+05:30" &&
			req.Latitude == 12.9716 &&
			req.Longitude == 77.5946 &&
			req.Timezone == "Asia/Kolkata" &&
			req.Region == "global" &&
			req.CalculationMethod == "traditional" &&
			req.Locale == "en"
	})).Return(expectedResponse, nil)

	server := &GatewayServer{}
	handler := server.handleCurrentTithi(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/tithi/current?date=2024-01-15T12:30:00Z&latitude=12.9716&longitude=77.5946&timezone=Asia/Kolkata&calendar_system=Purnimanta", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "2024-01-15", result["date"])

	tithi := result["tithi"].(map[string]interface{})
	assert.Equal(t, "Thuthiya", tithi["name"])
	assert.Equal(t, "Thuthiya", tithi["traditional_name"])
	assert.Equal(t, float64(18), tithi["number"])
	assert.Equal(t, "Krishna", tithi["paksha"])
	assert.Equal(t, float64(3), tithi["paksha_day"])
	assert.Equal(t, "Jaya", tithi["type"])

	panchaAnga := result["pancha_anga"].(map[string]interface{})
	assert.Equal(t, "Shravana (22)", panchaAnga["nakshatra"])
	assert.Equal(t, "Vaidhriti (27)", panchaAnga["yoga"])
	assert.Equal(t, "Bava (2)", panchaAnga["karana"])
	assert.Equal(t, "Somavara", panchaAnga["vara"])

	day := result["day"].(map[string]interface{})
	assert.Equal(t, "06:45:32", day["sunrise_time"])
	assert.Equal(t, "18:21:47", day["sunset_time"])

	calculation := result["calculation"].(map[string]interface{})
	assert.Equal(t, "Asia/Kolkata", calculation["timezone"])
	assert.Equal(t, "global", calculation["region"])
	assert.Equal(t, "Purnimanta", calculation["calendar_system"])

	_, err = time.Parse(time.RFC3339, result["generated_at"].(string))
	assert.NoError(t, err)
	_, err = time.Parse(time.RFC3339, result["next_refresh_at"].(string))
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestMakeCurrentTithiResponseUsesTithiWindowEvents(t *testing.T) {
	generatedAt := time.Date(2026, 6, 4, 16, 30, 0, 0, time.UTC)
	referenceAt := generatedAt
	data := &ppb.PanchangamData{
		Date:        "2026-06-04",
		Tithi:       "Chathurthi - Krishna Paksha Day 4 (Purnimanta)",
		Nakshatra:   "Uttara Ashadha (21)",
		Yoga:        "Brahma (25)",
		Karana:      "Kaulava (4)",
		SunriseTime: "05:47:41",
		SunsetTime:  "20:23:56",
		Events: []*ppb.PanchangamEvent{
			{
				Name:      "Tithi starts",
				Time:      "09:05:00",
				EventType: "TITHI_START",
			},
			{
				Name:      "Tithi ends",
				Time:      "23:00:00",
				EventType: "TITHI_END",
			},
			{
				Name:      "Raasi: Makara",
				Time:      "11:00:00",
				EventType: "RAASI",
			},
		},
	}

	result := makeCurrentTithiResponse(data, "America/Los_Angeles", "global", "Purnimanta", "Drik", "en", generatedAt, referenceAt)
	encoded, err := json.Marshal(result)
	require.NoError(t, err)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &payload))

	assert.Equal(t, "2026-06-04T16:05:00Z", result.Tithi.StartTime)
	assert.Equal(t, "2026-06-05T06:00:00Z", result.Tithi.EndTime)
	assert.Equal(t, "2026-06-05T06:00:00Z", result.NextRefreshAt)
	assert.Equal(t, "Makara", payload["raasi"])
}

func TestMakeCurrentTithiResponseRollsShortClockWindowToNextDay(t *testing.T) {
	generatedAt := time.Date(2026, 6, 4, 22, 30, 0, 0, time.UTC)
	referenceAt := time.Date(2026, 6, 4, 22, 28, 4, 0, time.UTC)
	data := &ppb.PanchangamData{
		Date:        "2026-06-04",
		Tithi:       "Panchami - Krishna Paksha Day 5 (Purnimanta)",
		Nakshatra:   "Shravana (22)",
		Yoga:        "Brahma (25)",
		Karana:      "Kaulava (4)",
		SunriseTime: "05:47:41",
		SunsetTime:  "20:23:56",
		Events: []*ppb.PanchangamEvent{
			{
				Name:      "Tithi starts",
				Time:      "11:15:37",
				EventType: "TITHI_START",
			},
			{
				Name:      "Tithi ends",
				Time:      "12:03:01",
				EventType: "TITHI_END",
			},
		},
	}

	result := makeCurrentTithiResponse(data, "America/Los_Angeles", "global", "Purnimanta", "Drik", "en", generatedAt, referenceAt)

	assert.Equal(t, "2026-06-04T18:15:37Z", result.Tithi.StartTime)
	assert.Equal(t, "2026-06-05T19:03:01Z", result.Tithi.EndTime)
	assert.Equal(t, "2026-06-05T19:03:01Z", result.NextRefreshAt)
}

func TestParseQueryFloat(t *testing.T) {
	tests := []struct {
		name      string
		query     map[string][]string
		names     []string
		want      float64
		wantValue string
		wantOK    bool
		wantErr   bool
	}{
		{
			name:      "trims value",
			query:     map[string][]string{"lat": {" 12.5 "}},
			names:     []string{"lat"},
			want:      12.5,
			wantValue: "12.5",
			wantOK:    true,
		},
		{
			name:      "uses fallback alias",
			query:     map[string][]string{"latitude": {""}, "lat": {"12.5"}},
			names:     []string{"latitude", "lat"},
			want:      12.5,
			wantValue: "12.5",
			wantOK:    true,
		},
		{
			name:    "missing value",
			query:   map[string][]string{"lat": {""}},
			names:   []string{"lat"},
			wantOK:  false,
			wantErr: false,
		},
		{
			name:      "invalid value keeps original input",
			query:     map[string][]string{"lng": {"bad"}},
			names:     []string{"lng"},
			wantValue: "bad",
			wantOK:    true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotValue, gotOK, err := parseQueryFloat(tt.query, tt.names...)

			assert.Equal(t, tt.wantOK, gotOK)
			assert.Equal(t, tt.wantValue, gotValue)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.wantOK {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
