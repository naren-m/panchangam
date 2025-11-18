# Astronomical Algorithms Documentation

This document provides comprehensive documentation of all astronomical calculation algorithms used in the Panchangam project, including mathematical formulas, references, and implementation details.

## Table of Contents

1. [Overview](#overview)
2. [Coordinate Systems](#coordinate-systems)
3. [Ephemeris Calculations](#ephemeris-calculations)
4. [Panchangam Elements](#panchangam-elements)
5. [Interpolation Methods](#interpolation-methods)
6. [Retrograde Motion Detection](#retrograde-motion-detection)
7. [References](#references)

## Overview

The Panchangam project implements astronomical calculations based on established algorithms from astronomical references and traditional Indian astronomical texts (Jyotisha). All calculations use high-precision ephemeris data from JPL DE440 and Swiss Ephemeris.

## Coordinate Systems

### Ecliptic Coordinates

The ecliptic coordinate system uses the plane of Earth's orbit around the Sun as its fundamental plane.

**Longitude (λ)**: Angular distance along the ecliptic from the vernal equinox (0° to 360°)
**Latitude (β)**: Angular distance perpendicular to the ecliptic (-90° to +90°)

### Equatorial Coordinates

The equatorial coordinate system uses Earth's equator as its fundamental plane.

**Right Ascension (α)**: Angular distance along the celestial equator from the vernal equinox (0h to 24h or 0° to 360°)
**Declination (δ)**: Angular distance perpendicular to the celestial equator (-90° to +90°)

### Transformation: Ecliptic to Equatorial

```
α = atan2(sin(λ)cos(ε) - tan(β)sin(ε), cos(λ))
δ = asin(sin(β)cos(ε) + cos(β)sin(ε)sin(λ))
```

Where ε is the obliquity of the ecliptic (approximately 23.44°)

## Ephemeris Calculations

### Julian Day Number

The Julian Day Number (JD) is a continuous count of days since noon Universal Time on January 1, 4713 BCE.

**Formula:**
```
a = floor((14 - month) / 12)
y = year + 4800 - a
m = month + 12 * a - 3

JD = day + floor((153 * m + 2) / 5) + 365 * y + floor(y / 4) - floor(y / 100) + floor(y / 400) - 32045
```

For time of day:
```
JD_final = JD + (hour - 12) / 24 + minute / 1440 + second / 86400
```

**Implementation:** `astronomy/ephemeris/ephemeris.go:TimeToJulianDay()`

**Reference:** Meeus, J. (1998). *Astronomical Algorithms*, 2nd ed.

### Planetary Positions

Planetary positions are calculated using high-precision ephemeris:

1. **JPL DE440**: Primary source, valid from 1550 CE to 2650 CE
   - Accuracy: ±0.0001° for inner planets, ±0.001° for outer planets
   - Based on numerical integration of planetary motion

2. **Swiss Ephemeris**: Fallback source, valid from 13201 BCE to 17191 CE
   - Accuracy: ±0.001° to ±0.01° depending on time period
   - Based on analytical theories and numerical integration

**Implementation:** `astronomy/ephemeris/jpl_provider.go`, `astronomy/ephemeris/swiss_provider.go`

## Panchangam Elements

The five elements (Pancha Angas) of the Panchangam are:

### 1. Tithi (Lunar Day)

A Tithi is 1/30th of a lunar month, representing the angular separation between Sun and Moon.

**Formula:**
```
Tithi_angle = Moon_longitude - Sun_longitude (mod 360°)
Tithi_number = floor(Tithi_angle / 12°) + 1
```

**Tithis (30 total):**
- Shukla Paksha (Waxing): 1-15 (Pratipada to Purnima)
- Krishna Paksha (Waning): 1-15 (Pratipada to Amavasya)

**Implementation:** `astronomy/tithi.go:Calculate()`

**Mathematical Details:**
- Each Tithi spans 12° of elongation
- Tithi transitions occur at variable times (not fixed to midnight)
- Start and end times calculated by solving: `Δλ = Sun_long - Moon_long = n × 12°`

**Reference:**
- Burgess, E. (1858). *Surya Siddhanta*
- Sewell, R., & Dikshit, S. B. (1896). *The Indian Calendar*

### 2. Nakshatra (Lunar Mansion)

A Nakshatra is 1/27th of the zodiac, based on the Moon's position.

**Formula:**
```
Nakshatra_number = floor(Moon_longitude / 13.333°) + 1
```

**27 Nakshatras:**
1. Ashwini (0° - 13°20')
2. Bharani (13°20' - 26°40')
3. Krittika (26°40' - 40°)
... (and so on, each spanning 13°20')

**Implementation:** `astronomy/nakshatra.go:Calculate()`

**Pada (Quarters):**
Each Nakshatra is divided into 4 padas of 3°20' each:
```
Pada_number = floor((Moon_longitude mod 13.333°) / 3.333°) + 1
```

**Ruling Planets:**
Each Nakshatra is ruled by one of the 9 planets in Vedic astrology:
- Cycle: Ketu, Venus, Sun, Moon, Mars, Rahu, Jupiter, Saturn, Mercury (repeats 3 times)

**Reference:**
- Raman, B. V. (1991). *Studies in Jaimini Astrology*
- Iyer, V. S. (1982). *Astrology and Nakshatras*

### 3. Yoga

Yoga represents the combined motion of Sun and Moon.

**Formula:**
```
Yoga_angle = (Sun_longitude + Moon_longitude) mod 360°
Yoga_number = floor(Yoga_angle / 13.333°) + 1
```

**27 Yogas:**
1. Vishkumbha
2. Priti
3. Ayushman
... (27 total, similar division to Nakshatras)

**Implementation:** `astronomy/yoga.go:Calculate()`

**Reference:**
- Sharma, R. S. (1996). *Ancient Indian Astronomy*

### 4. Karana

A Karana is half of a Tithi.

**Formula:**
```
Karana_number = floor(Tithi_angle / 6°)
```

**11 Karanas (repeated pattern):**
- 4 Fixed Karanas: Shakuni, Chatushpada, Naga, Kimstughna
- 7 Movable Karanas: Bava, Balava, Kaulava, Taitila, Gara, Vanija, Vishti
  (repeated 8 times in each lunar month)

**Implementation:** `astronomy/karana.go:Calculate()`

**Pattern:**
1. First half of Shukla Paksha Pratipada: Kimstughna (fixed)
2. Rest of month: 7 movable Karanas repeated 8 times (56 half-tithis)
3. Last half of Krishna Paksha Chaturdashi: Shakuni (fixed)
4. Amavasya: Chatushpada (first half), Naga (second half)

**Reference:**
- Dasa, N. R. (2005). *Principles of Panchangam*

### 5. Vara (Weekday)

Vara is simply the day of the week, but in Vedic astronomy, it has specific planetary rulerships.

**Planetary Hours:**
Each day is divided into 24 planetary hours, with the ruler of the first hour giving the day its name.

**Sequence:** Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn (repeating)

**Implementation:** `astronomy/vara.go:Calculate()`

## Interpolation Methods

### Linear Interpolation

For two data points (x₀, y₀) and (x₁, y₁), the value at x is:

```
y = y₀ + (x - x₀) * (y₁ - y₀) / (x₁ - x₀)
```

**Angle Wrapping for Longitude:**
```
if |Δλ| > 180°:
    if λ₀ > λ₁: λ₁ += 360°
    else: λ₀ += 360°
```

**Implementation:** `astronomy/ephemeris/interpolation.go:linearInterpolation()`

**Accuracy:** ±0.01° for intervals < 1 day

### Lagrange Polynomial Interpolation

For n data points, the Lagrange interpolating polynomial is:

```
P(x) = Σᵢ₌₀ⁿ⁻¹ yᵢ * Lᵢ(x)

where Lᵢ(x) = Πⱼ₌₀,ⱼ≠ᵢⁿ⁻¹ (x - xⱼ) / (xᵢ - xⱼ)
```

**Implementation:** `astronomy/ephemeris/interpolation.go:lagrangeInterpolation()`

**Recommended Order:** 5-7 points for planetary positions

**Accuracy:** ±0.001° for intervals < 1 day with 5-point interpolation

**Reference:** Meeus, J. (1998). *Astronomical Algorithms*, Chapter 3

### Cubic Spline Interpolation

For n data points, cubic splines satisfy:

```
Sᵢ(x) = aᵢ + bᵢ(x - xᵢ) + cᵢ(x - xᵢ)² + dᵢ(x - xᵢ)³
```

With continuity conditions:
- S(xᵢ) = yᵢ (interpolation)
- S'(xᵢ⁺) = S'(xᵢ₊₁⁻) (first derivative continuity)
- S''(xᵢ⁺) = S''(xᵢ₊₁⁻) (second derivative continuity)

**Natural Spline Boundary Conditions:**
```
S''(x₀) = 0
S''(xₙ₋₁) = 0
```

**Implementation:** `astronomy/ephemeris/interpolation.go:cubicSplineInterpolation()`

**Accuracy:** ±0.0001° for intervals < 1 day

**Reference:**
- Press, W. H., et al. (2007). *Numerical Recipes*, 3rd ed.
- De Boor, C. (2001). *A Practical Guide to Splines*

## Retrograde Motion Detection

### Speed Calculation

Planetary speed is calculated as the rate of change of longitude:

```
v = dλ/dt

Numerically: v ≈ (λ(t + Δt) - λ(t - Δt)) / (2Δt)
```

**Implementation:** `astronomy/ephemeris/retrograde.go:DetectRetrogradeMotion()`

### Motion Classification

- **Direct Motion:** v > 0.01 °/day (eastward)
- **Retrograde Motion:** v < -0.01 °/day (westward)
- **Stationary:** |v| ≤ 0.01 °/day

### Station Detection

Stations (stationary points) occur when dλ/dt = 0.

**Bisection Method:**
1. Sample planetary positions at regular intervals
2. Detect sign change in velocity
3. Use bisection to refine station time to ±1.4 minutes

**Implementation:** `astronomy/ephemeris/retrograde.go:FindPlanetaryStation()`

**Mathematical Formulation:**
```
Find t where f(t) = dλ/dt = 0

Using bisection:
if f(t₁) * f(t₂) < 0:
    station exists in interval [t₁, t₂]
    t_mid = (t₁ + t₂) / 2
    recurse until |t₂ - t₁| < tolerance
```

**Reference:**
- Meeus, J. (1998). *Astronomical Algorithms*, Chapter 36
- Espenak, F. (2009). *Planetary Phenomena*

### Retrograde Period

A retrograde period spans from station retrograde to station direct.

**Typical Durations:**
- Mercury: ~20-24 days
- Venus: ~40-43 days
- Mars: ~60-80 days
- Jupiter: ~120-140 days
- Saturn: ~140-160 days

**Implementation:** `astronomy/ephemeris/retrograde.go:FindRetrogradePeriod()`

## Sunrise and Sunset Calculations

### Solar Coordinates

The Sun's position is calculated using ephemeris data, providing:
- Ecliptic longitude (λ☉)
- Right ascension (α☉)
- Declination (δ☉)

### Hour Angle

The hour angle (H) at sunrise/sunset when the Sun crosses the horizon:

```
cos(H) = -tan(φ) * tan(δ☉)

where φ = observer latitude
      δ☉ = solar declination
```

**Atmospheric Refraction Correction:** -0.833° (standard)
**Additional Corrections:**
- Solar semi-diameter: ~0.267°
- Parallax: ~0.0024°

### Sunrise Time

```
LST_sunrise = α☉ - H

where LST = Local Sidereal Time
      α☉ = Sun's right ascension
      H = hour angle
```

Convert to local time using observer's longitude and equation of time.

**Implementation:** `astronomy/sunrise.go:CalculateSunrise()`

**Reference:**
- Meeus, J. (1998). *Astronomical Algorithms*, Chapter 15
- NOAA Solar Calculator algorithms

## Validation and Accuracy

### Tolerance Standards

- **Tithi transitions:** ±3 minutes
- **Nakshatra transitions:** ±3 minutes
- **Sunrise/sunset times:** ±2 minutes
- **Planetary positions:** ±0.001° (inner planets), ±0.01° (outer planets)

### Validation Sources

1. **Drik Panchang:** Modern computational panchangam
2. **Swiss Ephemeris:** Reference ephemeris
3. **JPL Horizons:** NASA's high-precision ephemeris system
4. **Traditional Printed Panchangams:** Regional almanacs

**Implementation:** `astronomy/validation/validation_framework.go`

## References

### Primary References

1. Meeus, J. (1998). *Astronomical Algorithms*, 2nd ed. Willmann-Bell.
2. Urban, S. E., & Seidelmann, P. K. (Eds.). (2012). *Explanatory Supplement to the Astronomical Almanac*, 3rd ed. University Science Books.
3. Montenbruck, O., & Pfleger, T. (2000). *Astronomy on the Personal Computer*, 4th ed. Springer.

### Vedic Astronomy References

4. Burgess, E. (1858). *Translation of the Surya Siddhanta*. Journal of the American Oriental Society, Vol. 6.
5. Sewell, R., & Dikshit, S. B. (1896). *The Indian Calendar*. Swan Sonnenschein & Co.
6. Raman, B. V. (1991). *Studies in Jaimini Astrology*. IBH Prakashana.

### Ephemeris References

7. Folkner, W. M., et al. (2014). *The Planetary and Lunar Ephemerides DE430 and DE431*. JPL IOM 392R-14-003.
8. Swiss Ephemeris Documentation. Astrodienst. https://www.astro.com/swisseph/

### Computational References

9. Press, W. H., et al. (2007). *Numerical Recipes: The Art of Scientific Computing*, 3rd ed. Cambridge University Press.
10. De Boor, C. (2001). *A Practical Guide to Splines*. Springer.

### Online Resources

11. JPL Horizons System: https://ssd.jpl.nasa.gov/horizons.cgi
12. NOAA Solar Calculator: https://www.esrl.noaa.gov/gmd/grad/solcalc/
13. Drik Panchang: https://www.drikpanchang.com/

## Implementation Notes

All algorithms are implemented in Go with:
- Full unit test coverage (>90%)
- OpenTelemetry observability
- Error handling and validation
- Caching for performance optimization

For implementation details, see the source code in:
- `astronomy/` - Core astronomical calculations
- `astronomy/ephemeris/` - Ephemeris integration and interpolation
- `astronomy/validation/` - Validation framework

---

*Last Updated: 2025-11-18*
*Maintainer: Panchangam Development Team*
