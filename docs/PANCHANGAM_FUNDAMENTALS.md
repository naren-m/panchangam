# Panchangam: A Mathematical Guide to the Hindu Astronomical Almanac

## Introduction

This document explains the Panchangam (Sanskrit: पञ्चाङ्ग, "five limbs") from a mathematical and engineering perspective. If you understand vectors, coordinate systems, and modular arithmetic, you already have the tools to understand ancient Hindu astronomy.

The Panchangam is essentially a **time-series function** that maps any instant to five discrete astronomical states. Think of it as a lookup table where the key is (date, time, location) and the values are derived from celestial mechanics.

---

## 1. The Coordinate System: Earth as Origin

### 1.1 Geocentric Model

Panchangam calculations use a **geocentric coordinate system**—Earth at the origin, with celestial bodies projected onto a conceptual sphere called the **celestial sphere**.

```
                    Celestial North Pole
                           │
                           │
                    ╭──────┴──────╮
                 ╱                  ╲
               ╱    Stars "fixed"     ╲
              │     at ∞ distance      │
              │                        │
        ──────┼────── EARTH ──────────┼────── → Vernal Equinox (0°)
              │         (0,0)          │           (Aries Point)
              │                        │
               ╲                      ╱
                 ╲                  ╱
                    ╰──────────────╯
                           │
                    Celestial South Pole
```

**Key insight**: While we know the Sun doesn't orbit Earth, the geocentric model is mathematically equivalent for computing **apparent positions**—what we observe from Earth. This is the same principle used in navigation and satellite tracking.

### 1.2 The Ecliptic: The Reference Plane

The **ecliptic** is the apparent path of the Sun through the sky over a year. It's tilted 23.4° from the celestial equator (Earth's equator projected onto the celestial sphere).

```
Side View of Celestial Sphere:

         Celestial
          Equator ─────────────────────  (0° declination)
              ╲      ╱
               ╲    ╱  23.4°
                ╲  ╱
    Ecliptic ────╳────────────────────  (Sun's path)
                  ╲
                   ╲
                    Earth
```

**Why the ecliptic matters**: All Panchangam calculations use **ecliptic longitude** (λ), measured in degrees from 0° to 360° along this plane.

### 1.3 Longitude Measurement

Ecliptic longitude is measured counter-clockwise from the **Vernal Equinox point** (where the ecliptic crosses the celestial equator, around March 20-21).

```
          90° (Summer Solstice)
               │
               │
    180° ──────┼────── 0° (Vernal Equinox)
               │
               │
          270° (Winter Solstice)

Moon at 45°:  λₘ = 45°  →  In "Taurus" region
Sun at 280°:  λₛ = 280° →  In "Sagittarius" region
```

**Normalization function**: All angles are normalized to [0°, 360°):
```
normalize(θ) = ((θ mod 360) + 360) mod 360
```

---

## 2. The Five Elements (Pancha Anga)

The Panchangam derives five elements from the positions of the Sun (λₛ) and Moon (λₘ):

| Element | Symbol | Input | Formula | Divisions |
|---------|--------|-------|---------|-----------|
| **Tithi** | तिथि | λₘ - λₛ | Angular difference | 30 |
| **Nakshatra** | नक्षत्र | λₘ | Moon's absolute position | 27 |
| **Yoga** | योग | λₘ + λₛ | Angular sum | 27 |
| **Karana** | करण | λₘ - λₛ | Half-tithi | 60 |
| **Vara** | वार | Date | Day of week | 7 |

Notice the elegant pattern:
- **Tithi and Karana**: Based on **difference** (relative position)
- **Nakshatra**: Based on **Moon's position** (absolute)
- **Yoga**: Based on **sum** (combined influence)
- **Vara**: Simple modular arithmetic on date

---

## 3. Tithi: The Lunar Day

### 3.1 Mathematical Definition

**Tithi** measures the angular separation between the Moon and Sun, quantized into 30 discrete values.

