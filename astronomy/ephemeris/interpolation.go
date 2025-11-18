package ephemeris

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// InterpolationMethod defines the type of interpolation to use
type InterpolationMethod string

const (
	// InterpolationLinear uses simple linear interpolation
	InterpolationLinear InterpolationMethod = "linear"

	// InterpolationLagrange uses Lagrange polynomial interpolation
	InterpolationLagrange InterpolationMethod = "lagrange"

	// InterpolationCubicSpline uses cubic spline interpolation
	InterpolationCubicSpline InterpolationMethod = "cubic_spline"
)

// InterpolationConfig holds configuration for interpolation operations
type InterpolationConfig struct {
	Method       InterpolationMethod
	Order        int     // Order for polynomial interpolation (3-7 recommended)
	Tolerance    float64 // Maximum acceptable error in degrees
	MaxCacheSize int     // Maximum number of cached data points
}

// DefaultInterpolationConfig returns the default configuration
func DefaultInterpolationConfig() InterpolationConfig {
	return InterpolationConfig{
		Method:       InterpolationCubicSpline,
		Order:        5,
		Tolerance:    0.0001, // ~0.36 arcseconds
		MaxCacheSize: 100,
	}
}

// Interpolator provides interpolation methods for planetary positions
type Interpolator struct {
	manager  *Manager
	config   InterpolationConfig
	observer observability.ObserverInterface
	cache    map[string]*interpolationCache
}

// interpolationCache stores data points for interpolation
type interpolationCache struct {
	jdPoints       []float64
	positionPoints []Position
	maxSize        int
}

// NewInterpolator creates a new interpolator
func NewInterpolator(manager *Manager, config InterpolationConfig) *Interpolator {
	return &Interpolator{
		manager:  manager,
		config:   config,
		observer: observability.Observer(),
		cache:    make(map[string]*interpolationCache),
	}
}

// InterpolatePlanetaryPosition calculates planetary position at a specific JD using interpolation
func (i *Interpolator) InterpolatePlanetaryPosition(ctx context.Context, jd JulianDay, planet string) (*Position, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.InterpolatePlanetaryPosition")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("planet", planet),
		attribute.String("method", string(i.config.Method)),
	)

	// Get surrounding data points
	points, err := i.getDataPoints(ctx, jd, planet)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get data points: %w", err)
	}

	// Perform interpolation based on method
	var position Position
	switch i.config.Method {
	case InterpolationLinear:
		position, err = i.linearInterpolation(points, float64(jd))
	case InterpolationLagrange:
		position, err = i.lagrangeInterpolation(points, float64(jd))
	case InterpolationCubicSpline:
		position, err = i.cubicSplineInterpolation(points, float64(jd))
	default:
		return nil, fmt.Errorf("unsupported interpolation method: %s", i.config.Method)
	}

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Normalize longitude to 0-360 range
	position.Longitude = normalizeAngle(position.Longitude)

	span.SetAttributes(
		attribute.Float64("interpolated_longitude", position.Longitude),
		attribute.Float64("interpolated_latitude", position.Latitude),
		attribute.Bool("success", true),
	)

	return &position, nil
}

// getDataPoints retrieves surrounding data points for interpolation
func (i *Interpolator) getDataPoints(ctx context.Context, jd JulianDay, planet string) ([]dataPoint, error) {
	numPoints := i.config.Order
	if i.config.Method == InterpolationLinear {
		numPoints = 2
	}

	// Calculate surrounding Julian days
	// For better accuracy, center the point around the target JD
	offset := float64(numPoints-1) / 2.0
	startJD := float64(jd) - offset

	points := make([]dataPoint, 0, numPoints)

	for j := 0; j < numPoints; j++ {
		currentJD := JulianDay(startJD + float64(j))

		// Get planetary positions
		positions, err := i.manager.GetPlanetaryPositions(ctx, currentJD)
		if err != nil {
			return nil, fmt.Errorf("failed to get positions for JD %f: %w", currentJD, err)
		}

		// Extract the specific planet position
		pos, err := i.extractPlanetPosition(positions, planet)
		if err != nil {
			return nil, err
		}

		points = append(points, dataPoint{
			jd:       float64(currentJD),
			position: *pos,
		})
	}

	// Sort points by JD
	sort.Slice(points, func(a, b int) bool {
		return points[a].jd < points[b].jd
	})

	return points, nil
}

// dataPoint holds a single data point for interpolation
type dataPoint struct {
	jd       float64
	position Position
}

// extractPlanetPosition extracts the position for a specific planet
func (i *Interpolator) extractPlanetPosition(positions *PlanetaryPositions, planet string) (*Position, error) {
	switch planet {
	case "sun":
		return &positions.Sun, nil
	case "moon":
		return &positions.Moon, nil
	case "mercury":
		return &positions.Mercury, nil
	case "venus":
		return &positions.Venus, nil
	case "mars":
		return &positions.Mars, nil
	case "jupiter":
		return &positions.Jupiter, nil
	case "saturn":
		return &positions.Saturn, nil
	case "uranus":
		return &positions.Uranus, nil
	case "neptune":
		return &positions.Neptune, nil
	case "pluto":
		return &positions.Pluto, nil
	default:
		return nil, fmt.Errorf("unknown planet: %s", planet)
	}
}

