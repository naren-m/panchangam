# Coding Standards

This document outlines the coding standards and best practices for the Panchangam project. Following these standards ensures code consistency, maintainability, and quality across the codebase.

## General Principles

1. **Write Clear, Self-Documenting Code**: Use descriptive names and clear logic
2. **Avoid Meaningless Duplication**: Extract shared code only when it removes real repeated logic
3. **Keep APIs Focused**: Each public function should have one clear job and clear errors
4. **Keep Functions Small**: Each function should do one thing well
5. **Comment Why, Not What**: Code shows what; comments explain why

## Go Backend Standards

### File Organization

```
package/
├── package.go          # Main package implementation
├── package_test.go     # Unit tests
├── types.go           # Type definitions
├── errors.go          # Error definitions
└── internal/          # Internal implementation details
```

### Naming Conventions

**Packages**
- Use lowercase, single-word names
- Avoid underscores or mixedCaps
- Example: `astronomy`, `ephemeris`, `services`

**Functions/Methods**
- Use camelCase for private functions: `calculateTithi`
- Use PascalCase for exported functions: `NewPanchangamServer`
- Start boolean functions with `Is`, `Has`, or `Can`
- Example: `IsAuspicious`, `HasRetrograde`, `CanCalculate`

**Variables**
- Use camelCase for private variables
- Use PascalCase for exported variables
- Use descriptive names, avoid single letters except for loops
- Example: `sunLongitude`, `moonPosition`, `tithiDuration`

**Constants**
- Use PascalCase or SCREAMING_SNAKE_CASE for exported constants
- Group related constants using `const` blocks
```go
const (
    DegreesPer Tithi = 12.0
    NakshatraCount = 27
    YogaCount = 27
)
```

### Code Structure

**Error Handling**
```go
// Always check errors explicitly
result, err := CalculateSunrise(location, date)
if err != nil {
    return nil, fmt.Errorf("calculate sunrise: %w", err)
}

// Use error wrapping to preserve context
return nil, fmt.Errorf("failed to get ephemeris data: %w", err)
```

**Function Design**
```go
// Good: Clear function signature with named return values for complex returns
func CalculateTithi(sunLong, moonLong float64) (tithi int, progress float64, err error) {
    if sunLong < 0 || sunLong > 360 {
        return 0, 0, ErrInvalidLongitude
    }

    diff := moonLong - sunLong
    if diff < 0 {
        diff += 360
    }

    tithi = int(diff/DegreesPerTithi) + 1
    progress = math.Mod(diff, DegreesPerTithi) / DegreesPerTithi

    return tithi, progress, nil
}
```

**Context Usage**
```go
// Always pass context as the first parameter
func FetchEphemerisData(ctx context.Context, date time.Time) (*EphemerisData, error) {
    // Check context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Your implementation
}
```

**Interfaces**
```go
// Add interfaces only when a caller needs them.
// Keep interfaces small and local to the caller.
type EphemerisProvider interface {
    GetPlanetaryPosition(ctx context.Context, planet string, date time.Time) (Position, error)
}
```

### Documentation

**Package Documentation**
```go
// Package astronomy provides astronomical calculations for Panchangam.
//
// This package implements calculations for Tithi, Nakshatra, Yoga, Karana,
// and other Panchangam elements based on Swiss Ephemeris data.
package astronomy
```

**Function Documentation**
```go
// CalculateSunrise computes the sunrise time for a given location and date.
// It accounts for atmospheric refraction and the observer's elevation.
//
// Parameters:
//   - location: Geographic coordinates (latitude, longitude, elevation)
//   - date: The date for which to calculate sunrise
//
// Returns the sunrise time in local solar time and any error encountered.
func CalculateSunrise(location Location, date time.Time) (time.Time, error) {
    // Implementation
}
```

### Testing

```go
// Use table-driven tests
func TestCalculateTithi(t *testing.T) {
    tests := []struct {
        name        string
        sunLong     float64
        moonLong    float64
        wantTithi   int
        wantProgress float64
        wantErr     bool
    }{
        {
            name:        "New Moon",
            sunLong:     0,
            moonLong:    0,
            wantTithi:   1,
            wantProgress: 0,
            wantErr:     false,
        },
        // More test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tithi, progress, err := CalculateTithi(tt.sunLong, tt.moonLong)

            if (err != nil) != tt.wantErr {
                t.Errorf("CalculateTithi() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if tithi != tt.wantTithi {
                t.Errorf("CalculateTithi() tithi = %v, want %v", tithi, tt.wantTithi)
            }

            if math.Abs(progress-tt.wantProgress) > 0.0001 {
                t.Errorf("CalculateTithi() progress = %v, want %v", progress, tt.wantProgress)
            }
        })
    }
}
```

### Service Construction

Use direct constructors. Add options or interfaces only after there is a real caller that needs them.

```go
server := NewPanchangamServer()
response, err := server.Get(ctx, req)
```

## TypeScript Frontend Standards

### File Organization

```
components/
├── Calendar/
│   ├── MonthNavigation.tsx
│   ├── CalendarDisplayManager.tsx
│   └── __tests__/
├── Settings/
│   ├── SettingsPanel.tsx
│   ├── ApiHealthCheck.tsx
│   └── index.ts
```

### Naming Conventions

**Files**
- Use PascalCase for component files: `MonthNavigation.tsx`
- Use camelCase for utilities: `formatDate.ts`
- Use kebab-case for standalone style files when needed