```
Given:
  λₛ = Sun's ecliptic longitude [0°, 360°)
  λₘ = Moon's ecliptic longitude [0°, 360°)

Angular Difference:
  Δ = normalize(λₘ - λₛ)

Tithi Number:
  T = floor(Δ / 12°) + 1    where T ∈ {1, 2, ..., 30}
```

**Why 12°?** The Moon completes one orbit (360°) relative to the Sun in approximately 29.5 days (synodic month). Dividing by 30 tithis:
```
360° ÷ 30 = 12° per tithi
```

### 3.2 Visual Representation

```
New Moon (T=1)                    Full Moon (T=15)
     Δ = 0°                           Δ = 180°
        │                                │
        ▼                                ▼
    ☀️───🌑                          ☀️───────────🌕
    Sun  Moon                        Sun          Moon
    (conjunction)                    (opposition)

                The Tithi Cycle

    T=1 ─────── T=8 ─────── T=15 ─────── T=22 ─────── T=30 ─── T=1
     │           │           │            │            │        │
     ▼           ▼           ▼            ▼            ▼        ▼
    🌑          🌓          🌕           🌗           🌑       🌑
   New      Quarter       Full        Quarter       New      New
   Moon      Moon         Moon         Moon        Moon     Moon

    └──── Shukla Paksha ────┴──── Krishna Paksha ────┘
          (Waxing)                 (Waning)
```

### 3.3 Paksha (Lunar Phase)

```
Shukla Paksha (Bright Half):  T ∈ {1, 2, ..., 15}   when Δ ∈ [0°, 180°)
Krishna Paksha (Dark Half):   T ∈ {16, 17, ..., 30} when Δ ∈ [180°, 360°)
```

### 3.4 Variable Duration

Tithis don't have constant duration because:
1. **Elliptical orbits**: Moon and Earth don't move at constant angular velocity
2. **Orbital perturbations**: Gravitational effects from Sun and planets

```
Average Tithi Duration:
  T_avg = Synodic Month ÷ 30 = 29.53 days ÷ 30 ≈ 23.6 hours

Actual Range:
  T_min ≈ 19h 59m  (when Moon at perigee, moving fastest)
  T_max ≈ 26h 47m  (when Moon at apogee, moving slowest)
```

This is analogous to how Earth's orbital speed varies (faster at perihelion, slower at aphelion), following Kepler's Second Law.

---

## 4. Nakshatra: The Lunar Mansion

### 4.1 Mathematical Definition

**Nakshatra** divides the ecliptic into 27 equal segments based on the Moon's absolute position.

```
Nakshatra Number:
  N = floor(λₘ / 13.333°) + 1    where N ∈ {1, 2, ..., 27}

Or equivalently:
  N = floor(λₘ × 27 / 360) + 1
```

**Why 27?** The Moon completes one sidereal orbit (360° relative to stars) in approximately 27.3 days. Ancient astronomers assigned one "mansion" per day.

### 4.2 The 27 Nakshatras

```
Each Nakshatra spans: 360° ÷ 27 = 13°20' = 13.333°

     0°        13.33°      26.67°      40°
     │           │           │          │
     ▼           ▼           ▼          ▼
┌─────────┬─────────┬─────────┬─────────┬───
│ Ashwini │ Bharani │Krittika │ Rohini  │ ...
│  1      │   2     │   3     │   4     │
└─────────┴─────────┴─────────┴─────────┴───
     │
     Moon at λₘ = 45° → Nakshatra = floor(45/13.33)+1 = 4 (Rohini)
```

### 4.3 Pada Subdivision

Each Nakshatra is further divided into 4 equal **padas** (quarters):

```
Pada span = 13.333° ÷ 4 = 3.333° = 3°20'

Total padas = 27 × 4 = 108 (a sacred number)

For Moon at λₘ:
  Pada = floor((λₘ mod 13.333°) / 3.333°) + 1    where Pada ∈ {1, 2, 3, 4}
```

---

