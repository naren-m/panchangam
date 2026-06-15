package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func createTestTithiCalculator() (*TithiCalculator, *MockEphemerisProvider, *MockCache) {
	observability.NewLocalObserver()

	mockProvider := &MockEphemerisProvider{}
	mockCache := &MockCache{}

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
