# Domain Context: Panchangam Astronomy

This document provides essential domain knowledge about Panchangam, Hindu astronomy, and the astronomical concepts used in this project. Understanding these concepts is crucial for implementing features and making architectural decisions.

## What is Panchangam?

**Panchangam** (Sanskrit: पञ्चाङ्ग, literally "five limbs") is the Hindu astronomical calendar system that provides daily astronomical and astrological data. It is essential for determining auspicious times (muhurtas) for various activities and religious observances.

### The Five Elements (Pancha Angas)

1. **Tithi** - Lunar day
2. **Vara** - Weekday
3. **Nakshatra** - Lunar mansion/constellation
4. **Yoga** - Luni-solar combination
5. **Karana** - Half-tithi

## Core Astronomical Concepts

### 1. Tithi (Lunar Day)

**Definition**: A tithi is 1/30th of a lunar month, representing a 12° angular separation between the Sun and Moon.

**Mathematical Formula**:
```
Tithi Number = floor((Moon_Longitude - Sun_Longitude) / 12°) + 1
Progress = ((Moon_Longitude - Sun_Longitude) % 12°) / 12°
```

**Key Points**:
- There are 30 tithis in a lunar month (15 in Shukla Paksha/waxing, 15 in Krishna Paksha/waning)
- Duration varies: ~19h 59m to ~26h 47m due to elliptical orbits
- A tithi can span less than or more than one solar day
- Transitions can occur at any time, not at midnight

**Tithi Names (Shukla Paksha - Waxing Moon)**:
1. Pratipada (प्रतिपदा)
2. Dwitiya (द्वितीया)
3. Tritiya (तृतीया)
4. Chaturthi (चतुर्थी)
5. Panchami (पञ्चमी)
6. Shashthi (षष्ठी)
7. Saptami (सप्तमी)
8. Ashtami (अष्टमी)
9. Navami (नवमी)
10. Dashami (दशमी)
11. Ekadashi (एकादशी)
12. Dwadashi (द्वादशी)
13. Trayodashi (त्रयोदशी)
14. Chaturdashi (चतुर्दशी)
15. Purnima (पूर्णिमा) - Full Moon

**Krishna Paksha (Waning Moon)**: Same names (1-14) + Amavasya (अमावस्या) - New Moon (15)

**Implementation Consideration**:
```go
// Tithi calculation must handle angle wrapping
func CalculateTithi(sunLong, moonLong float64) (int, float64, error) {
    diff := moonLong - sunLong
    if diff < 0 {
        diff += 360.0  // Handle angle wrapping
    }

    tithiNumber := int(diff/12.0) + 1
    progress := math.Mod(diff, 12.0) / 12.0

    if tithiNumber > 30 {
        tithiNumber = 30
    }

    return tithiNumber, progress, nil
}
```

### 2. Nakshatra (Lunar Mansion)

**Definition**: The 27 (or 28) divisions of the zodiac, each spanning 13°20' (13.333°).

**Mathematical Formula**:
```
Nakshatra Number = floor(Moon_Longitude / 13.333°) + 1
```

**27 Nakshatras**:
1. Ashwini (अश्विनी) - 0° to 13°20'
2. Bharani (भरणी) - 13°20' to 26°40'
3. Krittika (कृत्तिका) - 26°40' to 40°
4. Rohini (रोहिणी) - 40° to 53°20'
5. Mrigashira (मृगशीर्ष) - 53°20' to 66°40'
6. Ardra (आर्द्रा) - 66°40' to 80°
7. Punarvasu (पुनर्वसु) - 80° to 93°20'
8. Pushya (पुष्य) - 93°20' to 106°40'
9. Ashlesha (आश्लेषा) - 106°40' to 120°
10. Magha (मघा) - 120° to 133°20'
11. Purva Phalguni (पूर्व फाल्गुनी) - 133°20' to 146°40'
12. Uttara Phalguni (उत्तर फाल्गुनी) - 146°40' to 160°
13. Hasta (हस्त) - 160° to 173°20'
14. Chitra (चित्रा) - 173°20' to 186°40'
15. Swati (स्वाति) - 186°40' to 200°
16. Vishakha (विशाखा) - 200° to 213°20'
17. Anuradha (अनुराधा) - 213°20' to 226°40'
18. Jyeshtha (ज्येष्ठा) - 226°40' to 240°
19. Mula (मूल) - 240° to 253°20'
20. Purva Ashadha (पूर्व आषाढा) - 253°20' to 266°40'
21. Uttara Ashadha (उत्तर आषाढा) - 266°40' to 280°
22. Shravana (श्रवण) - 280° to 293°20'
23. Dhanishta (धनिष्ठा) - 293°20' to 306°40'
24. Shatabhisha (शतभिषा) - 306°40' to 320°
25. Purva Bhadrapada (पूर्व भाद्रपदा) - 320° to 333°20'
26. Uttara Bhadrapada (उत्तर भाद्रपदा) - 333°20' to 346°40'
27. Revati (रेवती) - 346°40' to 360°

