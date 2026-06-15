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
- **Test doubles**: Use simple fakes only when a real dependency makes the test slow or flaky

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

#### 3. Test the Current Service API

Service tests should call `Panchangam.Get` through `NewPanchangamServer()`.

```go
func TestPanchangamGet(t *testing.T) {
    ctx := context.Background()
    server := NewPanchangamServer()
    req := &ppb.GetPanchangamRequest{
        Date:      "2024-01-01",
        Latitude:  19.0760,
        Longitude: 72.8777,
        Timezone:  "Asia/Kolkata",
        Region:    "IN",
        Locale:    "en",
    }

    resp, err := server.Get(ctx, req)

    require.NoError(t, err)
    require.NotNil(t, resp)
    require.NotNil(t, resp.PanchangamData)
    assert.NotEmpty(t, resp.PanchangamData.Tithi)
}
```

#### 4. Test Concurrency

```go
func TestConcurrentCalculations(t *testing.T) {
    ctx := context.Background()
    server := NewPanchangamServer()
    dates := []time.Time{
        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
        time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
    }

    results := make(chan *ppb.PanchangamData, len(dates))

    for _, date := range dates {
        go func(d time.Time) {
            req := &ppb.GetPanchangamRequest{
                Date:      d.Format("2006-01-02"),
                Latitude:  19.0760,
                Longitude: 72.8777,
                Timezone:  "Asia/Kolkata",
                Region:    "IN",
                Locale:    "en",
            }
            resp, err := server.Get(ctx, req)
            require.NoError(t, err)
            results <- resp.PanchangamData
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
    ctx := context.Background()
    server := NewPanchangamServer()
    req := &ppb.GetPanchangamRequest{
        Date:      "2024-01-01",
        Latitude:  19.0760,
        Longitude: 72.8777,
        Timezone:  "Asia/Kolkata",
        Region:    "IN",
        Locale:    "en",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = server.Get(ctx, req)
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
npm test MonthNavigation.test.tsx
```

### Test File Structure

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MonthNavigation } from '../Calendar/MonthNavigation';
import type { Settings } from '../../types/panchangam';

describe('MonthNavigation', () => {
    const onPrevMonth = vi.fn();
    const onNextMonth = vi.fn();

    const settings: Settings = {
        calculation_method: 'Drik',
        locale: 'en',
        region: 'Karnataka',
        time_format: '12',
        location: {
            latitude: 12.9716,
            longitude: 77.5946,
            timezone: 'Asia/Kolkata',
            name: 'Bangalore',
            region: 'Karnataka',
        },
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders the selected month and location', () => {
        render(
            <MonthNavigation
                year={2024}
                month={0}
                settings={settings}
                onPrevMonth={onPrevMonth}
                onNextMonth={onNextMonth}
                onToday={vi.fn()}
                onLocationClick={vi.fn()}
                onSettingsClick={vi.fn()}
            />
        );

        expect(screen.getByText(/January/i)).toBeInTheDocument();
        expect(screen.getByText(/Bangalore/i)).toBeInTheDocument();
    });

    it('calls navigation handlers', () => {
        render(
            <MonthNavigation
                year={2024}
                month={0}
                settings={settings}
                onPrevMonth={onPrevMonth}
                onNextMonth={onNextMonth}
                onToday={vi.fn()}
                onLocationClick={vi.fn()}
                onSettingsClick={vi.fn()}
            />
        );

        fireEvent.click(screen.getByRole('button', { name: /previous/i }));
        fireEvent.click(screen.getByRole('button', { name: /next/i }));

        expect(onPrevMonth).toHaveBeenCalledTimes(1);
        expect(onNextMonth).toHaveBeenCalledTimes(1);
    });
});
```

### Testing Best Practices - TypeScript

#### 1. Test User Interactions

```typescript
it('should change date when user clicks next button', async () => {
    const onNextMonth = vi.fn();

    render(
        <MonthNavigation
            year={2024}
            month={0}
            settings={settings}
            onPrevMonth={vi.fn()}
            onNextMonth={onNextMonth}
            onToday={vi.fn()}
            onLocationClick={vi.fn()}
            onSettingsClick={vi.fn()}
        />
    );

    const nextButton = screen.getByRole('button', { name: /next/i });
    fireEvent.click(nextButton);

    expect(onNextMonth).toHaveBeenCalledTimes(1);
});
```

#### 2. Test Async Operations

```typescript
it('should load Panchangam data', async () => {
    panchangamApiClient.getPanchangam.mockResolvedValue({
        date: '2024-01-15',
        tithi: 'Panchami',
        nakshatra: 'Rohini',
        yoga: 'Vishkumbha',
        karana: 'Bava',
        sunrise_time: '06:30:00',
        sunset_time: '18:15:00',
        events: [],
    });

    const { result } = renderHook(() => usePanchangam(date, settings));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
        expect(result.current.loading).toBe(false);
    });

    expect(result.current.data?.tithi).toBe('Panchami');
});
```

#### 3. Test Error States

```typescript
it('should show error message when API fails', async () => {
    panchangamApiClient.getPanchangam.mockRejectedValue(new Error('API Error'));

    const { result } = renderHook(() => usePanchangam(date, settings));

    await waitFor(() => {
        expect(result.current.loading).toBe(false);
    });

    expect(result.current.errorState.hasError).toBe(true);
    expect(result.current.errorState.message).toBe('API Error');
});
```

#### 4. Test Custom Hooks

```typescript
import { renderHook, waitFor } from '@testing-library/react';

describe('usePanchangam', () => {
    it('should fetch data on mount', async () => {
        const data = { date: '2024-01-15', tithi: 'Panchami' };
        panchangamApiClient.getPanchangam.mockResolvedValue(data);

        const { result } = renderHook(() =>
            usePanchangam(new Date('2024-01-15'), settings)
        );

        expect(result.current.loading).toBe(true);

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
            expect(result.current.data).toEqual(data);
            expect(result.current.error).toBeNull();
        });
    });
});
```

#### 5. Keep Component Tests Focused

```typescript
it('renders settings choices', () => {
    render(<SettingsPanel settings={settings} onSettingsChange={vi.fn()} onClose={vi.fn()} />);

    expect(screen.getByText('Calculation Method')).toBeInTheDocument();
    expect(screen.getByText('Language')).toBeInTheDocument();
});
```

### Coverage Configuration

Keep coverage settings in `ui/vite.config.ts` under `test.coverage`:

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
func TestPanchangamGetIntegration(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    ctx := context.Background()
    server := NewPanchangamServer()
    req := &ppb.GetPanchangamRequest{
        Date:      "2024-01-01",
        Latitude:  19.0760,
        Longitude: 72.8777,
        Timezone:  "Asia/Kolkata",
        Region:    "IN",
        Locale:    "en",
    }

    result, err := server.Get(ctx, req)

    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotNil(t, result.PanchangamData)

    assert.NotEmpty(t, result.PanchangamData.Date)
    assert.NotEmpty(t, result.PanchangamData.Tithi)
    assert.NotEmpty(t, result.PanchangamData.Nakshatra)
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

PASS **DO Test:**
- All public functions and methods
- Edge cases and boundary conditions
- Error handling and validation
- Business logic and calculations
- User interactions and workflows
- Async operations and promises
- State management

FAIL **DON'T Test:**
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
- [ ] Test doubles are simple and only used when needed
- [ ] Integration tests pass
- [ ] Performance benchmarks run (if applicable)

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Vitest Documentation](https://vitest.dev/)
- [Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Playwright Documentation](https://playwright.dev/)
