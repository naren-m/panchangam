# Panchangam CLI User Guide

A comprehensive guide for using the Panchangam CLI tool for astronomical calculations based on Hindu calendar systems.

## Table of Contents
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Command Reference](#command-reference)
- [Common Use Cases](#common-use-cases)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/naren-m/panchangam.git
cd panchangam

# Build the CLI
go build -o panchangam-cli ./cmd/panchangam-cli

# Optional: Install to system path
sudo mv panchangam-cli /usr/local/bin/
```

### Using Go Install
```bash
go install github.com/naren-m/panchangam/cmd/panchangam-cli@latest
```

### Verify Installation
```bash
panchangam-cli version
```

## Quick Start

### 1. Get Today's Tithi
```bash
# For Mumbai
panchangam-cli tithi -l mumbai

# For custom location
panchangam-cli tithi --lat 19.0760 --lon 72.8777
```

### 2. Check Sunrise/Sunset Times
```bash
# Basic sun times
panchangam-cli sun -l london

# Detailed with solar noon
panchangam-cli sun -l london --detailed
```

### 3. Check Service Health
```bash
panchangam-cli health
```

## Core Concepts

### Astronomical Terms

**Tithi (तिथि)**: Lunar day in Hindu calendar. Each tithi is approximately 12° of angular distance between the Moon and Sun.

**Nakshatra (नक्षत्र)**: Lunar mansion or constellation. The zodiac is divided into 27 nakshatras.

**Yoga (योग)**: A combination of the Sun and Moon's positions. There are 27 yogas.

**Karana (करण)**: Half of a tithi. There are 11 karanas in total.

**Paksha (पक्ष)**: Lunar fortnight. Shukla Paksha (waxing) and Krishna Paksha (waning).

### Location Handling

The CLI supports two ways to specify locations:

1. **Predefined Locations**: Use `-l` flag with city code
2. **Custom Coordinates**: Use `--lat` and `--lon` flags

### Output Formats

- **table** (default): Human-readable format with Unicode formatting
- **json**: Machine-readable JSON format
- **yaml**: YAML format for configuration files
- **csv**: CSV format for spreadsheet import

## Command Reference

### Essential Commands

#### `tithi` - Lunar Day Calculation
```bash
# Basic usage
panchangam-cli tithi -l mumbai

# Specific date
panchangam-cli tithi -l mumbai -d 2024-06-21

# Detailed information
panchangam-cli tithi -l mumbai --detailed

# JSON output
panchangam-cli tithi -l mumbai -o json
```

**Output includes:**
- Tithi number (1-30)
- Tithi name (Sanskrit)
- Tithi type (Nanda, Bhadra, Jaya, Rikta, Purna)
- Paksha (Shukla/Krishna)
- Start/End times (with --detailed)
- Duration and Moon-Sun angle (with --detailed)

#### `sun` - Sun Timing Information
```bash
# Basic sunrise/sunset
panchangam-cli sun -l tokyo

# Detailed with solar noon
panchangam-cli sun -l tokyo --detailed

# Custom location
panchangam-cli sun --lat 35.6762 --lon 139.6503 --tz "Asia/Tokyo"
```

**Output includes:**
- Sunrise time
- Sunset time
- Day length
- Solar noon (with --detailed)
- Night length (with --detailed)

#### `locations` - List Available Cities
```bash
panchangam-cli locations
```

**Available location codes:**
- `nyc` - New York, USA
- `london` - London, UK
- `tokyo` - Tokyo, Japan
- `sydney` - Sydney, Australia
- `mumbai` - Mumbai, India
- `capetown` - Cape Town, South Africa
- `paris` - Paris, France
- `moscow` - Moscow, Russia
- `beijing` - Beijing, China
- `cairo` - Cairo, Egypt
- `rio` - Rio de Janeiro, Brazil
- `losangeles` - Los Angeles, USA

### Service Commands

#### `version` - Version Information
```bash
# Table format
panchangam-cli version

# JSON format
panchangam-cli version -o json
```

#### `health` - Service Health Check
```bash
# Check service status
panchangam-cli health

# JSON output for monitoring
panchangam-cli health -o json
```

#### `validate` - Validate Server Connection
```bash
panchangam-cli validate
```

### Upcoming Features (Framework Ready)

These commands have the framework implemented and will be fully functional in future releases:

- `nakshatra` - Lunar mansion calculations
- `yoga` - Combined Sun/Moon position calculations
- `karana` - Half-tithi calculations
- `ephemeris` - Planetary position data
- `events` - Festivals and special occasions
- `muhurta` - Auspicious time periods
- `range` - Multi-day calculations

## Common Use Cases

### Daily Panchangam Check

Create a morning routine script:
```bash
#!/bin/bash
# morning-panchangam.sh

echo "🌅 Good Morning! Today's Panchangam"
echo "===================================="
date
echo

echo "📍 Location: Mumbai"
echo

echo "🌙 Tithi:"
panchangam-cli tithi -l mumbai | grep -E "Name|Type|Paksha"
echo

echo "☀️ Sun Times:"
panchangam-cli sun -l mumbai | grep -E "Sunrise|Sunset|Day Length"
```

### Compare Sunrise Across Cities

```bash
#!/bin/bash
# sunrise-comparison.sh

cities=("nyc" "london" "tokyo" "mumbai" "sydney")

echo "🌅 Sunrise Times Around the World"
echo "================================="
echo

for city in "${cities[@]}"; do
    result=$(panchangam-cli sun -l $city | grep "Sunrise" | awk '{print $2}')
    echo "$city: $result"
done
```

### Monthly Tithi Calendar

```bash
#!/bin/bash
# monthly-tithis.sh

year=2024
month=06

echo "📅 Tithi Calendar for $month/$year"
echo "=================================="

for day in {01..30}; do
    date="$year-$month-$day"
    tithi=$(panchangam-cli tithi -l mumbai -d $date -o json | jq -r '.name')
    echo "$date: $tithi"
done
```

### Export Data for Analysis

```bash
# Export week's data to JSON
for i in {0..6}; do
    date=$(date -d "+$i days" +%Y-%m-%d)
    panchangam-cli tithi -l mumbai -d $date -o json
done > week_data.json

# Export to CSV for spreadsheet
echo "Date,Tithi,Type,Sunrise,Sunset" > panchangam_data.csv
for i in {0..30}; do
    date=$(date -d "+$i days" +%Y-%m-%d)
    # Parse and format data (example)
    panchangam-cli get -l mumbai -d $date -o json | \
    jq -r '[.panchangam_data.date, .panchangam_data.tithi, "type", .panchangam_data.sunrise_time, .panchangam_data.sunset_time] | @csv' \
    >> panchangam_data.csv
done
```

## Advanced Usage

### Custom Server Configuration

```bash
# Connect to remote server
panchangam-cli tithi -s "api.example.com:8080" -l mumbai

# With custom timeout
panchangam-cli tithi -s "remote:8080" -t 30s -l mumbai
```

### Timezone Handling

```bash
# Explicit timezone specification
panchangam-cli tithi --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"

# List available timezones (Linux/Mac)
timedatectl list-timezones | grep Asia
```

### Integration with Other Tools

```bash
# Send daily tithi to Slack
tithi=$(panchangam-cli tithi -l mumbai -o json | jq -r '.name')
curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\"Today's Tithi: $tithi\"}" \
    YOUR_SLACK_WEBHOOK_URL

# Log to file with timestamp
echo "$(date): $(panchangam-cli tithi -l mumbai)" >> ~/panchangam.log
```

### Programmatic Usage

```python
# Python example
import subprocess
import json

def get_tithi(location):
    cmd = ["panchangam-cli", "tithi", "-l", location, "-o", "json"]
    result = subprocess.run(cmd, capture_output=True, text=True)
    return json.loads(result.stdout)

tithi_data = get_tithi("mumbai")
print(f"Today's Tithi: {tithi_data['name']} ({tithi_data['type']})")
```

```javascript
// Node.js example
const { exec } = require('child_process');

function getTithi(location) {
    return new Promise((resolve, reject) => {
        exec(`panchangam-cli tithi -l ${location} -o json`, (error, stdout) => {
            if (error) reject(error);
            else resolve(JSON.parse(stdout));
        });
    });
}

getTithi('mumbai').then(data => {
    console.log(`Today's Tithi: ${data.name} (${data.type})`);
});
```

## Troubleshooting

### Common Issues

#### 1. Connection Refused
```bash
# Error: connection refused
# Solution: Check if server is running
panchangam-cli validate

# Try with different server
panchangam-cli tithi -s "localhost:8080" -l mumbai
```

#### 2. Invalid Date Format
```bash
# Error: invalid date format
# Correct format: YYYY-MM-DD
panchangam-cli tithi -d 2024-06-21  # ✓ Correct
panchangam-cli tithi -d 06/21/2024  # ✗ Wrong
```

#### 3. Unknown Location
```bash
# Error: unknown location
# Solution: Use 'locations' command to see available codes
panchangam-cli locations

# Or use coordinates
panchangam-cli tithi --lat 19.0760 --lon 72.8777
```

#### 4. Timezone Issues
```bash
# If timezone is not recognized
# Solution: Use full timezone identifier
panchangam-cli sun --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"

# List valid timezones
timedatectl list-timezones
```

### Debug Mode

```bash
# Enable verbose output
panchangam-cli tithi -l mumbai -v

# Enable debug output
panchangam-cli tithi -l mumbai --debug
```

### Performance Issues

```bash
# Increase timeout for slow connections
panchangam-cli tithi -l mumbai -t 30s

# Use JSON output for faster parsing
panchangam-cli tithi -l mumbai -o json | jq '.name'
```

## Best Practices

### 1. Use Location Presets When Possible
Location presets include timezone information automatically:
```bash
# Preferred
panchangam-cli tithi -l mumbai

# Instead of
panchangam-cli tithi --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"
```

### 2. Cache Results for Multiple Queries
```bash
# Save result and query multiple times
panchangam-cli tithi -l mumbai -o json > today_tithi.json
tithi_name=$(jq -r '.name' today_tithi.json)
tithi_type=$(jq -r '.type' today_tithi.json)
```

### 3. Use JSON for Scripting
```bash
# Easy to parse with jq
panchangam-cli sun -l tokyo -o json | jq -r '.sunrise'
```

### 4. Handle Errors Gracefully
```bash
# Check command success
if panchangam-cli tithi -l mumbai > /dev/null 2>&1; then
    echo "Success"
else
    echo "Failed to get tithi"
fi
```

## Additional Resources

- **GitHub Repository**: https://github.com/naren-m/panchangam
- **API Documentation**: See `/docs` in the repository
- **Issue Tracking**: https://github.com/naren-m/panchangam/issues

## Contributing

To contribute to the CLI development:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

See the LICENSE file in the main repository for details.