**Components**
- Use PascalCase: `MonthNavigation`, `SettingsPanel`, `ApiHealthCheck`
- Hooks start with `use`: `usePanchangam`, `useProgressivePanchangam`

**Variables/Functions**
- Use camelCase: `tithiData`, `calculateProgress`, `formatTime`
- Use SCREAMING_SNAKE_CASE for constants: `MAX_RETRIES`, `API_BASE_URL`

### TypeScript Best Practices

**Type Definitions**
```typescript
// Define explicit types for all function parameters and return values
interface PanchangamData {
    date: string;
    tithi: string;
    nakshatra: string;
    yoga: string;
    karana: string;
    sunrise_time: string;
    sunset_time: string;
    events: Event[];
}

// Use type instead of interface for unions and intersections
type PanchangamElement = 'tithi' | 'nakshatra' | 'yoga' | 'karana';

// Use proper return types
function loadPanchangamData(date: Date, settings: Settings): Promise<PanchangamData> {
    // Implementation
}
```

**Component Structure**
```typescript
// Use functional components with TypeScript
interface MonthNavigationProps {
    year: number;
    month: number;
    settings: Settings;
    onPrevMonth: () => void;
    onNextMonth: () => void;
}

export function MonthNavigation({
    year,
    month,
    settings,
    onPrevMonth,
    onNextMonth,
}: MonthNavigationProps): JSX.Element {
    const monthName = getMonthName(month, settings.locale);

    return (
        // JSX
    );
}
```

**Hooks**
```typescript
// Custom hooks should start with 'use'
function usePanchangam(date: Date, settings: Settings) {
    const [data, setData] = useState<PanchangamData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        let cancelled = false;

        async function fetchData() {
            try {
                setLoading(true);
                const response = await panchangamApiClient.getPanchangam({
                    date: formatDateForApi(date),
                    latitude: settings.location.latitude,
                    longitude: settings.location.longitude,
                    timezone: settings.location.timezone,
                    region: settings.region,
                    calculation_method: settings.calculation_method,
                    locale: settings.locale,
                });

                if (!cancelled) {
                    setData(response);
                    setError(null);
                }
            } catch (err) {
                if (!cancelled) {
                    setError(err as Error);
                }
            } finally {
                if (!cancelled) {
                    setLoading(false);
                }
            }
        }

        fetchData();

        return () => {
            cancelled = true;
        };
    }, [date, settings]);

    return { data, loading, error };
}
```

**Error Handling**
```typescript
// Always handle API errors and keep messages clear
async function loadPanchangamData(date: Date, settings: Settings): Promise<PanchangamData> {
    try {
        return await panchangamApiClient.getPanchangam({
            date: formatDateForApi(date),
            latitude: settings.location.latitude,
            longitude: settings.location.longitude,
            timezone: settings.location.timezone,
            region: settings.region,
            calculation_method: settings.calculation_method,
            locale: settings.locale,
        });
    } catch (error) {
        console.error('Failed to load Panchangam data:', error);
        throw error;
    }
}
```

### React Best Practices

**Component Design**
- Keep components small and focused
- Extract reusable logic into custom hooks
- Use composition over inheritance
- Memoize expensive calculations with `useMemo`
- Memoize callbacks with `useCallback` when passed to child components

**Props**
- Use destructuring for props
- Define prop types with TypeScript interfaces
- Mark optional props with `?`
- Provide default values where appropriate

**State Management**
- Keep state as local as possible
- Lift state up only when necessary
- Use Context for truly global state
- Consider reducer pattern for complex state logic

### Code Quality

**ESLint Configuration**
- Follow the project's ESLint rules
- No unused variables or imports
- Use `const` over `let` when possible
- Prefer arrow functions for consistency

**Import Organization**
```typescript
// 1. External dependencies
import React, { useState, useEffect } from 'react';

// 2. Internal modules
import { usePanchangam } from '@/hooks/usePanchangam';
import { panchangamApiClient } from '@/services/api/panchangamApiClient';
import { formatDateForApi, getMonthName } from '@/utils/dateHelpers';

// 3. Types
import type { PanchangamData, Settings } from '@/types/panchangam';
```

## Security Considerations

### Input Validation
- Always validate user input on both frontend and backend
- Sanitize data before use
- Use parameterized queries for database operations
- Validate date ranges and geographic coordinates

### API Security
- Use CORS appropriately
- Keep allowed origins explicit
- Return clear validation errors
- Log request IDs with errors

### Error Messages
- Don't expose internal implementation details
- Use generic error messages for users
- Log detailed errors for debugging

## Performance Guidelines

### Backend
- Use context for request cancellation
- Implement caching for expensive calculations
- Reuse cache and gRPC clients where possible
- Use appropriate buffer sizes

### Frontend
- Lazy load components
- Optimize images and assets
- Minimize bundle size
- Use code splitting
- Implement virtual scrolling for large lists

## Code Review Checklist

Before submitting code:
- [ ] Follows naming conventions
- [ ] Has appropriate error handling
- [ ] Includes unit tests with 90%+ coverage
- [ ] Has clear documentation
- [ ] No hardcoded values (use constants)
- [ ] No console.log or debug statements
- [ ] Handles edge cases
- [ ] Uses appropriate types
- [ ] Follows project structure
- [ ] Passes all linters and formatters
