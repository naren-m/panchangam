package implementations

import (
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

func TestCalculationMethodPluginTithiHelpers(t *testing.T) {
	plugin := NewCalculationMethodPlugin(&ephemeris.Manager{})
	start := time.Date(2024, time.January, 15, 6, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	tithi := &astronomy.TithiInfo{
		Number:    16,
		StartTime: start,
		EndTime:   end,
		Type:      astronomy.TithiTypeJaya,
	}

	if got := plugin.calculateTithiPercentage(tithi, start.Add(6*time.Hour)); got != 25 {
		t.Fatalf("expected 25 percent, got %f", got)
	}
	if !plugin.isTithiRunning(tithi, start.Add(time.Hour)) {
		t.Fatal("expected tithi to be running")
	}
	if plugin.isTithiRunning(tithi, start.Add(-time.Hour)) {
		t.Fatal("expected tithi before start to be inactive")
	}
	if got := plugin.getTithiLord(tithi.Number); got != "Agni" {
		t.Fatalf("expected repeated paksha lord Agni, got %s", got)
	}
	if got := plugin.getTithiQuality(tithi.Type); got != "victorious" {
		t.Fatalf("expected victorious quality, got %s", got)
	}
}

func TestCalculationMethodPluginNakshatraHelpers(t *testing.T) {
	plugin := NewCalculationMethodPlugin(&ephemeris.Manager{})
	date := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)

	nakshatra := plugin.calculateNakshatraFromLongitude(0, date)
	if nakshatra.Number != 1 || nakshatra.Name != "Ashwini" || nakshatra.Pada != 1 {
		t.Fatalf("unexpected first nakshatra: %+v", nakshatra)
	}
	if nakshatra.Lord != "Ketu" || nakshatra.Deity != "Ashwini Kumaras" || nakshatra.Symbol != "Horse's head" {
		t.Fatalf("unexpected Ashwini metadata: %+v", nakshatra)
	}

	nakshatra = plugin.calculateNakshatraFromLongitude(359, date)
	if nakshatra.Number != 27 || nakshatra.Name != "Revati" {
		t.Fatalf("unexpected final nakshatra: %+v", nakshatra)
	}
}
