# Panchangam CLI Commands Reference

Complete reference for all available commands in the Panchangam CLI.

## Command Overview

```
panchangam-cli [command] [flags]
```

### Core Astronomical Commands

| Command | Description | Status |
|---------|-------------|---------|
| `tithi` | Calculate Tithi (lunar day) | ✅ Working |
| `nakshatra` | Calculate Nakshatra (lunar mansion) | 📋 Framework Ready |
| `yoga` | Calculate Yoga | 📋 Framework Ready |
| `karana` | Calculate Karana (half-tithi) | 📋 Framework Ready |
| `sun` | Detailed sun timing information | ✅ Working |
| `ephemeris` | Planetary position data | 📋 Framework Ready |

### Calendar & Event Commands

| Command | Description | Status |
|---------|-------------|---------|
| `events` | Festivals and special occasions | 📋 Framework Ready |
| `muhurta` | Auspicious time periods | 📋 Framework Ready |
| `range` | Multi-day Panchangam data | 📋 Framework Ready |

### Service Commands

| Command | Description | Status |
|---------|-------------|---------|
| `get` | Basic Panchangam data (legacy) | ✅ Working |
| `health` | Service health check | ✅ Working |
| `version` | Version information | ✅ Working |
| `validate` | Validate server connectivity | ✅ Working |
| `benchmark` | Performance testing | ✅ Working |
| `locations` | List predefined locations | ✅ Working |

## Detailed Command Reference

### `tithi` - Tithi Calculation

Calculate the lunar day (Tithi) for a given date and location.

```bash
panchangam-cli tithi [flags]
```

**Flags:**
- `-d, --date string` - Date in YYYY-MM-DD format (default: today)
- `-l, --location string` - Predefined location (e.g., mumbai, london)
- `--lat float` - Latitude (-90 to 90)
- `--lon float` - Longitude (-180 to 180)
- `--tz string` - Timezone (default: "Asia/Kolkata")
- `--detailed` - Show detailed Tithi information

**Examples:**
```bash
# Basic tithi for Mumbai
panchangam-cli tithi -l mumbai

# Detailed tithi for specific date
panchangam-cli tithi -d 2024-06-21 -l mumbai --detailed

# Custom coordinates
panchangam-cli tithi --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"

# JSON output
panchangam-cli tithi -l mumbai -o json
```

### `sun` - Sun Timing Information

Calculate detailed sun timing information including sunrise, sunset, solar noon, and day length.

```bash
panchangam-cli sun [flags]
```

**Flags:**
- `-d, --date string` - Date in YYYY-MM-DD format (default: today)
- `-l, --location string` - Predefined location
- `--lat float` - Latitude (-90 to 90)
- `--lon float` - Longitude (-180 to 180)
- `--tz string` - Timezone
- `--detailed` - Show solar noon and night length

**Examples:**
```bash
# Basic sun times for Tokyo
panchangam-cli sun -l tokyo

# Detailed sun information
panchangam-cli sun -l london --detailed

# Summer solstice in Paris
panchangam-cli sun -l paris -d 2024-06-21 --detailed
```

### `nakshatra` - Nakshatra Calculation

Calculate the lunar mansion (Nakshatra) with pada, deity, and characteristics.

```bash
panchangam-cli nakshatra [flags]
```

**Status:** Framework ready, implementation pending

### `yoga` - Yoga Calculation

Calculate Yoga based on combined Sun and Moon positions.

```bash
panchangam-cli yoga [flags]
```

**Status:** Framework ready, implementation pending

### `karana` - Karana Calculation

Calculate Karana (half-tithi) information.

```bash
panchangam-cli karana [flags]
```

**Status:** Framework ready, implementation pending

### `ephemeris` - Planetary Positions

Get planetary position data from ephemeris providers.

```bash
panchangam-cli ephemeris [flags]
```

**Flags:**
- `-d, --date string` - Date in YYYY-MM-DD format
- `-p, --planet string` - Planet name (sun, moon, mars, mercury, jupiter, venus, saturn, all)
- `--provider string` - Ephemeris provider (swiss, jpl)
- `--detailed` - Show detailed ephemeris information

**Status:** Framework ready, implementation pending

### `events` - Festivals and Events

Get festivals, religious events, and special occasions.

```bash
panchangam-cli events [flags]
```

**Flags:**
- `-d, --date string` - Date in YYYY-MM-DD format
- `-l, --location string` - Predefined location
- `--type string` - Event type (festival, ekadashi, amavasya, purnima, all)
- `--region string` - Regional variation (tamil_nadu, kerala, bengal)