## 5. Yoga: The Luni-Solar Combination

### 5.1 Mathematical Definition

**Yoga** uses the **sum** of Sun and Moon longitudes, creating a unique rhythmic pattern.

```
Combined Longitude:
  Σ = normalize(λₘ + λₛ)

Yoga Number:
  Y = floor(Σ / 13.333°) + 1    where Y ∈ {1, 2, ..., 27}
```

### 5.2 Physical Interpretation

While Tithi measures **relative** motion (Moon gaining on Sun), Yoga measures the **combined** celestial influence. Think of it as a phase relationship in a two-oscillator system.

```
Example Phase Diagram:

    λₛ = 100°, λₘ = 200°

    Tithi:     Δ = 200° - 100° = 100° → T = floor(100/12) + 1 = 9
    Yoga:      Σ = 200° + 100° = 300° → Y = floor(300/13.33) + 1 = 23
    Nakshatra: N = floor(200/13.33) + 1 = 16 (Vishakha)
```

### 5.3 Yoga Progression Rate

The combined longitude advances at:
```
dΣ/dt = dλₘ/dt + dλₛ/dt

Where:
  dλₘ/dt ≈ 13.2°/day  (Moon's average motion)
  dλₛ/dt ≈ 1.0°/day   (Sun's average motion)

Therefore:
  dΣ/dt ≈ 14.2°/day

Yoga duration ≈ 13.333° ÷ 14.2°/day ≈ 22.5 hours (average)
```

---

## 6. Karana: The Half-Tithi

### 6.1 Mathematical Definition

**Karana** is simply half a tithi, giving 60 divisions per lunar month.

```
Karana Number:
  K = floor(Δ / 6°) + 1    where K ∈ {1, 2, ..., 60}

Or:
  K = 2×(T-1) + (Δ mod 12° ≥ 6° ? 2 : 1)
```

### 6.2 Karana Types

```
There are 11 named Karanas, arranged in a specific pattern:

Fixed Karanas (occur once per month):
  1. Kimstughna  (K=1)
  2. Shakuni     (K=58)
  3. Chatushpada (K=59)
  4. Nagava      (K=60)

Rotating Karanas (cycle 7 times through the month):
  Bava → Balava → Kaulava → Taitila → Gara → Vanija → Vishti

Pattern: K=2 to K=57 cycles through the 7 rotating karanas:
  K = 2:  Bava
  K = 3:  Balava
  K = 4:  Kaulava
  K = 5:  Taitila
  K = 6:  Gara
  K = 7:  Vanija
  K = 8:  Vishti
  K = 9:  Bava (cycle repeats)
  ...
```

---

## 7. Vara: The Weekday

### 7.1 Simple Modular Arithmetic

```
Vara = (Julian Day Number + 1) mod 7

Mapping:
  0 → Ravivara    (Sunday)    ☀️
  1 → Somavara    (Monday)    ☽
  2 → Mangalavara (Tuesday)   ♂
  3 → Budhavara   (Wednesday) ☿
  4 → Guruvara    (Thursday)  ♃
  5 → Shukravara  (Friday)    ♀
  6 → Shanivara   (Saturday)  ♄
```

The names derive from the seven classical planets visible to the naked eye.

---

## 8. Rashi (Zodiac Signs): The Solar Context

While not part of the "five limbs," Rashi provides context for the Sun's position.

### 8.1 Mathematical Definition

```
Rashi divisions: 360° ÷ 12 = 30° per sign

Rashi Number:
  R = floor(λₛ / 30°) + 1    where R ∈ {1, 2, ..., 12}

Mapping:
  R=1:  Mesha (Aries)        0° - 30°
  R=2:  Vrishabha (Taurus)   30° - 60°
  R=3:  Mithuna (Gemini)     60° - 90°
  ...
  R=12: Meena (Pisces)       330° - 360°
```

### 8.2 Ayanamsa: The Sidereal Correction