**Each Nakshatra has**:
- 4 Padas (quarters) of 3°20' each
- Ruling planet
- Deity
- Symbol
- Gender, nature, quality classifications

**Special Nakshatras**:
- **Rohini, Pushya, Hasta**: Generally auspicious (good for most activities)
- **Mula, Ashlesha, Ardra**: Considered inauspicious for certain activities

### 3. Yoga

**Definition**: The sum of solar and lunar longitudes divided into 27 equal parts.

**Mathematical Formula**:
```
Yoga = floor((Sun_Longitude + Moon_Longitude) / 13.333°) + 1
```

**27 Yogas**:
1. Vishkambha (विष्कम्भ)
2. Priti (प्रीति)
3. Ayushman (आयुष्मान)
4. Saubhagya (सौभाग्य)
5. Shobhana (शोभन)
6. Atiganda (अतिगण्ड)
7. Sukarma (सुकर्मा)
8. Dhriti (धृति)
9. Shula (शूल)
10. Ganda (गण्ड)
11. Vriddhi (वृद्धि)
12. Dhruva (ध्रुव)
13. Vyaghata (व्याघात)
14. Harshana (हर्षण)
15. Vajra (वज्र)
16. Siddhi (सिद्धि)
17. Vyatipata (व्यतीपात)
18. Variyan (वरीयान्)
19. Parigha (परिघ)
20. Shiva (शिव)
21. Siddha (सिद्ध)
22. Sadhya (साध्य)
23. Shubha (शुभ)
24. Shukla (शुक्ल)
25. Brahma (ब्रह्म)
26. Indra (इन्द्र)
27. Vaidhriti (वैधृति)

