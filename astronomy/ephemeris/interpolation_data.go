package ephemeris

import (
	"context"
	"fmt"
	"sort"
)

// getDataPoints retrieves surrounding data points for interpolation.
func (i *Interpolator) getDataPoints(ctx context.Context, jd JulianDay, planet string) ([]dataPoint, error) {
	numPoints := i.config.Order
	if i.config.Method == InterpolationLinear {
		numPoints = 2
	}

	offset := float64(numPoints-1) / 2.0
	startJD := float64(jd) - offset
	points := make([]dataPoint, 0, numPoints)

	for j := 0; j < numPoints; j++ {
		currentJD := JulianDay(startJD + float64(j))

		positions, err := i.manager.GetPlanetaryPositions(ctx, currentJD)
		if err != nil {
			return nil, fmt.Errorf("failed to get positions for JD %f: %w", currentJD, err)
		}

		pos, err := i.extractPlanetPosition(positions, planet)
		if err != nil {
			return nil, err
		}

		points = append(points, dataPoint{
			jd:       float64(currentJD),
			position: *pos,
		})
	}

	sort.Slice(points, func(a, b int) bool {
		return points[a].jd < points[b].jd
	})

	return points, nil
}

// extractPlanetPosition extracts the position for a specific planet.
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