**Status:** Framework ready, implementation pending

### `muhurta` - Auspicious Timings

Calculate auspicious and inauspicious time periods.

```bash
panchangam-cli muhurta [flags]
```

**Flags:**
- `-d, --date string` - Date in YYYY-MM-DD format
- `-l, --location string` - Predefined location
- `--purpose string` - Purpose (marriage, business, travel, all)
- `--quality string` - Quality filter (auspicious, inauspicious, all)

**Status:** Framework ready, implementation pending

### `range` - Date Range Calculations

Get Panchangam data for multiple consecutive dates.

```bash
panchangam-cli range [flags]
```

**Flags:**
- `--start string` - Start date in YYYY-MM-DD format
- `--end string` - End date in YYYY-MM-DD format
- `-l, --location string` - Predefined location
- `--method string` - Calculation method (drik, vakya)
- `--region string` - Regional variation

**Status:** Framework ready, implementation pending

### `health` - Service Health Check

Check the health status of the Panchangam service and ephemeris providers.

```bash
panchangam-cli health [flags]
```

**Examples:**
```bash
# Basic health check
panchangam-cli health

# JSON output for monitoring
panchangam-cli health -o json
```

### `version` - Version Information

Display version information and supported features.

```bash
panchangam-cli version [flags]
```

**Examples:**
```bash
# Table format
panchangam-cli version

# JSON format for parsing
panchangam-cli version -o json
```

## Global Flags

These flags are available for all commands:

- `-s, --server string` - gRPC server address (default: "localhost:8080")
- `-o, --output string` - Output format: table, json, yaml, csv (default: "table")
- `-t, --timeout duration` - Request timeout (default: 10s)
- `-v, --verbose` - Enable verbose output
- `--debug` - Enable debug output
- `-h, --help` - Help for any command

## Output Formats

### Table (Default)
Human-readable table format with Unicode characters and formatting.

### JSON
Machine-readable JSON format for parsing and integration.
```bash
panchangam-cli tithi -l mumbai -o json
```

### YAML
YAML format for configuration files and readable data exchange.
```bash
panchangam-cli sun -l tokyo -o yaml
```

### CSV
Comma-separated values for spreadsheet import (where applicable).
```bash
panchangam-cli range --start 2024-06-01 --end 2024-06-07 -o csv
```

## Location Presets

Use `-l` or `--location` flag with these preset codes:

- `nyc` - New York, USA (40.7128°N, 74.0060°W)
- `london` - London, UK (51.5074°N, 0.1278°W)
- `tokyo` - Tokyo, Japan (35.6762°N, 139.6503°E)
- `sydney` - Sydney, Australia (33.8688°S, 151.2093°E)
- `mumbai` - Mumbai, India (19.0760°N, 72.8777°E)
- `capetown` - Cape Town, South Africa (33.9249°S, 18.4241°E)
- `paris` - Paris, France (48.8566°N, 2.3522°E)
- `moscow` - Moscow, Russia (55.7558°N, 37.6176°E)
- `beijing` - Beijing, China (39.9042°N, 116.4074°E)
- `cairo` - Cairo, Egypt (30.0444°N, 31.2357°E)
- `rio` - Rio de Janeiro, Brazil (22.9068°S, 43.1729°W)
- `losangeles` - Los Angeles, USA (34.0522°N, 118.2437°W)

## Examples by Use Case

### Daily Panchangam Check
```bash
# Morning routine check for Mumbai
panchangam-cli tithi -l mumbai
panchangam-cli sun -l mumbai --detailed
panchangam-cli events -l mumbai
```

### Planning Travel
```bash
# Check muhurtas for travel
panchangam-cli muhurta --purpose travel -d 2024-06-25

# Check tithi for multiple days
for day in 25 26 27; do
  panchangam-cli tithi -d 2024-06-$day -l mumbai
done
```

### Festival Calendar
```bash
# Get Tamil Nadu festivals
panchangam-cli events --region tamil_nadu

# Export to CSV
panchangam-cli events --region kerala -o csv > kerala_festivals.csv
```

### Astronomical Research
```bash
# Track sun times throughout the year
for month in {01..12}; do
  panchangam-cli sun -d 2024-$month-15 -l mumbai --detailed
done

# Compare planetary positions
panchangam-cli ephemeris --planet sun -o json
panchangam-cli ephemeris --planet moon -o json
```