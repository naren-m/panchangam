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

// MockEphemerisProvider is a mock implementation of EphemerisProvider for testing
type MockEphemerisProvider struct {
	mock.Mock
}

func (m *MockEphemerisProvider) GetPlanetaryPositions(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.PlanetaryPositions, error) {
	args := m.Called(ctx, jd)
	return args.Get(0).(*ephemeris.PlanetaryPositions), args.Error(1)
}

func (m *MockEphemerisProvider) GetSunPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.SolarPosition, error) {
	args := m.Called(ctx, jd)
	return args.Get(0).(*ephemeris.SolarPosition), args.Error(1)
}

func (m *MockEphemerisProvider) GetMoonPosition(ctx context.Context, jd ephemeris.JulianDay) (*ephemeris.LunarPosition, error) {
	args := m.Called(ctx, jd)
	return args.Get(0).(*ephemeris.LunarPosition), args.Error(1)
}

func (m *MockEphemerisProvider) IsAvailable(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockEphemerisProvider) GetDataRange() (startJD, endJD ephemeris.JulianDay) {
	args := m.Called()
	return args.Get(0).(ephemeris.JulianDay), args.Get(1).(ephemeris.JulianDay)
}

func (m *MockEphemerisProvider) GetHealthStatus(ctx context.Context) (*ephemeris.HealthStatus, error) {
	args := m.Called(ctx)
	return args.Get(0).(*ephemeris.HealthStatus), args.Error(1)
}

func (m *MockEphemerisProvider) GetProviderName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEphemerisProvider) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEphemerisProvider) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockCache is a mock implementation of Cache for testing
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (interface{}, bool) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, duration time.Duration) {
	m.Called(ctx, key, value, duration)
}

func (m *MockCache) Delete(ctx context.Context, key string) bool {
	args := m.Called(ctx, key)
	return args.Bool(0)
}

func (m *MockCache) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) GetStats(ctx context.Context) *ephemeris.CacheStats {
	args := m.Called(ctx)
	return args.Get(0).(*ephemeris.CacheStats)
}

func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Helper function to create a test TithiCalculator with mocked dependencies
func createTestTithiCalculator() (*TithiCalculator, *MockEphemerisProvider, *MockCache) {
	// Initialize observability for testing
	observability.NewLocalObserver()

	mockProvider := &MockEphemerisProvider{}
	mockCache := &MockCache{}

	// Set up basic expectations for provider metadata
	mockProvider.On("GetProviderName").Return("MockProvider")
	mockProvider.On("GetVersion").Return("1.0.0")

	manager := ephemeris.NewManager(mockProvider, nil, mockCache)
	calculator := NewTithiCalculator(manager)

	return calculator, mockProvider, mockCache
}

func TestNewTithiCalculator(t *testing.T) {
	calculator, _, _ := createTestTithiCalculator()

	assert.NotNil(t, calculator)
	assert.NotNil(t, calculator.ephemerisManager)
	assert.NotNil(t, calculator.observer)
}

func TestGetTithiTypeDescription(t *testing.T) {
	tests := []struct {
		tithiType    TithiType
		expectedDesc string
	}{
		{TithiTypeNanda, "Joyful, good for celebrations and new beginnings"},
		{TithiTypeBhadra, "Auspicious, good for all activities"},
		{TithiTypeJaya, "Victorious, good for achieving success"},
		{TithiTypeRikta, "Empty, avoid starting new ventures"},
		{TithiTypePurna, "Complete, excellent for completion of tasks"},
		{TithiType("Invalid"), "Unknown Tithi type"},
	}

	for _, test := range tests {
		t.Run(string(test.tithiType), func(t *testing.T) {
			desc := GetTithiTypeDescription(test.tithiType)
			assert.Equal(t, test.expectedDesc, desc)
		})
	}
}

