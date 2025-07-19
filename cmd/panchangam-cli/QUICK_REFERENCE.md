# Panchangam CLI Quick Reference

## Installation
```bash
go build -o panchangam-cli ./cmd/panchangam-cli
```

## Essential Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `tithi` | Get lunar day | `panchangam-cli tithi -l mumbai` |
| `sun` | Sunrise/sunset times | `panchangam-cli sun -l tokyo --detailed` |
| `health` | Check service status | `panchangam-cli health` |
| `version` | Show version info | `panchangam-cli version` |
| `locations` | List city presets | `panchangam-cli locations` |

## Quick Examples

### Daily Basics
```bash
# Today's Tithi for Mumbai
panchangam-cli tithi -l mumbai

# Detailed sun times for London
panchangam-cli sun -l london --detailed

# Check what's available
panchangam-cli --help
```

### Custom Location
```bash
# Using coordinates
panchangam-cli tithi --lat 19.0760 --lon 72.8777 --tz "Asia/Kolkata"

# Specific date
panchangam-cli sun -l tokyo -d 2024-06-21
```

### Output Formats
```bash
# JSON output
panchangam-cli tithi -l mumbai -o json

# YAML output  
panchangam-cli version -o yaml
```

## Location Codes
| Code | City | Timezone |
|------|------|----------|
| `mumbai` | Mumbai, India | Asia/Kolkata |
| `london` | London, UK | Europe/London |
| `tokyo` | Tokyo, Japan | Asia/Tokyo |
| `nyc` | New York, USA | America/New_York |
| `sydney` | Sydney, Australia | Australia/Sydney |

**See all**: `panchangam-cli locations`

## Global Flags
| Flag | Description | Example |
|------|-------------|---------|
| `-o` | Output format (table/json/yaml) | `-o json` |
| `-s` | Server address | `-s "remote:8080"` |
| `-t` | Timeout | `-t 30s` |
| `-v` | Verbose output | `-v` |
| `--debug` | Debug mode | `--debug` |

## Common Patterns

### Morning Check Script
```bash
#!/bin/bash
echo "Today's Panchangam:"
panchangam-cli tithi -l mumbai
panchangam-cli sun -l mumbai
```

### JSON Parsing
```bash
# Get just the tithi name
panchangam-cli tithi -l mumbai -o json | jq -r '.name'

# Get sunrise time
panchangam-cli sun -l london -o json | jq -r '.sunrise'
```

### Multi-City Comparison
```bash
for city in mumbai london tokyo; do
  echo "$city: $(panchangam-cli sun -l $city | grep Sunrise)"
done
```

## Troubleshooting
- **Connection issues**: `panchangam-cli validate`
- **Unknown location**: `panchangam-cli locations`
- **Wrong date format**: Use `YYYY-MM-DD`
- **Timezone errors**: Use full identifier like `Asia/Kolkata`

## Help
```bash
# General help
panchangam-cli --help

# Command-specific help
panchangam-cli tithi --help
panchangam-cli sun --help
```