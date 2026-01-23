# Regional Variations in Panchangam Calculations

This document describes regional variations in Panchangam calculations across different parts of India and how they are implemented in this system.

## Overview

Indian astronomical calculations vary by region due to:
1. **Calendar Systems** (Amanta vs Purnimanta)
2. **Calculation Methods** (Drik Ganita vs Vakya Ganita)
3. **Festival Dates and Observances**
4. **Local Traditions and Customs**

## Calendar Systems

### Amanta System (New Moon Ending)

**Used in:** South India (Tamil Nadu, Karnataka, Andhra Pradesh, Kerala)

**Characteristics:**
- Lunar month ends on Amavasya (New Moon)
- Month naming based on the solar month in which Amavasya occurs
- Krishna Paksha comes before Shukla Paksha within the month

**Example:**
```
Chaitra month:
- Chaitra Krishna Paksha (waning) - days 1-15
- Chaitra Shukla Paksha (waxing) - days 1-15
- Ends with Amavasya (New Moon)
```

**Implementation:** `astronomy/traditional.go:AmantaCalendar`

### Purnimanta System (Full Moon Ending)

**Used in:** North India (Maharashtra, Gujarat, Madhya Pradesh, Rajasthan, UP)

**Characteristics:**
- Lunar month ends on Purnima (Full Moon)
- Month naming based on the solar month in which Purnima occurs
- Shukla Paksha comes before Krishna Paksha

**Example:**
```
Chaitra month:
- Chaitra Shukla Paksha (waxing) - days 1-15
- Chaitra Krishna Paksha (waning) - days 1-15
- Ends with Purnima (Full Moon)
```

**Implementation:** `astronomy/traditional.go:PurnimantaCalendar`

### Conversion Between Systems

```go
// Purnimanta to Amanta
if isKrishnaPaksha(tithi) {
    amantaMonth = purnimantaMonth + 1
} else {
    amantaMonth = purnimantaMonth
}

// Amanta to Purnimanta
if isKrishnaPaksha(tithi) {
    purnimantaMonth = amantaMonth - 1
} else {
    purnimantaMonth = amantaMonth
}
```

## Calculation Methods

### Drik Ganita (Observational Method)

**Characteristics:**
- Based on actual astronomical observations
- Uses precise ephemeris calculations
- More accurate, matches modern astronomical data
- Preferred in modern Panchangams

**Advantages:**
- Accuracy: ±2-3 minutes for sunrise
- Accounts for precession of equinoxes
- Uses actual planetary positions

**Regions:** Widely adopted across India in modern times

**Implementation:** Default calculation method in `astronomy/` package

### Vakya Ganita (Traditional Method)

**Characteristics:**
- Based on traditional astronomical tables
- Uses mean positions rather than true positions
- Preserves historical continuity
- Still used in traditional temples

**Differences from Drik Ganita:**
- Can differ by 1-2 days for festival dates
- Does not account for precession
- Uses fixed parameters from Surya Siddhanta

**Regions:** Traditional temples and orthodox communities

**Implementation:** `astronomy/traditional.go:VakyaCalculator`

### Comparison Table

| Aspect | Drik Ganita | Vakya Ganita |
|--------|-------------|--------------|
| Basis | Modern Ephemeris | Traditional Tables |
| Accuracy | ±2-3 minutes | ±10-30 minutes |
| Precession | Accounted | Not accounted |
| Planetary Positions | True | Mean |
| Usage | Modern Panchangams | Traditional Temples |

## Regional Festival Variations

### Tamil Nadu

**Unique Aspects:**
- Solar calendar (Vaisakhi system) for festivals
- Tamil month names: Chithirai, Vaikasi, Aani, etc.
- Pongal based on solar transit
- Specific muhurta rules for weddings

**Key Festivals:**
- Thai Pongal (solar)
- Tamil New Year (Chithirai 1)
- Aadi Perukku

**Implementation:** `api/examples/tamil_nd_plugin.go`

### Kerala

**Unique Aspects:**
- Malayalam calendar (Kollavarsham)
- Era: Kollam Era (CE year + 824/825)
- Solar months with unique names
- Specific ayanamsa correction

**Key Festivals:**
- Onam (Thiruvonam in Chingam)
- Vishu (Mesha Sankranti)

**Implementation:** `api/examples/kerala_plugin.go`

### Maharashtra and Gujarat

**Unique Aspects:**
- Purnimanta calendar system
- Emphasis on tithis for festivals
- Specific rules for adhika masa
- Distinct festival dates

**Key Festivals:**
- Gudi Padwa (Chaitra Shukla Pratipada)
- Ganesh Chaturthi
- Makar Sankranti

### Bengal and Odisha

**Unique Aspects:**
- Bengali calendar (Bangla Saal)
- Solar-lunar hybrid system
- New Year on Poila Baisakh
- Unique nakshatra-based festivals

**Key Festivals:**
- Pohela Boishakh
- Durga Puja (specific tithi rules)
- Rath Yatra

## Regional Sunrise and Sunset

Sunrise and sunset times vary by location due to:

### Latitude Effects

```
Sunrise time variation = f(latitude, declination)

For India (8°N to 35°N):
- Summer solstice: Earlier sunrise in north
- Winter solstice: Earlier sunrise in south
```

### Longitude Effects

```
Time difference = (Longitude₁ - Longitude₂) / 15° × 60 minutes

Example:
Ahmedabad (72.5°E) vs Kolkata (88.4°E)
Difference = (88.4 - 72.5) / 15 × 60 ≈ 64 minutes
```

