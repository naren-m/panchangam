package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateTithiFromLongitudes(t *testing.T) {
	calculator, _, _ := createTestTithiCalculator()
	ctx := context.Background()
	referenceDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		sunLong        float64
		moonLong       float64
		expectedTithi  int
		expectedShukla bool
		expectedType   TithiType
	}{
		{
			name:           "New Moon (Amavasya)",
			sunLong:        100.0,
			moonLong:       100.0,
			expectedTithi:  1,
			expectedShukla: true,
			expectedType:   TithiTypeNanda,
		},
		{
			name:           "First Quarter",
			sunLong:        100.0,
			moonLong:       190.0,
			expectedTithi:  8,
			expectedShukla: true,
			expectedType:   TithiTypeJaya,
		},
		{
			name:           "Full Moon (Purnima)",
			sunLong:        100.0,
			moonLong:       268.0,
			expectedTithi:  15,
			expectedShukla: true,
			expectedType:   TithiTypePurna,
		},
		{
			name:           "Third Quarter",
			sunLong:        100.0,
			moonLong:       10.0,
			expectedTithi:  23,
			expectedShukla: false,
			expectedType:   TithiTypeJaya,
		},
		{
			name:           "Cross Zero Longitude",
			sunLong:        350.0,
			moonLong:       10.0,
			expectedTithi:  2,
			expectedShukla: true,
			expectedType:   TithiTypeBhadra,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tithi, err := calculator.calculateTithiFromLongitudes(ctx, test.sunLong, test.moonLong, referenceDate, "traditional")

			require.NoError(t, err)
			require.NotNil(t, tithi)

			assert.Equal(t, test.expectedTithi, tithi.Number)
			assert.Equal(t, test.expectedShukla, tithi.IsShukla)
			assert.Equal(t, test.expectedType, tithi.Type)
			assert.Equal(t, TithiNames[test.expectedTithi], tithi.Name)

			err = ValidateTithiCalculation(tithi)
			assert.NoError(t, err)

			assert.True(t, tithi.EndTime.After(tithi.StartTime))
			assert.True(t, tithi.Duration > 0 && tithi.Duration < 48)
		})
	}
}

func TestGetTithiFromLongitudes(t *testing.T) {
	calculator, _, _ := createTestTithiCalculator()
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tithi, err := calculator.GetTithiFromLongitudes(ctx, 100.0, 190.0, date)

	require.NoError(t, err)
	require.NotNil(t, tithi)

	assert.Equal(t, 8, tithi.Number)
	assert.Equal(t, TithiNames[8], tithi.Name)
	assert.True(t, tithi.IsShukla)
	assert.Equal(t, TithiTypeJaya, tithi.Type)
}

func TestTithiNumberFromLongitudes(t *testing.T) {
	assert.Equal(t, 1, tithiNumberFromLongitudes(100, 100))
	assert.Equal(t, 8, tithiNumberFromLongitudes(100, 190))
	assert.Equal(t, 30, tithiNumberFromLongitudes(355, 354))
	assert.Equal(t, 1, tithiNumberFromLongitudes(354, 355))
}

func TestTithiCalculation_EdgeCases(t *testing.T) {
	calculator, _, _ := createTestTithiCalculator()
	ctx := context.Background()
	referenceDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		sunLong  float64
		moonLong float64
	}{
		{
			name:     "Exact boundary - 360 degrees",
			sunLong:  0.0,
			moonLong: 360.0,
		},
		{
			name:     "Large longitude values",
			sunLong:  720.0,
			moonLong: 800.0,
		},
		{
			name:     "Negative longitude (should be normalized)",
			sunLong:  350.0,
			moonLong: -10.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tithi, err := calculator.calculateTithiFromLongitudes(ctx, test.sunLong, test.moonLong, referenceDate, "traditional")

			require.NoError(t, err)
			require.NotNil(t, tithi)

			err = ValidateTithiCalculation(tithi)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkCalculateTithiFromLongitudes(b *testing.B) {
	calculator, _, _ := createTestTithiCalculator()
	ctx := context.Background()
	referenceDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculator.calculateTithiFromLongitudes(ctx, 100.0, 190.0, referenceDate, "traditional")
		if err != nil {
			b.Fatal(err)
		}
	}
}
