package skyview

import (
	"math"
	"testing"
)

func TestEclipticToEquatorial(t *testing.T) {
	tests := []struct {
		name     string
		ecliptic EclipticCoordinates
		wantRA   float64
		wantDec  float64
		deltaRA  float64
		deltaDec float64
	}{
		{
			name: "Vernal equinox",
			ecliptic: EclipticCoordinates{
				Longitude: 0,
				Latitude:  0,
				Distance:  1,
			},
			wantRA:   0,
			wantDec:  0,
			deltaRA:  1,
			deltaDec: 1,
		},
		{
			name: "Summer solstice",
			ecliptic: EclipticCoordinates{
				Longitude: 90,
				Latitude:  0,
				Distance:  1,
			},
			wantRA:   90,
			wantDec:  23.4,
			deltaRA:  1,
			deltaDec: 1,
		},
		{
			name: "Ecliptic north pole",
			ecliptic: EclipticCoordinates{
				Longitude: 0,
				Latitude:  90,
				Distance:  1,
			},
			wantRA:   270,
			wantDec:  66.6,
			deltaRA:  180,
			deltaDec: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := eclipticToEquatorial(tt.ecliptic)

			if math.Abs(eq.RightAscension-tt.wantRA) > tt.deltaRA {
				t.Errorf("RA = %v, want %v ± %v", eq.RightAscension, tt.wantRA, tt.deltaRA)
			}
			if math.Abs(eq.Declination-tt.wantDec) > tt.deltaDec {
				t.Errorf("Dec = %v, want %v ± %v", eq.Declination, tt.wantDec, tt.deltaDec)
			}
			if eq.Distance != tt.ecliptic.Distance {
				t.Errorf("Distance = %v, want %v", eq.Distance, tt.ecliptic.Distance)
			}
		})
	}
}

func TestEquatorialToHorizontal(t *testing.T) {
	observer := Observer{
		Latitude:  40.7128,
		Longitude: -74.006,
		Altitude:  0,
		Timezone:  "America/New_York",
	}

	tests := []struct {
		name       string
		equatorial EquatorialCoordinates
		lst        float64
	}{
		{
			name: "Object at meridian",
			equatorial: EquatorialCoordinates{
				RightAscension: 180,
				Declination:    40,
				Distance:       1,
			},
			lst: 180,
		},
		{
			name: "Object rising",
			equatorial: EquatorialCoordinates{
				RightAscension: 90,
				Declination:    0,
				Distance:       1,
			},
			lst: 270,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hz := equatorialToHorizontal(tt.equatorial, observer, tt.lst)

			if hz.Azimuth < 0 || hz.Azimuth > 360 {
				t.Errorf("Azimuth = %v, want value in [0, 360]", hz.Azimuth)
			}
			if hz.Altitude < -90 || hz.Altitude > 90 {
				t.Errorf("Altitude = %v, want value in [-90, 90]", hz.Altitude)
			}
			if hz.Distance != tt.equatorial.Distance {
				t.Errorf("Distance = %v, want %v", hz.Distance, tt.equatorial.Distance)
			}
		})
	}
}
