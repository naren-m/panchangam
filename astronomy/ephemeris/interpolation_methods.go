package ephemeris

import (
	"fmt"
	"math"
)

// linearInterpolation performs simple linear interpolation.
func (i *Interpolator) linearInterpolation(points []dataPoint, jd float64) (Position, error) {
	if len(points) < 2 {
		return Position{}, fmt.Errorf("need at least 2 points for linear interpolation")
	}

	var p0, p1 dataPoint
	for idx := 0; idx < len(points)-1; idx++ {
		if points[idx].jd <= jd && points[idx+1].jd >= jd {
			p0 = points[idx]
			p1 = points[idx+1]
			break
		}
	}

	if p0.jd == 0 && p1.jd == 0 {
		if jd < points[0].jd {
			p0 = points[0]
			p1 = points[1]
		} else {
			p0 = points[len(points)-2]
			p1 = points[len(points)-1]
		}
	}

	t := (jd - p0.jd) / (p1.jd - p0.jd)
	lon0 := p0.position.Longitude
	lon1 := p1.position.Longitude

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

// lagrangeInterpolation performs Lagrange polynomial interpolation.
func (i *Interpolator) lagrangeInterpolation(points []dataPoint, jd float64) (Position, error) {
	n := len(points)
	if n < 2 {
		return Position{}, fmt.Errorf("need at least 2 points for Lagrange interpolation")
	}

	var longitude, latitude, distance, speed float64

	for j := 0; j < n; j++ {
		term := 1.0
		for m := 0; m < n; m++ {
			if m != j {
				term *= (jd - points[m].jd) / (points[j].jd - points[m].jd)
			}
		}

		lon := points[j].position.Longitude
		if j > 0 {
			prevLon := points[j-1].position.Longitude
			if math.Abs(lon-prevLon) > 180 && prevLon > lon {
				lon += 360
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

// cubicSplineInterpolation performs cubic spline interpolation.
func (i *Interpolator) cubicSplineInterpolation(points []dataPoint, jd float64) (Position, error) {
	n := len(points)
	if n < 4 {
		return i.lagrangeInterpolation(points, jd)
	}

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

// cubicSplineComponent performs cubic spline interpolation for a single component.
func (i *Interpolator) cubicSplineComponent(points []dataPoint, jd float64, getValue func(dataPoint) float64) float64 {
	n := len(points)

	var idx int
	for idx = 0; idx < n-1; idx++ {
		if jd >= points[idx].jd && jd <= points[idx+1].jd {
			break
		}
	}
	if idx >= n-1 {
		idx = n - 2
	}

	x := make([]float64, n)
	y := make([]float64, n)
	for j := 0; j < n; j++ {
		x[j] = points[j].jd
		y[j] = getValue(points[j])
	}

	for j := 1; j < n; j++ {
		if math.Abs(y[j]-y[j-1]) > 180 {
			if y[j-1] > y[j] {
				y[j] += 360
			} else {
				y[j-1] += 360
			}
		}
	}

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

	dx := jd - x[idx]
	return y[idx] + b[idx]*dx + c[idx]*dx*dx + d[idx]*dx*dx*dx
}

// normalizeAngle normalizes an angle to the range [0, 360).
func normalizeAngle(angle float64) float64 {
	result := math.Mod(angle, 360.0)
	if result < 0 {
		result += 360.0
	}
	return result
}
