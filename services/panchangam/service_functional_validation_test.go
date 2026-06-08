package panchangam

import (
	"context"
	"testing"

	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceFunctionalInputValidation(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	testCases := []struct {
		name          string
		request       *ppb.GetPanchangamRequest
		expectedError codes.Code
		description   string
	}{
		{
			name: "Invalid_Latitude_High",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  91.0,
				Longitude: 77.5946,
			},
			expectedError: codes.InvalidArgument,
			description:   "Latitude above 90 should be rejected",
		},
		{
			name: "Invalid_Latitude_Low",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  -91.0,
				Longitude: 77.5946,
			},
			expectedError: codes.InvalidArgument,
			description:   "Latitude below -90 should be rejected",
		},
		{
			name: "Invalid_Longitude_High",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  12.9716,
				Longitude: 181.0,
			},
			expectedError: codes.InvalidArgument,
			description:   "Longitude above 180 should be rejected",
		},
		{
			name: "Invalid_Longitude_Low",
			request: &ppb.GetPanchangamRequest{
				Date:      "2024-01-15",
				Latitude:  12.9716,
				Longitude: -181.0,
			},
			expectedError: codes.InvalidArgument,
			description:   "Longitude below -180 should be rejected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.Get(ctx, tc.request)

			assert.Error(t, err, tc.description)
			assert.Nil(t, resp, "Response should be nil for invalid request")

			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be a gRPC status error")
			assert.Equal(t, tc.expectedError, st.Code(), "Error code should match expected")
		})
	}
}

func TestServiceFunctionalDateValidation(t *testing.T) {
	observability.NewLocalObserver()

	server := NewPanchangamServer()
	ctx := context.Background()

	invalidDateCases := []string{
		"invalid-date",
		"2024-13-01",
		"2024-01-32",
		"24-01-01",
		"2024/01/01",
	}

	for _, invalidDate := range invalidDateCases {
		t.Run("Invalid_Date_"+invalidDate, func(t *testing.T) {
			req := &ppb.GetPanchangamRequest{
				Date:      invalidDate,
				Latitude:  12.9716,
				Longitude: 77.5946,
			}

			resp, err := server.Get(ctx, req)

			assert.Error(t, err, "Invalid date should cause error")
			assert.Nil(t, resp, "Response should be nil for invalid date")

			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be a gRPC status error")
			assert.Equal(t, codes.InvalidArgument, st.Code(), "Should return InvalidArgument for bad date")
		})
	}
}
