package panchangam

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPanchangamServer_EmptyDate(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	for _, date := range []string{"", "   "} {
		t.Run(fmt.Sprintf("date=%q", date), func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      date,
				Latitude:  40.7128,
				Longitude: -74.0060,
			}

			resp, err := server.Get(context.Background(), req)
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			assert.Contains(t, err.Error(), "date parameter is required")
		})
	}
}

func TestPanchangamServer_LongTimezone(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timezone:  "This/Is/A/Very/Long/Invalid/Timezone/String/That/Should/Cause/An/Error",
	}

	resp, err := server.Get(context.Background(), req)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "invalid timezone")
	assert.Nil(t, resp)
}

func TestPanchangamServer_BoundaryCoordinates(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		shouldErr bool
	}{
		{"Valid boundary - North Pole", 90.0, 0.0, false},
		{"Valid boundary - South Pole", -90.0, 0.0, false},
		{"Valid boundary - East boundary", 0.0, 180.0, false},
		{"Valid boundary - West boundary", 0.0, -180.0, false},
		{"Just over north boundary", 90.1, 0.0, true},
		{"Just over south boundary", -90.1, 0.0, true},
		{"Just over east boundary", 0.0, 180.1, true},
		{"Just over west boundary", 0.0, -180.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      "2024-06-21",
				Latitude:  tt.latitude,
				Longitude: tt.longitude,
			}

			resp, err := server.Get(context.Background(), req)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestPanchangamServer_Initialization(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	assert.NotNil(t, server)
	assert.NotNil(t, server.observer)
	assert.IsType(t, &PanchangamServer{}, server)
}

func TestPanchangamServer_DateFormats(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	invalidDates := []string{
		"2024-13-01",
		"2024-02-30",
		"2024-06-32",
		"24-06-21",
		"2024/06/21",
		"2024-6-21",
		"2024-06-1",
		"invalid-date",
	}

	for _, date := range invalidDates {
		t.Run(fmt.Sprintf("Invalid date: %s", date), func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      date,
				Latitude:  40.7128,
				Longitude: -74.0060,
			}

			resp, err := server.Get(context.Background(), req)
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestParsePanchangamDateInputTrimsWhitespace(t *testing.T) {
	parsedDate, parsedDateTime, hasClockTime, err := parsePanchangamDateInput(" 2024-06-21 ")
	require.NoError(t, err)
	assert.False(t, hasClockTime)
	assert.True(t, parsedDateTime.IsZero())
	assert.Equal(t, time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC), parsedDate)

	parsedDate, parsedDateTime, hasClockTime, err = parsePanchangamDateInput(" 2024-06-21T17:45:00Z ")
	require.NoError(t, err)
	assert.True(t, hasClockTime)
	assert.True(t, parsedDate.IsZero())
	assert.Equal(t, time.Date(2024, time.June, 21, 17, 45, 0, 0, time.UTC), parsedDateTime)
}
