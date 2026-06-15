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

func TestTraceAttribute(t *testing.T) {
	attr := traceAttribute("test_key", "test_value")
	assert.Equal(t, "test_key", string(attr.Key))
	assert.Equal(t, "test_value", attr.Value.AsString())
}

func TestTraceAttributes(t *testing.T) {
	attrs := traceAttributes("key1", "value1", "key2", "value2")
	assert.NotNil(t, attrs)
	assert.Len(t, attrs, 1)

	attrs = traceAttributes("key1", "value1", "key2")
	assert.Nil(t, attrs)

	attrs = traceAttributes()
	assert.NotNil(t, attrs)
	assert.Len(t, attrs, 1)
}

func TestPanchangamServer_LoggingPaths(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	req := &ppb.GetPanchangamRequest{
		Date:              "2024-06-21",
		Latitude:          40.7128,
		Longitude:         -74.0060,
		Timezone:          "America/New_York",
		Region:            "North America",
		CalculationMethod: "Drik",
		Locale:            "en-US",
	}

	resp, err := server.Get(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.PanchangamData)
}

func TestPanchangamServer_LoggingEmptyFields(t *testing.T) {
	observability.NewLocalObserver()
	server := NewPanchangamServer()

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-06-21",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	resp, err := server.Get(context.Background(), req)
	if err != nil && status.Code(err) == codes.Internal {
		assert.Contains(t, err.Error(), "failed to fetch panchangam data")
	} else {
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}
}
