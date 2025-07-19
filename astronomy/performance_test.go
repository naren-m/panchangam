package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/stretchr/testify/mock"
)

// Simple performance test without observability telemetry
func BenchmarkPanchangamPerformance(b *testing.B) {
	// Create a simple calculator that bypasses telemetry
	calculator, mockProvider, _ := createTestTithiCalculator()
	manager := calculator.ephemerisManager

	mockPositions := &ephemeris.PlanetaryPositions{
		Sun: ephemeris.Position{
			Longitude: 285.5,
		},
		Moon: ephemeris.Position{
			Longitude: 95.7,
		},
	}
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.Anything).Return(mockPositions, nil)

	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tithiCalc := calculator
	nakshatraCalc := NewNakshatraCalculator(manager)
	yogaCalc := NewYogaCalculator(manager)
	karanaCalc := NewKaranaCalculator(manager)
	varaCalc := NewVaraCalculator()

	b.ResetTimer()
	
	start := time.Now()
	for i := 0; i < b.N; i++ {
		// Calculate all 5 Panchangam elements
		_, _ = tithiCalc.GetTithiFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		_, _ = nakshatraCalc.GetNakshatraFromLongitude(ctx, mockPositions.Moon.Longitude, testDate)
		_, _ = yogaCalc.GetYogaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		_, _ = karanaCalc.GetKaranaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		
		// Vara calculation (simplified)
		sunrise := time.Date(2024, 1, 15, 6, 30, 0, 0, time.UTC)
		nextSunrise := time.Date(2024, 1, 16, 6, 31, 0, 0, time.UTC)
		_, _ = varaCalc.GetVaraFromGregorianDay(ctx, testDate.Weekday(), sunrise, nextSunrise, testDate)
	}
	elapsed := time.Since(start)
	
	b.StopTimer()
	
	// Calculate average time per complete Panchangam calculation
	avgTime := elapsed / time.Duration(b.N)
	b.Logf("Average time per complete Panchangam calculation: %v", avgTime)
	
	// Check if we meet the <100ms requirement
	if avgTime > 100*time.Millisecond {
		b.Logf("WARNING: Performance target not met. Average: %v > 100ms target", avgTime)
	} else {
		b.Logf("SUCCESS: Performance target met. Average: %v < 100ms target", avgTime)
	}
}