**Important**: Hindu astronomy uses the **sidereal** zodiac (fixed to stars), while Western astronomy uses the **tropical** zodiac (fixed to equinoxes).

Due to Earth's axial precession (26,000-year cycle), these differ by the **Ayanamsa** (currently ~24°):

```
λ_sidereal = λ_tropical - Ayanamsa

Example:
  Tropical Sun longitude: λₛ = 45°
  Ayanamsa (Lahiri): 24.2°
  Sidereal longitude: 45° - 24.2° = 20.8° (still in Mesha/Aries)
```

---

## 9. Putting It All Together: The Calculation Pipeline

### 9.1 System Architecture

```
┌─────────────┐    ┌──────────────────┐    ┌────────────────────┐
│   INPUT     │    │   EPHEMERIS      │    │   PANCHANGAM       │
│             │    │   ENGINE         │    │   CALCULATOR       │
│ • Date      │───▶│                  │───▶│                    │
│ • Time      │    │ • Swiss Ephemeris│    │ • Tithi Formula    │
│ • Location  │    │ • JPL Data       │    │ • Nakshatra Formula│
│             │    │ • Ayanamsa       │    │ • Yoga Formula     │
└─────────────┘    └──────────────────┘    │ • Karana Formula   │
                           │                │ • Vara Calculation │
                           ▼                └────────────────────┘
                   ┌──────────────────┐              │
                   │  λₛ (Sun)        │              │
                   │  λₘ (Moon)       │              ▼
                   │  Julian Day      │     ┌────────────────────┐
                   └──────────────────┘     │     OUTPUT         │
                                            │                    │
                                            │ Tithi: Shukla 9    │
                                            │ Nakshatra: Rohini  │
                                            │ Yoga: Siddha       │
                                            │ Karana: Vishti     │
                                            │ Vara: Mangala      │
                                            └────────────────────┘
```

### 9.2 Complete Calculation Example

```
Given:
  Date: January 15, 2025, 12:00 UTC
  Location: Mumbai (19.076°N, 72.877°E)

Step 1: Convert to Julian Day
  JD = 2460691.0

Step 2: Get Planetary Positions (from ephemeris)
  λₛ (tropical) = 294.85°
  λₘ (tropical) = 87.42°
  Ayanamsa = 24.18°

Step 3: Convert to Sidereal
  λₛ (sidereal) = 294.85° - 24.18° = 270.67°
  λₘ (sidereal) = 87.42° - 24.18° = 63.24°

Step 4: Calculate Panchangam Elements

  Tithi:
    Δ = normalize(63.24° - 270.67°) = normalize(-207.43°) = 152.57°
    T = floor(152.57° / 12°) + 1 = 12 + 1 = 13 (Shukla Trayodashi)

  Nakshatra:
    N = floor(63.24° / 13.333°) + 1 = 4 + 1 = 5 (Mrigashira)

  Yoga:
    Σ = normalize(63.24° + 270.67°) = 333.91°
    Y = floor(333.91° / 13.333°) + 1 = 25 + 1 = 26 (Uthara Bhadrapada... wait)
    Actually: Y = floor(333.91° / 13.333°) + 1 = 26 (Uttarabhadra)

  Karana:
    K = floor(152.57° / 6°) + 1 = 25 + 1 = 26
    Rotating index = (26-2) mod 7 = 24 mod 7 = 3 → Kaulava

  Vara:
    (2460691 + 1) mod 7 = 2460692 mod 7 = 3 → Budhavara (Wednesday)

Result:
  Tithi: Shukla Trayodashi (13)
  Nakshatra: Mrigashira (5)
  Yoga: Uttarabhadra (26)
  Karana: Kaulava
  Vara: Wednesday
```

---

## 10. Engineering Considerations

### 10.1 Precision Requirements

| Component | Precision Needed | Rationale |
|-----------|-----------------|-----------|
| Ephemeris | ±0.001° | Tithi changes every 12°, need sub-degree accuracy |
| Time | ±1 minute | Tithi duration ~24h, transitions matter |
| Location | ±0.1° | Sunrise/sunset calculations |