**Inauspicious Yogas**: Vyaghata (#13) and Vyatipata (#17) - avoid important activities

### 4. Karana

**Definition**: Half of a tithi (6° of angular separation between Sun and Moon).

**Mathematical Formula**:
```
Karana = floor((Moon_Longitude - Sun_Longitude) / 6°) % 11 + 1
```

**11 Karanas**:
- **7 Movable (Chara)** - repeat 8 times each in a lunar month:
  1. Bava (बव)
  2. Balava (बालव)
  3. Kaulava (कौलव)
  4. Taitila (तैतिल)
  5. Gara (गर)
  6. Vanija (वणिज)
  7. Vishti/Bhadra (विष्टि/भद्रा) - Most inauspicious

- **4 Fixed (Sthira)** - occur once per lunar month:
  8. Shakuni (शकुनि) - Last half of 14th Krishna Paksha tithi
  9. Chatushpada (चतुष्पद) - First half of Amavasya
  10. Naga (नाग) - Second half of Amavasya
  11. Kimstughna (किंस्तुघ्न) - First half of Pratipada (Shukla Paksha)

**Special Karana**:
- **Vishti (Bhadra)**: Extremely inauspicious, avoid important activities

### 5. Vara (Weekday)

**Solar Days** (sunrise to sunrise, not midnight to midnight):

1. **Ravivara** (रविवार) - Sunday - Sun
2. **Somavara** (सोमवार) - Monday - Moon
3. **Mangalavara** (मंगलवार) - Tuesday - Mars
4. **Budhavara** (बुधवार) - Wednesday - Mercury
5. **Guruvara** (गुरुवार) - Thursday - Jupiter
6. **Shukravara** (शुक्रवार) - Friday - Venus
7. **Shanivara** (शनिवार) - Saturday - Saturn

**Important**: In Panchangam, the day starts at **sunrise**, not midnight!

## Astronomical Foundation

### Coordinate Systems

**1. Tropical vs. Sidereal Zodiac**

- **Tropical Zodiac**: Fixed to equinoxes (Western astrology)
- **Sidereal Zodiac**: Fixed to stars (Vedic astronomy)
- **Difference**: Ayanamsa (precession correction)

**Ayanamsa** (अयनांश):
```
Sidereal_Longitude = Tropical_Longitude - Ayanamsa
```

**Common Ayanamsa Systems**:
- **Lahiri** (Chitrapaksha): Most commonly used in India (23.85° for year 2000)
- **Raman**: Alternative system
- **Krishnamurthy**: KP astrology system

**Implementation**:
```go
// Calculate ayanamsa for a given date
func CalculateLahiriAyanamsa(jd float64) float64 {
    // T = centuries from J2000.0
    t := (jd - 2451545.0) / 36525.0

    // Lahiri formula
    ayanamsa := 23.85 + 0.013888889*t

    return ayanamsa
}
```

### Swiss Ephemeris

**Purpose**: Provides highly accurate planetary positions

**Key Features**:
- Precision: 0.001 arcsecond
- Time span: 30,000 years (13000 BCE to 17000 CE)
- Based on NASA JPL data
- Compressed from 2.8 GB to 99 MB

**Usage in Project**:
```go
// Get planetary position
sunPos := ephemeris.GetPlanetPosition("Sun", julianDay)
moonPos := ephemeris.GetPlanetPosition("Moon", julianDay)
```

### Time Calculations

**Julian Day**: Continuous day count from noon UTC on January 1, 4713 BCE

```go
func DateToJulianDay(date time.Time) float64 {
    // Standard formula
    year := date.Year()
    month := int(date.Month())
    day := date.Day()

    a := (14 - month) / 12
    y := year + 4800 - a
    m := month + 12*a - 3

    jd := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045

    // Add time fraction
    hour := date.Hour()
    minute := date.Minute()
    second := date.Second()
    dayFraction := (float64(hour) - 12) / 24.0 + float64(minute) / 1440.0 + float64(second) / 86400.0

    return float64(jd) + dayFraction
}
```

**Sunrise/Sunset Calculation**:
- Account for atmospheric refraction (~0.5833°)
- Adjust for observer elevation
- Use iterative method for precision

## Regional Variations

### Month Systems

**1. Amanta System** (South India)
- Month ends on New Moon (Amavasya)
- Used in: Tamil Nadu, Kerala, Karnataka, Andhra Pradesh

**2. Purnimanta System** (North India)
- Month ends on Full Moon (Purnima)
- Used in: North India, Gujarat, Maharashtra

**Important**: Same day can have different month names in different systems!

### Calendar Systems

**1. Vikrama Samvat**
- Starts: 57 BCE
- New Year: Chaitra Shukla Pratipada (March-April)
- Used in: North India

**2. Shaka Samvat**
- Starts: 78 CE
- Official calendar of India
- New Year: Chaitra 1 (March 22 typically)

**3. Kali Yuga**
- Starts: 3102 BCE
- Theoretical/scriptural reference

### Calculation Methods

**1. Drik Ganita** (Observational)
- Based on actual planetary positions
- Uses Swiss Ephemeris or similar
- **More accurate** (recommended)
- Used by: Modern panchangams

**2. Vakya Ganita** (Traditional)
- Based on ancient texts (Surya Siddhanta)
- Uses memorized verses
- Less accurate (can differ by up to 12 hours)
- Used by: Traditional Tamil panchangams (Pambu Panchangam)

## Muhurta (Auspicious Timing)

### Special Muhurtas

**1. Abhijit Muhurta**
- Time: 24 minutes before to 24 minutes after local solar noon
- Always auspicious except Wednesdays
- Duration: 48 minutes
- Calculation: Based on sunrise and sunset times

```go
func CalculateAbhijitMuhurta(sunrise, sunset time.Time) (start, end time.Time) {
    dayDuration := sunset.Sub(sunrise)
    muhurtaDuration := dayDuration / 15  // 15 muhurtas in a day
    midday := sunrise.Add(dayDuration / 2)

    start = midday.Add(-muhurtaDuration / 2)
    end = midday.Add(muhurtaDuration / 2)

    return start, end
}
```

**2. Brahma Muhurta**
- Time: 96 minutes (2 muhurtas) before sunrise
- Highly auspicious for spiritual practices
- Duration: 96 minutes (48 minutes each for two muhurtas)

**3. Rahu Kala** (Inauspicious Period)
- Daily 90-minute period ruled by Rahu (shadow planet)
- Varies by weekday
- Avoid starting new activities

| Day | Rahu Kala Position |
|-----|-------------------|
| Sunday | 4th period (afternoon) |
| Monday | 2nd period (morning) |
| Tuesday | 7th period (evening) |
| Wednesday | 5th period (afternoon) |
| Thursday | 6th period (late afternoon) |
| Friday | 3rd period (late morning) |
| Saturday | 1st period (early morning) |

### Muhurta Selection Criteria

**Auspicious Factors**:
- Favorable tithi (avoid Amavasya, Chaturdashi for most activities)
- Auspicious nakshatra (Rohini, Pushya, Hasta)
- Good yoga (avoid Vyaghata, Vyatipata)
- Favorable karana (avoid Vishti/Bhadra)
- Appropriate weekday for activity
- Good planetary transits

**Activity-Specific**:
- **Marriage**: Complex calculations, requires birth charts
- **Griha Pravesh**: House warming - Rohini, Uttara nakshatra preferred
- **Business**: Strong Mercury, Jupiter
- **Education**: Saraswati worship days, Pushya nakshatra
- **Travel**: Avoid Vishti karana, inauspicious yogas

## Festivals and Observances

### Major Festivals (Tithi-based)

**Fixed to Tithi**:
- **Diwali**: Krishna Paksha Amavasya (Kartik month)
- **Holi**: Phalguna Purnima
- **Janmashtami**: Krishna Paksha Ashtami (Shravana/Bhadrapada)
- **Ram Navami**: Shukla Paksha Navami (Chaitra)
- **Ekadashi**: 11th tithi (both pakshas) - fasting day

**Solar-based**:
- **Makar Sankranti**: Sun enters Capricorn (~Jan 14-15)
- **Pongal**: Tamil harvest festival (3-4 days from Makar Sankranti)
- **Vishu**: Malayalam New Year (Sun enters Aries)

### Festival Calculation Complexity

**Challenge**: Festivals can span multiple civil days

Example: If Ashtami tithi spans 2 days, which day to celebrate Janmashtami?

**Rules**:
- Consider which day has midnight during the tithi
- Check regional customs
- Some festivals observe on previous/next day based on specific rules

## Data Accuracy Requirements

### Precision Needs

**Planetary Positions**:
- Minimum: 0.01° (36 arcseconds)
- Recommended: 0.001° (3.6 arcseconds)
- Swiss Ephemeris provides: 0.001 arcsecond

**Time Precision**:
- Tithi transitions: ±1 minute acceptable
- Sunrise/Sunset: ±1 minute acceptable
- Muhurta boundaries: ±2 minutes acceptable

### Validation Sources

**Historical Validation**:
- Cross-reference with traditional printed panchangams
- Drik Panchang (online reference)
- Indian Astronomical Ephemeris
- Festival dates from government calendars

**Test Cases**:
- Known eclipse dates
- Historical festival dates
- Solstices and equinoxes
- Known astronomical events

## Common Pitfalls for Developers

### 1. Time Zones
```go
// WRONG: Using UTC for local calculations
sunrise := CalculateSunrise(location, time.Now().UTC())

// CORRECT: Use local time zone
localTime := time.Now().In(location.TimeZone)
sunrise := CalculateSunrise(location, localTime)
```

### 2. Angle Wrapping
```go
// Handle angles > 360° or < 0°
func NormalizeAngle(angle float64) float64 {
    angle = math.Mod(angle, 360.0)
    if angle < 0 {
        angle += 360.0
    }
    return angle
}
```

### 3. Sunrise as Day Start
```go
// WRONG: Day starts at midnight
panchangam := GetPanchangam(date.Truncate(24 * time.Hour))

// CORRECT: Day starts at sunrise
sunriseTime := CalculateSunrise(location, date)
panchangam := GetPanchangam(sunriseTime)
```

### 4. Ayanamsa Application
```go
// Always apply ayanamsa to convert tropical to sidereal
siderealLongitude := tropicalLongitude - ayanamsa
```

### 5. Leap Months (Adhik Masa)
- Occurs when no solar transition happens in a lunar month
- Must be handled in festival calculations
- Affects month numbering

## Glossary

- **Adhik Masa**: Leap month in Hindu calendar
- **Amavasya**: New Moon (30th tithi)
- **Ayanamsa**: Precession correction angle
- **Drik**: Observational/actual calculation method
- **Julian Day**: Continuous day count used in astronomy
- **Karana**: Half-tithi (6° angular separation)
- **Nakshatra**: Lunar mansion (27 divisions of zodiac)
- **Paksha**: Lunar fortnight (Shukla=waxing, Krishna=waning)
- **Purnima**: Full Moon (15th tithi)
- **Tithi**: Lunar day (12° angular separation)
- **Vakya**: Traditional verse-based calculation method
- **Yoga**: Luni-solar combination
- **Sidereal**: Fixed to stars
- **Tropical**: Fixed to seasons (equinoxes)

## References

- **Surya Siddhanta**: Ancient astronomical text (4th-5th century CE)
- **Siddhanta Shiromani**: Bhaskaracharya's work (12th century CE)
- **Swiss Ephemeris**: Modern high-precision ephemeris
- **Indian Astronomical Ephemeris**: Annual publication by Positional Astronomy Centre
- **Drik Panchang**: Online reference implementation

## For LLM Agents

When implementing Panchangam features:

1. **Always use sidereal (Nirayana) coordinates** with appropriate ayanamsa
2. **Remember solar days start at sunrise**, not midnight
3. **Account for regional variations** (Amanta vs Purnimanta)
4. **Validate against known values** (historical panchangams)
5. **Handle edge cases** (polar regions, angle wrapping, leap months)
6. **Precision matters** - use appropriate floating-point precision
7. **Time zones are critical** - always be explicit
8. **Test extensively** - astronomical calculations are complex

The domain has **deep cultural and religious significance** - accuracy and respect for tradition are paramount.
