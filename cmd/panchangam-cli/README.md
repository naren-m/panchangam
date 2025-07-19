# Panchangam CLI

A comprehensive command-line interface for astronomical calculations based on Hindu calendar systems.

## 📚 Documentation

- **[User Guide](USER_GUIDE.md)** - Complete usage guide with examples and best practices
- **[Quick Reference](QUICK_REFERENCE.md)** - Command cheat sheet for immediate use
- **[Commands Reference](COMMANDS.md)** - Detailed documentation for all commands

## 🚀 Enhanced Features (New!)

- 🌙 **Tithi Calculations** - Lunar day with detailed timing and characteristics
- 🌅 **Enhanced Sun Times** - Sunrise, sunset, solar noon, and day length
- 🏥 **Health Monitoring** - Service status and ephemeris health checks
- 📊 **Multiple Output Formats** - Table, JSON, YAML, CSV support
- 🌍 **Global Coverage** - Predefined locations and custom coordinates
- 📋 **Framework Ready** - Commands prepared for Nakshatra, Yoga, Karana, Events, Muhurtas

## ⚡ Quick Start

```bash
# Build the CLI
go build -o panchangam-cli .

# Get today's Tithi for Mumbai
./panchangam-cli tithi -l mumbai

# Get detailed sun times for London  
./panchangam-cli sun -l london --detailed

# Check service health
./panchangam-cli health

# See all commands
./panchangam-cli --help
```

## 📋 Available Commands

| Command | Status | Description |
|---------|--------|-------------|
| `tithi` | ✅ **Working** | Calculate Tithi (lunar day) with timing |
| `sun` | ✅ **Working** | Detailed sun timing information |
| `health` | ✅ **Working** | Service health and status check |
| `version` | ✅ **Working** | Version and feature information |
| `get` | ✅ **Working** | Basic Panchangam data (legacy) |
| `locations` | ✅ **Working** | List predefined city locations |
| `validate` | ✅ **Working** | Validate server connectivity |
| `benchmark` | ✅ **Working** | Performance testing |
| `nakshatra` | 📋 Framework Ready | Lunar mansion calculations |
| `yoga` | 📋 Framework Ready | Combined Sun/Moon positions |
| `karana` | 📋 Framework Ready | Half-tithi calculations |
| `ephemeris` | 📋 Framework Ready | Planetary position data |
| `events` | 📋 Framework Ready | Festivals and special occasions |
| `muhurta` | 📋 Framework Ready | Auspicious time periods |
| `range` | 📋 Framework Ready | Multi-day calculations |

## Features

- 🌍 **Predefined Locations**: Quick access to major cities worldwide
- 📅 **Historical Dates**: Test with any past or future date
- 🎯 **Multiple Output Formats**: Table, JSON, and YAML output
- 🔍 **Server Validation**: Test connectivity and basic functionality
- 📊 **Benchmarking**: Performance testing capabilities
- ⚡ **Custom Coordinates**: Use any latitude/longitude

## Installation

```bash
# Build the CLI
go build -o panchangam-cli cmd/panchangam-cli/main.go

# Or use via go run
go run cmd/panchangam-cli/main.go [command]
```

## Quick Start

```bash
# 1. Start the server (in another terminal)
make run

# 2. Validate connection
./panchangam-cli validate

# 3. Get data for New York
./panchangam-cli get -l nyc

# 4. See all available locations
./panchangam-cli locations
```

## Commands

### `get` - Get Panchangam Data

Retrieve sunrise/sunset times and panchangam data for a specific date and location.

```bash
# Using predefined locations
./panchangam-cli get -l london
./panchangam-cli get -l tokyo -d 2024-06-21

# Using custom coordinates
./panchangam-cli get --lat 37.7749 --lon -122.4194 --tz "America/Los_Angeles"

# With different output formats
./panchangam-cli get -l mumbai -o json
./panchangam-cli get -l sydney -o yaml
```

**Options:**
- `-l, --location`: Predefined location code (see `locations` command)
- `-d, --date`: Date in YYYY-MM-DD format (default: today)
- `--lat`: Latitude (-90 to 90)
- `--lon`: Longitude (-180 to 180)
- `--tz`: Timezone (e.g., America/New_York)
- `--region`: Regional system (e.g., Tamil Nadu)
- `--method`: Calculation method (e.g., Drik, Vakya)
- `--locale`: Language/locale (e.g., en, ta)

### `locations` - List Available Locations

Show all predefined locations with their coordinates and timezones.

```bash
./panchangam-cli locations
```

**Available Locations:**
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

### `validate` - Validate Server Connection

Test connectivity to the gRPC server and validate basic functionality.