// linearInterpolation performs simple linear interpolation
func (i *Interpolator) linearInterpolation(points []dataPoint, jd float64) (Position, error) {
	if len(points) < 2 {
		return Position{}, fmt.Errorf("need at least 2 points for linear interpolation")
	}

	// Find the two points bracketing the target JD
	var p0, p1 dataPoint
	for idx := 0; idx < len(points)-1; idx++ {
		if points[idx].jd <= jd && points[idx+1].jd >= jd {
			p0 = points[idx]
			p1 = points[idx+1]
			break
		}
	}

	if p0.jd == 0 && p1.jd == 0 {
		// Use first two points if target is before all points
		// or last two points if target is after all points
		if jd < points[0].jd {
			p0 = points[0]
			p1 = points[1]
		} else {
			p0 = points[len(points)-2]
			p1 = points[len(points)-1]
		}
	}

	// Calculate interpolation factor
	t := (jd - p0.jd) / (p1.jd - p0.jd)

	// Handle angle wrapping for longitude
	lon0 := p0.position.Longitude
	lon1 := p1.position.Longitude

	// Check if we cross the 0/360 boundary
	if math.Abs(lon1-lon0) > 180 {
		if lon0 > lon1 {
			lon1 += 360
		} else {
			lon0 += 360
		}
	}

	return Position{
		Longitude: lon0 + t*(lon1-lon0),
		Latitude:  p0.position.Latitude + t*(p1.position.Latitude-p0.position.Latitude),
		Distance:  p0.position.Distance + t*(p1.position.Distance-p0.position.Distance),
		Speed:     p0.position.Speed + t*(p1.position.Speed-p0.position.Speed),
	}, nil
}

// lagrangeInterpolation performs Lagrange polynomial interpolation
func (i *Interpolator) lagrangeInterpolation(points []dataPoint, jd float64) (Position, error) {
	n := len(points)
	if n < 2 {
		return Position{}, fmt.Errorf("need at least 2 points for Lagrange interpolation")
	}

	var longitude, latitude, distance, speed float64

	// Lagrange interpolation formula
	for j := 0; j < n; j++ {
		term := 1.0
		for m := 0; m < n; m++ {
			if m != j {
				term *= (jd - points[m].jd) / (points[j].jd - points[m].jd)
			}
		}

		// Handle longitude wrapping
		lon := points[j].position.Longitude
		if j > 0 {
			prevLon := points[j-1].position.Longitude
			if math.Abs(lon-prevLon) > 180 {
				if prevLon > lon {
					lon += 360
				}
			}
		}

		longitude += term * lon
		latitude += term * points[j].position.Latitude
		distance += term * points[j].position.Distance
		speed += term * points[j].position.Speed
	}

	return Position{
		Longitude: longitude,
		Latitude:  latitude,
		Distance:  distance,
		Speed:     speed,
	}, nil
}

// cubicSplineInterpolation performs cubic spline interpolation
func (i *Interpolator) cubicSplineInterpolation(points []dataPoint, jd float64) (Position, error) {
	n := len(points)
	if n < 4 {
		// Fall back to Lagrange for fewer points
		return i.lagrangeInterpolation(points, jd)
	}

	// Build cubic spline for each component
	longitude := i.cubicSplineComponent(points, jd, func(p dataPoint) float64 {
		return p.position.Longitude
	})

	latitude := i.cubicSplineComponent(points, jd, func(p dataPoint) float64 {
		return p.position.Latitude
	})

	distance := i.cubicSplineComponent(points, jd, func(p dataPoint) float64 {
		return p.position.Distance
	})

	speed := i.cubicSplineComponent(points, jd, func(p dataPoint) float64 {
		return p.position.Speed
	})

	return Position{
		Longitude: longitude,
		Latitude:  latitude,
		Distance:  distance,
		Speed:     speed,
	}, nil
}