func TestGetTithiType(t *testing.T) {
	tests := []struct {
		tithiNumber  int
		expectedType TithiType
	}{
		// Shukla Paksha
		{1, TithiTypeNanda},   // Pratipada
		{2, TithiTypeBhadra},  // Dwitiya
		{3, TithiTypeJaya},    // Tritiya
		{4, TithiTypeRikta},   // Chaturthi
		{5, TithiTypePurna},   // Panchami
		{6, TithiTypeNanda},   // Shashthi
		{7, TithiTypeBhadra},  // Saptami
		{8, TithiTypeJaya},    // Ashtami
		{9, TithiTypeRikta},   // Navami
		{10, TithiTypePurna},  // Dashami
		{11, TithiTypeNanda},  // Ekadashi
		{12, TithiTypeBhadra}, // Dwadashi
		{13, TithiTypeJaya},   // Trayodashi
		{14, TithiTypeRikta},  // Chaturdashi
		{15, TithiTypePurna},  // Purnima

		// Krishna Paksha (should follow same pattern)
		{16, TithiTypeNanda},  // Pratipada
		{17, TithiTypeBhadra}, // Dwitiya
		{18, TithiTypeJaya},   // Tritiya
		{19, TithiTypeRikta},  // Chaturthi
		{20, TithiTypePurna},  // Panchami
		{25, TithiTypePurna},  // Dashami
		{30, TithiTypePurna},  // Amavasya
	}

	for _, test := range tests {
		t.Run(TithiNames[test.tithiNumber], func(t *testing.T) {
			tithiType := getTithiType(test.tithiNumber)
			assert.Equal(t, test.expectedType, tithiType)
		})
	}
}

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
			moonLong:       100.0, // Same longitude = 0° difference
			expectedTithi:  1,     // First Tithi after new moon
			expectedShukla: true,
			expectedType:   TithiTypeNanda,
		},
		{
			name:           "First Quarter",
			sunLong:        100.0,
			moonLong:       190.0, // 90° difference
			expectedTithi:  8,     // 90° / 12° + 1 = 8th Tithi
			expectedShukla: true,
			expectedType:   TithiTypeJaya,
		},
		{
			name:           "Full Moon (Purnima)",
			sunLong:        100.0,
			moonLong:       268.0, // 168° difference (14th Tithi * 12° = 168°, so we're in 15th Tithi)
			expectedTithi:  15,    // 168° / 12° + 1 = 15th Tithi (Purnima)
			expectedShukla: true,
			expectedType:   TithiTypePurna,
		},
		{
			name:           "Third Quarter",
			sunLong:        100.0,
			moonLong:       10.0, // 270° difference (Moon ahead by 270°)
			expectedTithi:  23,   // 270° / 12° + 1 = 23rd Tithi (Krishna Paksha)
			expectedShukla: false,
			expectedType:   TithiTypeJaya,
		},
		{
			name:           "Cross Zero Longitude",
			sunLong:        350.0,
			moonLong:       10.0, // 20° difference when crossing 0°
			expectedTithi:  2,    // 20° / 12° + 1 = 2nd Tithi
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

			// Validate the calculation
			err = ValidateTithiCalculation(tithi)
			assert.NoError(t, err)

			// Check that times are reasonable
			assert.True(t, tithi.EndTime.After(tithi.StartTime))
			assert.True(t, tithi.Duration > 0 && tithi.Duration < 48) // Should be between 0 and 48 hours
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

func TestGetTithiForDate(t *testing.T) {
	calculator, mockProvider, mockCache := createTestTithiCalculator()
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Set up mock responses
	expectedPositions := &ephemeris.PlanetaryPositions{
		Sun: ephemeris.Position{
			Longitude: 295.5, // Sun in Capricorn
		},
		Moon: ephemeris.Position{
			Longitude: 385.5, // Moon ahead by 90° (but normalized to 25.5°)
		},
	}

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, false)
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration"))

	// Mock ephemeris provider
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.AnythingOfType("ephemeris.JulianDay")).
		Return(expectedPositions, nil)

	tithi, err := calculator.GetTithiForDate(ctx, date)

	require.NoError(t, err)
	require.NotNil(t, tithi)

	// Verify the calculation
	assert.True(t, tithi.Number >= 1 && tithi.Number <= 30)
	assert.NotEmpty(t, tithi.Name)
	assert.True(t, tithi.Duration > 0)
	assert.True(t, tithi.EndTime.After(tithi.StartTime))

	// Verify mock calls
	mockProvider.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetTithiForDate_EphemerisError(t *testing.T) {
	calculator, mockProvider, mockCache := createTestTithiCalculator()
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, false)

	// Mock ephemeris provider error
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.AnythingOfType("ephemeris.JulianDay")).
		Return((*ephemeris.PlanetaryPositions)(nil), assert.AnError)

	tithi, err := calculator.GetTithiForDate(ctx, date)

	assert.Error(t, err)
	assert.Nil(t, tithi)
	assert.Contains(t, err.Error(), "failed to get planetary positions")

	// Verify mock calls
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

func TestTithiNames(t *testing.T) {
	// Test that all Tithi numbers have names
	for i := 1; i <= 30; i++ {
		name, exists := TithiNames[i]
		assert.True(t, exists, "Tithi number %d should have a name", i)
		assert.NotEmpty(t, name, "Tithi name for %d should not be empty", i)
	}

	// Test specific important Tithis
	assert.Equal(t, "Pratipada", TithiNames[1])
	assert.Equal(t, "Purnima", TithiNames[15])
	assert.Equal(t, "Pratipada", TithiNames[16]) // Krishna Paksha Pratipada
	assert.Equal(t, "Amavasya", TithiNames[30])
}

func TestCalculateTithiTimes(t *testing.T) {
	calculator, _, _ := createTestTithiCalculator()
	ctx := context.Background()
	referenceDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		tithiFloat   float64
		expectMidday bool // If true, expect reference time to be between start and end
	}{
		{
			name:         "Beginning of Tithi",
			tithiFloat:   7.0, // Exactly at start of 8th Tithi
			expectMidday: false,
		},
		{
			name:         "Middle of Tithi",
			tithiFloat:   7.5, // Middle of 8th Tithi
			expectMidday: true,
		},
		{
			name:         "Near End of Tithi",
			tithiFloat:   7.9, // Near end of 8th Tithi
			expectMidday: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			startTime, endTime := calculator.calculateTithiTimes(ctx, test.tithiFloat, referenceDate)

			// Basic validations
			assert.True(t, endTime.After(startTime))

			duration := endTime.Sub(startTime)
			assert.True(t, duration.Hours() > 20 && duration.Hours() < 30) // Reasonable Tithi duration

			if test.expectMidday {
				noonRef := time.Date(referenceDate.Year(), referenceDate.Month(), referenceDate.Day(), 12, 0, 0, 0, referenceDate.Location())
				assert.True(t, noonRef.After(startTime) && noonRef.Before(endTime),
					"Reference time should be between start and end for middle of Tithi")
			}
		})
	}
}

// Benchmark tests
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

func BenchmarkGetTithiType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getTithiType((i % 30) + 1)
	}
}

// Edge case tests
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
			sunLong:  720.0, // Multiple full circles
			moonLong: 800.0,
		},
		{
			name:     "Negative longitude (should be normalized)",
			sunLong:  350.0,
			moonLong: -10.0, // Should normalize to 350.0
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tithi, err := calculator.calculateTithiFromLongitudes(ctx, test.sunLong, test.moonLong, referenceDate, "traditional")

			require.NoError(t, err)
			require.NotNil(t, tithi)

			// Should always produce valid results
			err = ValidateTithiCalculation(tithi)
			assert.NoError(t, err)
		})
	}
}