### Regional Locations

| City | Latitude | Longitude | Typical Sunrise (Jan) | Typical Sunrise (Jul) |
|------|----------|-----------|----------------------|----------------------|
| Delhi | 28.7°N | 77.1°E | 07:15 | 05:30 |
| Mumbai | 19.1°N | 72.9°E | 07:10 | 06:00 |
| Chennai | 13.1°N | 80.3°E | 06:45 | 05:55 |
| Kolkata | 22.6°N | 88.4°E | 06:45 | 05:05 |

**Implementation:** `astronomy/sunrise.go:CalculateSunrise()`

## Adhika Masa (Intercalary Month)

### Definition

An extra lunar month added approximately every 2.7 years to align the lunar and solar calendars.

### Regional Variations

**South India (Amanta):**
- Adhika masa named after following solar month
- No major festivals celebrated
- Considered inauspicious for ceremonies

**North India (Purnimanta):**
- Adhika masa named after solar month
- Similar restrictions on ceremonies
- Religious activities emphasized

### Calculation

```
Adhika Masa occurs when:
- No solar transit (Sankranti) occurs during a lunar month
- Typically happens when Sun is in same sign for 2 lunar months
```

**Implementation:** `astronomy/festivals.go:CalculateAdhikaMasa()`

## Kshaya Masa (Lost Month)

### Definition

A rare occurrence where two solar transits happen within one lunar month, causing a month to be "lost."

**Frequency:** Approximately once every 140 years

**Last Occurrence:** December 1983

**Next Expected:** 2124 CE

## Regional Plugin Architecture

### Plugin System

The system supports regional variations through plugins:

```go
type RegionalPlugin interface {
    GetRegion() string
    ModifyCalculation(ctx context.Context, calc Calculation) (Calculation, error)
    GetFestivals(year int) []Festival
    GetCalendarSystem() CalendarSystem
}
```

### Implemented Plugins

1. **North India Plugin**
   - Purnimanta calendar
   - Region: "north_india"
   - File: `api/examples/north_india_plugin.go`

2. **South India Plugin**
   - Amanta calendar
   - Region: "south_india"
   - File: `api/examples/south_india_plugin.go`

3. **Tamil Nadu Plugin**
   - Solar calendar integration
   - Region: "tamil_nadu"
   - File: `api/examples/tamil_nd_plugin.go`

4. **Kerala Plugin**
   - Malayalam calendar
   - Region: "kerala"
   - File: `api/examples/kerala_plugin.go`

5. **Bengal Plugin**
   - Bengali calendar
   - Region: "bengal"
   - File: `api/examples/bengal_plugin.go`

### Creating Custom Regional Plugins

```go
type CustomRegionPlugin struct {
    region string
    calendarSystem string
}

func (p *CustomRegionPlugin) GetRegion() string {
    return p.region
}

func (p *CustomRegionPlugin) ModifyCalculation(ctx context.Context, calc Calculation) (Calculation, error) {
    // Apply regional modifications
    return calc, nil
}

func (p *CustomRegionPlugin) GetFestivals(year int) []Festival {
    // Return regional festivals
    return []Festival{}
}
```

## Data Sources for Regional Variations

### Primary Sources

1. **Rashtriya Panchang** (National Calendar of India)
   - Published by: Positional Astronomy Centre, Kolkata
   - Authority: Government of India

2. **Regional Panchangams**
   - Kalnirnay (Maharashtra)
   - Vakya Panchangam (South India)
   - Bengali Panjika (Bengal)

3. **Temple Panchangams**
   - Tirupati Temple Panchangam
   - Kashi Vishwanath Temple
   - Jagannath Temple, Puri

### Academic References

4. **Sewell & Dikshit (1896):** *The Indian Calendar* - Comprehensive study of regional variations

5. **Burgess (1858):** *Surya Siddhanta* - Foundation of traditional calculations

6. **Ramakumar (1993):** *Regional Calendars of India* - Detailed regional analysis

## Configuration Examples

### API Request with Regional Settings

```json
{
  "date": "2024-01-15",
  "latitude": 13.0827,
  "longitude": 80.2707,
  "timezone": "Asia/Kolkata",
  "region": "tamil_nadu",
  "calculation_method": "drik",
  "calendar_system": "amanta",
  "locale": "ta"
}
```

### Plugin Registration

```go
manager := api.NewPluginManager()

// Register regional plugins
manager.RegisterPlugin("region_tamil", tamilNaduPlugin)
manager.RegisterPlugin("region_kerala", keralaPlugin)
manager.RegisterPlugin("region_bengal", bengalPlugin)

// Process with regional customization
result := manager.ProcessPanchangamData(ctx, data)
```

## Testing Regional Variations

Validation framework includes regional tests:

```go
// Test Tamil Nadu calculations
tamilRef := ReferenceData{
    Region: "tamil_nadu",
    Date: time.Date(2024, 4, 14, 0, 0, 0, 0, time.UTC),
    // ... reference data
}

suite := validator.ValidateRegionalCalculations(ctx, "tamil_nadu", []ReferenceData{tamilRef})
```

## Summary

Regional variations are handled through:
1. **Pluggable architecture** for regional customizations
2. **Calendar system selection** (Amanta/Purnimanta)
3. **Calculation method choice** (Drik/Vakya)
4. **Location-specific calculations** for sunrise/sunset
5. **Regional festival databases**

All variations are thoroughly tested and validated against established regional sources.

---

*Last Updated: 2025-11-18*
*Maintainer: Panchangam Development Team*
