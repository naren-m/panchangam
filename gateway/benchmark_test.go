package gateway

import (
	"net/http/httptest"
	"testing"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/mock"
)

func BenchmarkHandlePanchangam(b *testing.B) {
	mockClient := new(MockPanchangamClient)
	response := &ppb.GetPanchangamResponse{
		PanchangamData: &ppb.PanchangamData{
			Date:        "2024-01-15",
			Tithi:       "Test Tithi",
			Nakshatra:   "Test Nakshatra",
			Yoga:        "Test Yoga",
			Karana:      "Test Karana",
			SunriseTime: "06:45:32",
			SunsetTime:  "18:21:47",
		},
	}
	mockClient.On("Get", mock.Anything, mock.Anything).Return(response, nil)

	server := &GatewayServer{}
	handler := server.handlePanchangam(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler(w, req)
	}
}