// cubicSplineComponent performs cubic spline interpolation for a single component
func (i *Interpolator) cubicSplineComponent(points []dataPoint, jd float64, getValue func(dataPoint) float64) float64 {
	n := len(points)

	// Find the interval containing jd
	var idx int
	for idx = 0; idx < n-1; idx++ {
		if jd >= points[idx].jd && jd <= points[idx+1].jd {
			break
		}
	}

	// Ensure idx is valid
	if idx >= n-1 {
		idx = n - 2
	}

	// Extract values for the spline
	x := make([]float64, n)
	y := make([]float64, n)
	for j := 0; j < n; j++ {
		x[j] = points[j].jd
		y[j] = getValue(points[j])
	}

	// Handle angle wrapping for cyclic values (like longitude)
	for j := 1; j < n; j++ {
		if math.Abs(y[j]-y[j-1]) > 180 {
			if y[j-1] > y[j] {
				y[j] += 360
			} else {
				y[j-1] += 360
			}
		}
	}

	// Calculate second derivatives (natural spline boundary conditions)
	h := make([]float64, n-1)
	for j := 0; j < n-1; j++ {
		h[j] = x[j+1] - x[j]
	}

	alpha := make([]float64, n-1)
	for j := 1; j < n-1; j++ {
		alpha[j] = (3.0/h[j])*(y[j+1]-y[j]) - (3.0/h[j-1])*(y[j]-y[j-1])
	}

	l := make([]float64, n)
	mu := make([]float64, n)
	z := make([]float64, n)

	l[0] = 1.0

	for j := 1; j < n-1; j++ {
		l[j] = 2.0*(x[j+1]-x[j-1]) - h[j-1]*mu[j-1]
		mu[j] = h[j] / l[j]
		z[j] = (alpha[j] - h[j-1]*z[j-1]) / l[j]
	}

	l[n-1] = 1.0

	c := make([]float64, n)
	b := make([]float64, n-1)
	d := make([]float64, n-1)

	for j := n - 2; j >= 0; j-- {
		c[j] = z[j] - mu[j]*c[j+1]
		b[j] = (y[j+1]-y[j])/h[j] - h[j]*(c[j+1]+2.0*c[j])/3.0
		d[j] = (c[j+1] - c[j]) / (3.0 * h[j])
	}

	// Evaluate the spline at jd
	dx := jd - x[idx]
	result := y[idx] + b[idx]*dx + c[idx]*dx*dx + d[idx]*dx*dx*dx

	return result
}

// InterpolatePlanetaryPositions interpolates all planetary positions at once
func (i *Interpolator) InterpolatePlanetaryPositions(ctx context.Context, jd JulianDay) (*PlanetaryPositions, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.InterpolatePlanetaryPositions")
	defer span.End()

	span.SetAttributes(
		attribute.Float64("julian_day", float64(jd)),
		attribute.String("method", string(i.config.Method)),
	)

	planets := []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto"}
	positions := &PlanetaryPositions{
		JulianDay: jd,
	}

	for _, planet := range planets {
		pos, err := i.InterpolatePlanetaryPosition(ctx, jd, planet)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to interpolate %s: %w", planet, err)
		}

		// Assign to appropriate field
		switch planet {
		case "sun":
			positions.Sun = *pos
		case "moon":
			positions.Moon = *pos
		case "mercury":
			positions.Mercury = *pos
		case "venus":
			positions.Venus = *pos
		case "mars":
			positions.Mars = *pos
		case "jupiter":
			positions.Jupiter = *pos
		case "saturn":
			positions.Saturn = *pos
		case "uranus":
			positions.Uranus = *pos
		case "neptune":
			positions.Neptune = *pos
		case "pluto":
			positions.Pluto = *pos
		}
	}

	span.SetAttributes(attribute.Bool("success", true))
	return positions, nil
}

// ValidateInterpolation validates interpolation accuracy against direct calculation
func (i *Interpolator) ValidateInterpolation(ctx context.Context, jd JulianDay, planet string) (float64, error) {
	ctx, span := i.observer.CreateSpan(ctx, "interpolator.ValidateInterpolation")
	defer span.End()

	// Get interpolated position
	interpPos, err := i.InterpolatePlanetaryPosition(ctx, jd, planet)
	if err != nil {
		return 0, fmt.Errorf("interpolation failed: %w", err)
	}

	// Get actual position
	actualPositions, err := i.manager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		return 0, fmt.Errorf("failed to get actual positions: %w", err)
	}

	actualPos, err := i.extractPlanetPosition(actualPositions, planet)
	if err != nil {
		return 0, err
	}

	// Calculate error in degrees
	lonError := math.Abs(interpPos.Longitude - actualPos.Longitude)
	if lonError > 180 {
		lonError = 360 - lonError
	}

	latError := math.Abs(interpPos.Latitude - actualPos.Latitude)
	distError := math.Abs(interpPos.Distance-actualPos.Distance) / actualPos.Distance * 100 // percentage

	// Total error (weighted)
	totalError := lonError + latError*0.5 + distError*0.1

	span.SetAttributes(
		attribute.Float64("longitude_error_deg", lonError),
		attribute.Float64("latitude_error_deg", latError),
		attribute.Float64("distance_error_percent", distError),
		attribute.Float64("total_error", totalError),
		attribute.Bool("within_tolerance", totalError <= i.config.Tolerance),
	)

	return totalError, nil
}

// normalizeAngle normalizes an angle to the range [0, 360)
func normalizeAngle(angle float64) float64 {
	result := math.Mod(angle, 360.0)
	if result < 0 {
		result += 360.0
	}
	return result
}

// GetInterpolationMethod returns the current interpolation method
func (i *Interpolator) GetInterpolationMethod() InterpolationMethod {
	return i.config.Method
}

// SetInterpolationMethod sets the interpolation method
func (i *Interpolator) SetInterpolationMethod(method InterpolationMethod) {
	i.config.Method = method
}

// GetInterpolationConfig returns the current configuration
func (i *Interpolator) GetInterpolationConfig() InterpolationConfig {
	return i.config
}
