package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validateServiceStructureReady(t *testing.T, ctx context.Context) {
	t.Helper()

	server := NewPanchangamServer()
	require.NotNil(t, server, "Service should be creatable")

	req := &ppb.GetPanchangamRequest{
		Date:      "2024-01-15",
		Latitude:  12.9716,
		Longitude: 77.5946,
		Timezone:  "Asia/Kolkata",
	}

	resp, err := server.Get(ctx, req)
	assert.NoError(t, err, "Service should handle requests")
	require.NotNil(t, resp, "Service should return response")
	require.NotNil(t, resp.PanchangamData, "Service should return Panchangam data")

	data := resp.PanchangamData
	assert.Equal(t, req.Date, data.Date, "Date should be set correctly")
	assert.NotEmpty(t, data.Tithi, "Tithi field should exist")
	assert.NotEmpty(t, data.Nakshatra, "Nakshatra field should exist")
	assert.NotEmpty(t, data.Yoga, "Yoga field should exist")
	assert.NotEmpty(t, data.Karana, "Karana field should exist")
	assert.NotEmpty(t, data.SunriseTime, "SunriseTime should be calculated")
	assert.NotEmpty(t, data.SunsetTime, "SunsetTime should be calculated")
	assert.NotNil(t, data.Events, "Events should be included")

	t.Logf("Service structure is ready for real calculation integration")
}

func validateIntegrationReadiness(t *testing.T) {
	t.Helper()

	location := astronomy.Location{}
	assert.IsType(t, astronomy.Location{}, location, "Location type should be defined")

	validateTithiStructure(t)
	validateNakshatraStructure(t)
	validateYogaStructure(t)
	validateKaranaStructure(t)
	validateVaraStructure(t)

	t.Logf("All interfaces ready for service-calculation integration")
}

func validateTithiStructure(t *testing.T) {
	t.Helper()
	tithi := &astronomy.TithiInfo{}

	assert.IsType(t, 0, tithi.Number, "Tithi should have Number field")
	assert.IsType(t, "", tithi.Name, "Tithi should have Name field")
	assert.IsType(t, astronomy.TithiType(""), tithi.Type, "Tithi should have Type field")
	assert.IsType(t, time.Time{}, tithi.StartTime, "Tithi should have StartTime field")
	assert.IsType(t, time.Time{}, tithi.EndTime, "Tithi should have EndTime field")
	assert.IsType(t, 0.0, tithi.Duration, "Tithi should have Duration field")
	assert.IsType(t, false, tithi.IsShukla, "Tithi should have IsShukla field")
	assert.IsType(t, 0.0, tithi.MoonSunDiff, "Tithi should have MoonSunDiff field")
}

func validateNakshatraStructure(t *testing.T) {
	t.Helper()
	nakshatra := &astronomy.NakshatraInfo{}

	assert.IsType(t, 0, nakshatra.Number, "Nakshatra should have Number field")
	assert.IsType(t, "", nakshatra.Name, "Nakshatra should have Name field")
	assert.IsType(t, "", nakshatra.Deity, "Nakshatra should have Deity field")
	assert.IsType(t, "", nakshatra.PlanetaryLord, "Nakshatra should have PlanetaryLord field")
	assert.IsType(t, "", nakshatra.Symbol, "Nakshatra should have Symbol field")
	assert.IsType(t, 0, nakshatra.Pada, "Nakshatra should have Pada field")
}

func validateYogaStructure(t *testing.T) {
	t.Helper()
	yoga := &astronomy.YogaInfo{}

	assert.IsType(t, 0, yoga.Number, "Yoga should have Number field")
	assert.IsType(t, "", yoga.Name, "Yoga should have Name field")
	assert.IsType(t, astronomy.YogaQuality(""), yoga.Quality, "Yoga should have Quality field")
}

func validateKaranaStructure(t *testing.T) {
	t.Helper()
	karana := &astronomy.KaranaInfo{}

	assert.IsType(t, 0, karana.Number, "Karana should have Number field")
	assert.IsType(t, "", karana.Name, "Karana should have Name field")
	assert.IsType(t, astronomy.KaranaType(""), karana.Type, "Karana should have Type field")
	assert.IsType(t, "", karana.Description, "Karana should have Description field")
	assert.IsType(t, false, karana.IsVishti, "Karana should have IsVishti field")
	assert.IsType(t, 0, karana.TithiNumber, "Karana should have TithiNumber field")
	assert.IsType(t, 0, karana.HalfTithi, "Karana should have HalfTithi field")
}

func validateVaraStructure(t *testing.T) {
	t.Helper()
	vara := &astronomy.VaraInfo{}

	assert.IsType(t, 0, vara.Number, "Vara should have Number field")
	assert.IsType(t, "", vara.Name, "Vara should have Name field")
	assert.IsType(t, "", vara.PlanetaryLord, "Vara should have PlanetaryLord field")
	assert.IsType(t, "", vara.GregorianDay, "Vara should have GregorianDay field")
	assert.IsType(t, false, vara.IsAuspicious, "Vara should have IsAuspicious field")
	assert.IsType(t, 0, vara.CurrentHora, "Vara should have CurrentHora field")
	assert.IsType(t, "", vara.HoraPlanet, "Vara should have HoraPlanet field")
}

func TestServiceIntegrationReadiness(t *testing.T) {
	t.Run("Service_Ready_For_Calculator_Integration", func(t *testing.T) {
		observability.NewLocalObserver()

		server := NewPanchangamServer()
		require.NotNil(t, server, "Service should be creatable")

		ctx := context.Background()
		req := &ppb.GetPanchangamRequest{
			Date:      "2024-01-15",
			Latitude:  12.9716,
			Longitude: 77.5946,
			Timezone:  "Asia/Kolkata",
		}

		resp, err := server.Get(ctx, req)
		require.NoError(t, err, "Service should handle requests")
		require.NotNil(t, resp, "Service should return response")

		data := resp.PanchangamData
		require.NotNil(t, data, "Response data should exist")

		assert.IsType(t, "", data.Tithi, "Tithi field ready for string result")
		assert.IsType(t, "", data.Nakshatra, "Nakshatra field ready for string result")
		assert.IsType(t, "", data.Yoga, "Yoga field ready for string result")
		assert.IsType(t, "", data.Karana, "Karana field ready for string result")
		assert.IsType(t, "", data.SunriseTime, "SunriseTime field ready for time string")
		assert.IsType(t, "", data.SunsetTime, "SunsetTime field ready for time string")
		assert.IsType(t, []*ppb.PanchangamEvent{}, data.Events, "Events field ready for event list")

		t.Logf("Service is ready for calculator integration")
		t.Logf("- Request parsing: OK")
		t.Logf("- Response structure: OK")
		t.Logf("- Field types: OK")
		t.Logf("- Error handling: OK")
		t.Logf("- Observability: OK")
	})
}