### 10.2 Performance Optimization

```
Caching Strategy:
┌────────────────────────────────────────────────────────────┐
│  Date → Ephemeris Cache (expensive computation)            │
│                                                            │
│  Key: (JulianDay, Ayanamsa)                               │
│  Value: (λₛ, λₘ, sunrise, sunset)                          │
│  TTL: 1 hour (positions change slowly)                     │
└────────────────────────────────────────────────────────────┘

Batch Processing:
  For monthly calendar, fetch ephemeris data in one call
  Then compute derived values locally (O(1) per day)
```

### 10.3 Edge Cases

1. **Tithi Kshaya** (skipped tithi): A tithi that starts and ends within the same sunrise-to-sunrise day
2. **Adhika Tithi** (extra tithi): Two sunrises during the same tithi
3. **Longitude wrap-around**: Handle 359° → 0° transitions
4. **Timezone boundaries**: Ensure consistent date handling

---

## 11. Visualization: The Celestial Chart

The Panchangam UI renders these calculations as a geocentric chart:

```
                           0° (Aries/Mesha)
                                │
                        ┌───────┴───────┐
                     ╱ ╲  N1  │  N27 ╱  ╲
                   ╱     ╲────┼────╱      ╲
                 ╱    N2   ╲  │  ╱   N26    ╲
               ╱─────────────╲│╱─────────────╲
              │     NAKSHATRA RING (27)       │
              │   ┌───────────────────────┐   │
              │   │    RASHI RING (12)    │   │
    270° ─────│   │   ┌───────────────┐   │   │───── 90°
              │   │   │     ☀️ Sun     │   │   │
              │   │   │       │       │   │   │
              │   │   │   🌍 EARTH    │   │   │
              │   │   │       │       │   │   │
              │   │   │     🌙 Moon   │   │   │
              │   │   └───────────────┘   │   │
              │   └───────────────────────┘   │
              │                               │
               ╲                             ╱
                 ╲                         ╱
                   ╲                     ╱
                     ╲                 ╱
                        └───────┬───────┘
                                │
                           180° (Libra)

The arc from Sun to Moon = Tithi angle (Δ)
Moon's position on Nakshatra ring = Current Nakshatra
Sun's position on Rashi ring = Current Rashi (month)
```

---

## 12. Summary: The Elegant Mathematics

The Panchangam demonstrates elegant mathematical principles:

1. **Modular Arithmetic**: All calculations reduce to divisions and remainders on a circular domain [0°, 360°)

2. **Linear Transformations**: Each element is a linear function of celestial coordinates:
   - Tithi: T(λₘ, λₛ) = floor((λₘ - λₛ) / 12) + 1
   - Nakshatra: N(λₘ) = floor(λₘ / 13.333) + 1
   - Yoga: Y(λₘ, λₛ) = floor((λₘ + λₛ) / 13.333) + 1

3. **Hierarchical Decomposition**: Time is structured in nested cycles:
   - Tithi (30/month) → Paksha (2/month) → Masa (12/year) → Samvatsara (60-year cycle)

4. **Coordinate System Invariance**: The geocentric model, while not physically accurate, is mathematically equivalent for positional astronomy.

---

## References

### Mathematical Sources
- Meeus, Jean. "Astronomical Algorithms" (1991)
- Seidelmann, P. Kenneth. "Explanatory Supplement to the Astronomical Almanac" (1992)

### Traditional Sources
- Surya Siddhanta (ancient Sanskrit text on astronomy)
- Siddhanta Shiromani by Bhaskaracharya (12th century)

### Modern Implementation
- Swiss Ephemeris: https://www.astro.com/swisseph/
- JPL Horizons: https://ssd.jpl.nasa.gov/horizons/

---

*Document Version: 1.0.0*
*Last Updated: January 2026*
*Target Audience: Engineers, mathematicians, and developers working with Panchangam calculations*