```bash
./panchangam-cli validate
```

### `benchmark` - Performance Testing

Run performance benchmarks against the gRPC server.

```bash
# Default benchmark (100 requests, 10 workers)
./panchangam-cli benchmark

# Custom benchmark
./panchangam-cli benchmark -n 1000 -w 20
```

**Options:**
- `-n, --requests`: Number of requests to make (default: 100)
- `-w, --workers`: Number of concurrent workers (default: 10)

## Global Options

- `-s, --server`: gRPC server address (default: localhost:8080)
- `-o, --output`: Output format - table, json, yaml (default: table)
- `-t, --timeout`: Request timeout (default: 10s)

## Examples

### Basic Usage

```bash
# Get today's data for London
./panchangam-cli get -l london

# Get historical data for New York
./panchangam-cli get -l nyc -d 2020-01-15

# Get data for custom location
./panchangam-cli get --lat 51.5074 --lon -0.1278 --tz "Europe/London"
```

### Different Output Formats

```bash
# Table format (default)
./panchangam-cli get -l tokyo

# JSON format
./panchangam-cli get -l tokyo -o json

# YAML format
./panchangam-cli get -l tokyo -o yaml
```

### Testing and Validation

```bash
# Test server connectivity
./panchangam-cli validate

# Run performance benchmark
./panchangam-cli benchmark -n 50 -w 5

# Test with different server
./panchangam-cli -s localhost:9090 validate
```

### Advanced Usage

```bash
# Get data with all optional fields
./panchangam-cli get -l mumbai -d 2024-06-21 \
  --region "Maharashtra" \
  --method "Drik" \
  --locale "en"

# Benchmark with custom parameters
./panchangam-cli benchmark -n 500 -w 25 -t 30s
```

## Sample Output

### Table Format (Default)
```
🌅 Panchangam Data
═══════════════════════════════════════
📅 Date: 2024-07-17
📍 Location: 51.5074°N, -0.1278°E
🌐 Timezone: Europe/London
═══════════════════════════════════════

📊 Sunrise/Sunset Times (UTC):
┌─────────────────────────────────┐
│ 🌅 Sunrise: 04:15:32           │
│ 🌇 Sunset:  20:58:45           │
└─────────────────────────────────┘
☀️  Day Length: 16h43m13s

📜 Traditional Panchangam:
• Tithi: Some Tithi
• Nakshatra: Some Nakshatra
• Yoga: Some Yoga
• Karana: Some Karana
```

### JSON Format
```json
{
  "panchangam_data": {
    "date": "2024-07-17",
    "tithi": "Some Tithi",
    "nakshatra": "Some Nakshatra",
    "yoga": "Some Yoga",
    "karana": "Some Karana",
    "sunrise_time": "04:15:32",
    "sunset_time": "20:58:45",
    "events": [
      {
        "name": "Some Event 1",
        "time": "08:00:00"
      }
    ]
  }
}
```

## Integration with Other Tools

### Shell Scripts
```bash
#!/bin/bash
# Get sunrise time for multiple cities
for city in nyc london tokyo sydney; do
    echo "=== $city ==="
    ./panchangam-cli get -l $city -o json | jq -r '.panchangam_data.sunrise_time'
done
```

### JSON Processing with jq
```bash
# Extract just the sunrise time
./panchangam-cli get -l london -o json | jq -r '.panchangam_data.sunrise_time'

# Get day length calculation
./panchangam-cli get -l tokyo -o json | jq -r '.panchangam_data | "\(.sunrise_time) to \(.sunset_time)"'
```

## Error Handling

The CLI provides helpful error messages for common issues:

- **Connection errors**: Server not running or wrong address
- **Invalid coordinates**: Latitude/longitude out of range
- **Invalid dates**: Wrong date format or invalid date
- **Unknown locations**: Invalid location code
- **Timeout errors**: Server taking too long to respond

## Development and Testing

```bash
# Run with verbose output
./panchangam-cli -v get -l london

# Test with local development server
./panchangam-cli -s localhost:8080 validate

# Benchmark server performance
./panchangam-cli benchmark -n 100 -w 10
```

## Building

```bash
# Build for current platform
go build -o panchangam-cli cmd/panchangam-cli/main.go

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o panchangam-cli-linux cmd/panchangam-cli/main.go
GOOS=windows GOARCH=amd64 go build -o panchangam-cli.exe cmd/panchangam-cli/main.go
GOOS=darwin GOARCH=amd64 go build -o panchangam-cli-mac cmd/panchangam-cli/main.go
```

This CLI provides a comprehensive interface for testing and interacting with your Panchangam gRPC service!