#!/bin/bash

# Demo script showing various sunrise/sunset calculation examples
# Make sure the panchangam server is running before executing this script

echo "Panchangam Sunrise/Sunset Demo Examples"
echo "========================================="
echo ""

# Function to run demo with description
run_demo() {
    local description="$1"
    shift
    echo "Example: $description"
    echo "Command: go run ./cmd/sunrise-demo $*"
    echo "----------------------------------------"
    go run ./cmd/sunrise-demo "$@"
    echo ""
    echo "Press Enter to continue..."
    read -r
    echo ""
}

# Check if server is running
echo "Checking if panchangam server is running..."
if ! pgrep -f "panchangam" > /dev/null; then
    echo "WARN: server not detected. Please start the server first:"
    echo "   go run ./cmd/server"
    echo ""
    echo "Press Enter when server is running..."
    read -r
fi

echo "Starting demo examples..."
echo ""

# Example 1: Default location (New York)
run_demo "New York (Default)" 

# Example 2: Predefined locations
run_demo "London, UK" -location london

run_demo "Tokyo, Japan" -location tokyo

run_demo "Sydney, Australia" -location sydney

run_demo "Mumbai, India" -location mumbai

run_demo "Cape Town, South Africa" -location capetown

# Example 3: Custom coordinates
run_demo "San Francisco, USA (Custom coords)" -lat 37.7749 -lon -122.4194 -tz "America/Los_Angeles"

run_demo "Paris, France (Custom coords)" -lat 48.8566 -lon 2.3522 -tz "Europe/Paris"

# Example 4: Historical dates
run_demo "London on Winter Solstice 2020" -location london -date 2020-12-21

run_demo "New York on Summer Solstice 2024" -location nyc -date 2024-06-21

run_demo "Tokyo on validation date (Jan 15, 2020)" -location tokyo -date 2020-01-15

# Example 5: Extreme latitudes
run_demo "Reykjavik, Iceland (High latitude)" -lat 64.1466 -lon -21.9426 -tz "Atlantic/Reykjavik"

run_demo "Ushuaia, Argentina (Southern high latitude)" -lat -54.8019 -lon -68.3030 -tz "America/Argentina/Ushuaia"

# Example 6: Different seasons
run_demo "Equator on Equinox" -lat 0.0 -lon 0.0 -tz "UTC" -date 2024-03-20

echo "Demo completed."
echo ""
echo "Tips:"
echo "- Use -location for quick predefined locations"
echo "- Specify -date for historical dates (YYYY-MM-DD)"
echo "- All times are returned in UTC for consistency"
echo "- Day length is calculated automatically"
echo "- The algorithm handles polar regions correctly"
echo ""
echo "For more examples, see cmd/sunrise-demo/README.md"
