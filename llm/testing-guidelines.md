# Testing Guidelines

This document outlines testing requirements, best practices, and standards for the Panchangam project. Maintaining high test coverage ensures code reliability and prevents regressions.

## Testing Requirements

### Code Coverage Standards

**Critical Requirement**: All pull requests must maintain **minimum 90% code coverage**.

Coverage metrics:
- **Line Coverage**: 90% minimum
- **Function Coverage**: 90% minimum
- **Branch Coverage**: 85% minimum (recommended)

### Test Types

1. **Unit Tests**: Test individual functions and methods in isolation
2. **Integration Tests**: Test interactions between components
3. **End-to-End Tests**: Test complete user workflows (frontend)
4. **Performance Tests**: Test calculation performance and optimization

## Go Backend Testing

### Testing Framework

- **Standard Library**: `testing` package
- **Assertions**: `github.com/stretchr/testify/assert` and `testify/require`
- **Mocking**: `testify/mock` or interfaces for dependency injection

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./astronomy/...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

### Test File Structure

```go
package astronomy

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Table-driven tests are preferred
func TestCalculateTithi(t *testing.T) {
    tests := []struct {
        name        string
        sunLong     float64
        moonLong    float64
        want        int
        wantErr     bool
        errContains string
    }{
        {
            name:     "new moon - tithi 1",
            sunLong:  0,
            moonLong: 0,
            want:     1,
            wantErr:  false,
        },
        {
            name:     "full moon - tithi 15",
            sunLong:  0,
            moonLong: 180,
            want:     15,
            wantErr:  false,
        },
        {
            name:        "invalid longitude",
            sunLong:     -10,
            moonLong:    0,
            want:        0,
            wantErr:     true,
            errContains: "invalid longitude",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CalculateTithi(tt.sunLong, tt.moonLong)

            if tt.wantErr {
                require.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Testing Best Practices - Go

#### 1. Use Table-Driven Tests

```go
func TestNakshatraCalculation(t *testing.T) {
    tests := []struct {
        name      string
        moonLong  float64
        want      int
        wantName  string
    }{
        {"Ashwini start", 0.0, 1, "Ashwini"},
        {"Bharani start", 13.333, 2, "Bharani"},
        {"Revati end", 359.999, 27, "Revati"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            nakshatra, name := CalculateNakshatra(tt.moonLong)
            assert.Equal(t, tt.want, nakshatra)
            assert.Equal(t, tt.wantName, name)
        })
    }
}
```

#### 2. Test Edge Cases

```go
func TestSunriseCalculation(t *testing.T) {
    tests := []struct {
        name     string
        lat      float64
        lon      float64
        date     time.Time
        wantErr  bool
    }{
        {
            name: "normal case - Mumbai",
            lat:  19.0760,
            lon:  72.8777,
            date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        {
            name:    "polar region - midnight sun",
            lat:     80.0,
            lon:     0.0,
            date:    time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
            wantErr: true,
        },
        {
            name:    "invalid latitude",
            lat:     100.0,
            lon:     0.0,
            date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := CalculateSunrise(tt.lat, tt.lon, tt.date)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### 3. Use Mocks for External Dependencies

```go
// Define interface for dependency
type EphemerisProvider interface {
    GetPlanetPosition(planet string, date time.Time) (float64, error)
}

// Create mock
type MockEphemerisProvider struct {
    mock.Mock
}

func (m *MockEphemerisProvider) GetPlanetPosition(planet string, date time.Time) (float64, error) {
    args := m.Called(planet, date)
    return args.Get(0).(float64), args.Error(1)
}

// Test with mock
func TestPanchangamCalculation(t *testing.T) {
    mockProvider := new(MockEphemerisProvider)
    date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

    // Set expectations
    mockProvider.On("GetPlanetPosition", "Sun", date).Return(280.0, nil)
    mockProvider.On("GetPlanetPosition", "Moon", date).Return(320.0, nil)

    service := NewPanchangamService(mockProvider)
    result, err := service.CalculateDaily(date)

    require.NoError(t, err)
    assert.NotNil(t, result)

    // Verify all expectations were met
    mockProvider.AssertExpectations(t)
}
```

#### 4. Test Concurrency

```go
func TestConcurrentCalculations(t *testing.T) {
    dates := []time.Time{
        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
    }

    results := make(chan *PanchangamData, len(dates))

    for _, date := range dates {
        go func(d time.Time) {
            data, err := CalculatePanchangam(d)
            require.NoError(t, err)
            results <- data
        }(date)
    }

    // Collect results
    for i := 0; i < len(dates); i++ {
        result := <-results
        assert.NotNil(t, result)
    }
}
```

#### 5. Benchmark Tests

```go
func BenchmarkTithiCalculation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        CalculateTithi(123.45, 234.56)
    }
}

func BenchmarkPanchangamGeneration(b *testing.B) {
    date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    location := Location{Lat: 19.0760, Lon: 72.8777}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        CalculatePanchangam(date, location)
    }
}
```

### Coverage Analysis

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
go test -cover ./...
```

## TypeScript Frontend Testing

### Testing Framework

- **Test Runner**: Vitest
- **Testing Library**: @testing-library/react
- **User Interactions**: @testing-library/user-event
- **E2E Testing**: Playwright

### Running Tests

```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in UI mode
npm run test:ui

# Run tests in watch mode
npm test -- --watch

# Run specific test file
npm test PanchangamDisplay.test.tsx
```

### Test File Structure

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { PanchangamDisplay } from './PanchangamDisplay';

describe('PanchangamDisplay', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should render Panchangam data correctly', async () => {
        const mockData = {
            tithi: { number: 1, name: 'Pratipada', progress: 0.5 },
            nakshatra: { number: 1, name: 'Ashwini' },
        };

        // Mock API call
        vi.spyOn(global, 'fetch').mockResolvedValueOnce({
            ok: true,
            json: async () => mockData,
        } as Response);

        render(<PanchangamDisplay date={new Date('2024-01-01')} />);

        await waitFor(() => {
            expect(screen.getByText('Pratipada')).toBeInTheDocument();
            expect(screen.getByText('Ashwini')).toBeInTheDocument();
        });
    });

    it('should handle errors gracefully', async () => {
        vi.spyOn(global, 'fetch').mockRejectedValueOnce(
            new Error('Network error')
        );

        render(<PanchangamDisplay date={new Date('2024-01-01')} />);

        await waitFor(() => {
            expect(screen.getByText(/error/i)).toBeInTheDocument();
        });
    });
});
```

### Testing Best Practices - TypeScript

#### 1. Test User Interactions

```typescript
it('should change date when user clicks next button', async () => {
    const user = userEvent.setup();
    const onDateChange = vi.fn();

    render(
        <DateNavigator
            date={new Date('2024-01-01')}
            onDateChange={onDateChange}
        />
    );

    const nextButton = screen.getByRole('button', { name: /next/i });
    await user.click(nextButton);

    expect(onDateChange).toHaveBeenCalledWith(new Date('2024-01-02'));
});
```

#### 2. Test Async Operations

```typescript
it('should load Panchangam data on mount', async () => {
    const mockData = { tithi: { name: 'Pratipada' } };

    vi.spyOn(global, 'fetch').mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
    } as Response);

    render(<PanchangamDisplay date={new Date('2024-01-01')} />);

    // Initially shows loading
    expect(screen.getByText(/loading/i)).toBeInTheDocument();

    // Wait for data to load
    await waitFor(() => {
        expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
        expect(screen.getByText('Pratipada')).toBeInTheDocument();
    });
});
```

#### 3. Test Error States

```typescript
it('should show error message when API fails', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    vi.spyOn(global, 'fetch').mockRejectedValueOnce(
        new Error('API Error')
    );

    render(<PanchangamDisplay date={new Date('2024-01-01')} />);

    await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(/error/i);
    });

    consoleErrorSpy.mockRestore();
});
```

#### 4. Test Custom Hooks

```typescript
import { renderHook, waitFor } from '@testing-library/react';

describe('usePanchangamData', () => {
    it('should fetch data on mount', async () => {
        const mockData = { tithi: { name: 'Pratipada' } };

        vi.spyOn(global, 'fetch').mockResolvedValueOnce({
            ok: true,
            json: async () => mockData,
        } as Response);

        const { result } = renderHook(() =>
            usePanchangamData(new Date('2024-01-01'))
        );

        expect(result.current.loading).toBe(true);

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
            expect(result.current.data).toEqual(mockData);
            expect(result.current.error).toBeNull();
        });
    });
});
```

#### 5. Snapshot Testing (Use Sparingly)

```typescript
it('should match snapshot', () => {
    const { container } = render(
        <TithiCard tithi={{ number: 1, name: 'Pratipada', progress: 0.5 }} />
    );

    expect(container).toMatchSnapshot();
});
```

### Coverage Configuration

Add to `vitest.config.ts`:

```typescript
export default defineConfig({
    test: {
        coverage: {
            provider: 'v8',
            reporter: ['text', 'json', 'html'],
            exclude: [
                'node_modules/',
                'src/**/*.test.{ts,tsx}',
                'src/**/*.types.ts',
                '**/*.d.ts',
            ],
            statements: 90,
            branches: 85,
            functions: 90,
            lines: 90,
        },
    },
});
```

## Integration Testing

### Backend Integration Tests

```go
func TestPanchangamServiceIntegration(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup
    ephemeris := ephemeris.NewSwissEphemeris()
    service := NewPanchangamService(ephemeris)

    // Test
    date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    location := Location{Lat: 19.0760, Lon: 72.8777}

    result, err := service.CalculateDaily(date, location)

    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, date, result.Date)

    // Validate data ranges
    assert.GreaterOrEqual(t, result.Tithi.Number, 1)
    assert.LessOrEqual(t, result.Tithi.Number, 30)
}
```

### Frontend E2E Tests (Playwright)

```typescript
import { test, expect } from '@playwright/test';

test('should display Panchangam for selected date', async ({ page }) => {
    await page.goto('/');

    // Wait for initial load
    await expect(page.locator('[data-testid="panchangam-display"]')).toBeVisible();

    // Check that Tithi is displayed
    await expect(page.locator('[data-testid="tithi-name"]')).toContainText(/\w+/);

    // Change date
    await page.click('[data-testid="date-next"]');

    // Verify data updated
    await expect(page.locator('[data-testid="loading"]')).toBeVisible();
    await expect(page.locator('[data-testid="loading"]')).not.toBeVisible();
});
```

## Test Coverage Best Practices

### What to Test

✅ **DO Test:**
- All public functions and methods
- Edge cases and boundary conditions
- Error handling and validation
- Business logic and calculations
- User interactions and workflows
- Async operations and promises
- State management

❌ **DON'T Test:**
- Third-party library internals
- Simple getters/setters without logic
- Configuration files
- Type definitions (TypeScript handles this)

### Achieving 90% Coverage

1. **Write Tests First**: TDD approach helps reach coverage goals
2. **Test Edge Cases**: Cover error paths and boundary conditions
3. **Use Coverage Reports**: Identify untested code
4. **Avoid Test Duplication**: Don't test the same thing multiple ways
5. **Test Behavior, Not Implementation**: Focus on what, not how

### Coverage Gaps

If coverage falls below 90%:

1. Run coverage report to identify gaps
2. Add tests for uncovered lines
3. Consider if code is testable (may need refactoring)
4. Document any intentional exclusions

## Continuous Integration

### Pre-commit Checks

- Run tests locally before committing
- Ensure all tests pass
- Verify coverage meets threshold
- Run linters and formatters

### CI Pipeline

```yaml
# Example GitHub Actions workflow
- name: Run Go tests
  run: make test-coverage

- name: Run Frontend tests
  run: npm run test:coverage

- name: Check coverage threshold
  run: |
    go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | awk '{if ($1 < 90) exit 1}'
```

## Testing Checklist

Before submitting a PR:

- [ ] All tests pass locally
- [ ] New code has corresponding tests
- [ ] Coverage meets 90% threshold
- [ ] Edge cases are tested
- [ ] Error scenarios are tested
- [ ] No flaky tests
- [ ] Tests are well-named and clear
- [ ] Mocks are used appropriately
- [ ] Integration tests pass
- [ ] Performance benchmarks run (if applicable)

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Vitest Documentation](https://vitest.dev/)
- [Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Playwright Documentation](https://playwright.dev/)
