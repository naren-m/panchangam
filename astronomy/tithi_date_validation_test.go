package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTithiForDate(t *testing.T) {
	observability.NewLocalObserver()
	manager := ephemeris.NewManager(ephemeris.NewSwissProvider(), nil, ephemeris.NewMemoryCache(100, time.Hour))
	calculator := NewTithiCalculator(manager)
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tithi, err := calculator.GetTithiForDate(ctx, date)

	require.NoError(t, err)
	require.NotNil(t, tithi)

	assert.True(t, tithi.Number >= 1 && tithi.Number <= 30)
	assert.NotEmpty(t, tithi.Name)
	assert.True(t, tithi.Duration > 0)
	assert.True(t, tithi.EndTime.After(tithi.StartTime))
}

func TestGetTithiForDate_EphemerisError(t *testing.T) {
	calculator, mockProvider, mockCache := createTestTithiCalculator()
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, false)
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.AnythingOfType("ephemeris.JulianDay")).
		Return((*ephemeris.PlanetaryPositions)(nil), assert.AnError)

	tithi, err := calculator.GetTithiForDate(ctx, date)

	assert.Error(t, err)
	assert.Nil(t, tithi)
	assert.Contains(t, err.Error(), "failed to get planetary positions")

	mockProvider.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestValidateTithiCalculation(t *testing.T) {
	validTithi := &TithiInfo{
		Number:      8,
		Name:        "Ashtami",
		Type:        TithiTypeJaya,
		StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2024, 1, 15, 30, 0, 0, 0, time.UTC),
		Duration:    24.0,
		IsShukla:    true,
		MoonSunDiff: 90.0,
	}

	tests := []struct {
		name          string
		tithi         *TithiInfo
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid Tithi",
			tithi:       validTithi,
			expectError: false,
		},
		{
			name:          "Nil Tithi",
			tithi:         nil,
			expectError:   true,
			errorContains: "tithi cannot be nil",
		},
		{
			name: "Invalid Tithi Number - Too Low",
			tithi: &TithiInfo{
				Number:      0,
				Duration:    24.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 30, 0, 0, 0, time.UTC),
				MoonSunDiff: 90.0,
			},
			expectError:   true,
			errorContains: "invalid tithi number",
		},
		{
			name: "Invalid Tithi Number - Too High",
			tithi: &TithiInfo{
				Number:      31,
				Duration:    24.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 30, 0, 0, 0, time.UTC),
				MoonSunDiff: 90.0,
			},
			expectError:   true,
			errorContains: "invalid tithi number",
		},
		{
			name: "Invalid Moon-Sun Difference - Negative",
			tithi: &TithiInfo{
				Number:      8,
				Duration:    24.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 30, 0, 0, 0, time.UTC),
				MoonSunDiff: -10.0,
			},
			expectError:   true,
			errorContains: "invalid moon-sun difference",
		},
		{
			name: "Invalid Moon-Sun Difference - Too High",
			tithi: &TithiInfo{
				Number:      8,
				Duration:    24.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 30, 0, 0, 0, time.UTC),
				MoonSunDiff: 370.0,
			},
			expectError:   true,
			errorContains: "invalid moon-sun difference",
		},
		{
			name: "Invalid Duration - Zero",
			tithi: &TithiInfo{
				Number:      8,
				Duration:    0.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				MoonSunDiff: 90.0,
			},
			expectError:   true,
			errorContains: "invalid tithi duration",
		},
		{
			name: "Invalid Duration - Too Long",
			tithi: &TithiInfo{
				Number:      8,
				Duration:    50.0,
				StartTime:   time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 17, 8, 0, 0, 0, time.UTC),
				MoonSunDiff: 90.0,
			},
			expectError:   true,
			errorContains: "invalid tithi duration",
		},
		{
			name: "End Time Before Start Time",
			tithi: &TithiInfo{
				Number:      8,
				Duration:    24.0,
				StartTime:   time.Date(2024, 1, 15, 18, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
				MoonSunDiff: 90.0,
			},
			expectError:   true,
			errorContains: "tithi end time cannot be before start time",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateTithiCalculation(test.tithi)

			if test.expectError {
				assert.Error(t, err)
				if test.errorContains != "" {
					assert.Contains(t, err.Error(), test.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
