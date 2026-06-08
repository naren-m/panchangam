# Sunrise/Sunset Demo Client

A simple command-line client to demonstrate the sunrise/sunset calculation API.

## Quick Start

### 1. Start the server
```bash
# In the project root
go run ./cmd/server
```

### 2. Run the demo client
```bash
# Basic usage with default location (New York)
go run ./cmd/sunrise-demo

# Specify custom coordinates
go run ./cmd/sunrise-demo -lat 51.5074 -lon -0.1278 -tz "Europe/London"

# Use predefined locations
go run ./cmd/sunrise-demo -location london
go run ./cmd/sunrise-demo -location tokyo
go run ./cmd/sunrise-demo -location mumbai
```

## Usage Examples

### Predefined Locations
```bash
# New York
go run ./cmd/sunrise-demo -location nyc

# London
go run ./cmd/sunrise-demo -location london

# Tokyo
go run ./cmd/sunrise-demo -location tokyo

# Sydney
go run ./cmd/sunrise-demo -location sydney

# Mumbai
go run ./cmd/sunrise-demo -location mumbai

# Cape Town
go run ./cmd/sunrise-demo -location capetown
```

### Custom Coordinates
```bash
# San Francisco
go run ./cmd/sunrise-demo -lat 37.7749 -lon -122.4194 -tz "America/Los_Angeles"

# Paris
go run ./cmd/sunrise-demo -lat 48.8566 -lon 2.3522 -tz "Europe/Paris"

# Chennai
go run ./cmd/sunrise-demo -lat 13.0827 -lon 80.2707 -tz "Asia/Kolkata"
```

### Historical Dates
```bash
# January 15, 2020 (validation date)
go run ./cmd/sunrise-demo -location london -date 2020-01-15

# Summer solstice
go run ./cmd/sunrise-demo -location nyc -date 2024-06-21

# Winter solstice
go run ./cmd/sunrise-demo -location london -date 2024-12-21
```

## Command Line Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-address` | gRPC server address | `localhost:8080` | `-address localhost:9090` |
| `-date` | Date in YYYY-MM-DD format | Today | `-date 2024-06-21` |
| `-lat` | Latitude (-90 to 90) | `40.7128` | `-lat 51.5074` |
| `-lon` | Longitude (-180 to 180) | `-74.0060` | `-lon -0.1278` |
| `-tz` | Timezone | `America/New_York` | `-tz "Asia/Tokyo"` |
| `-location` | Predefined location | None | `-location london` |

## Available Predefined Locations

| Location | Coordinates | Timezone |
|----------|-------------|----------|
| `nyc`, `newyork` | 40.7128°N, 74.0060°W | America/New_York |
| `london` | 51.5074°N, 0.1278°W | Europe/London |
| `tokyo` | 35.6762°N, 139.6503°E | Asia/Tokyo |
| `sydney` | 33.8688°S, 151.2093°E | Australia/Sydney |
| `mumbai` | 19.0760°N, 72.8777°E | Asia/Kolkata |
| `capetown` | 33.9249°S, 18.4241°E | Africa/Johannesburg |

## Sample Output

```
Sunrise/Sunset Calculator
-------------------------
Date: 2024-07-17
Location: 51.5074°N, -0.1278°E
Timezone: Europe/London
Server: localhost:8080
-------------------------

Results:
  Sunrise: 05:15:32
  Sunset:  20:58:45
Day Length: 15h43m13s

Traditional Panchangam Data:
- Tithi: Some Tithi
- Nakshatra: Some Nakshatra
- Yoga: Some Yoga
- Karana: Some Karana

Events:
- Some Event 1 at 08:00:00
- Some Event 2 at 12:00:00

Calculation completed successfully.
```

## Building

```bash
# Build the demo client
go build -o sunrise-demo ./cmd/sunrise-demo

# Run the binary
./sunrise-demo -location tokyo -date 2024-06-21
```

## Error Handling

The client handles various error conditions:
- Invalid coordinates (lat/lon out of range)
- Invalid date format
- Server connection errors
- gRPC timeout errors
- Unknown predefined locations

## Integration with Your Applications

This demo shows how to:
1. Connect to the gRPC service
2. Create properly formatted requests
3. Handle responses and errors
4. Parse and display sunrise/sunset times
5. Calculate day length from the results

You can use this as a reference for integrating the sunrise/sunset API into your own applications.
