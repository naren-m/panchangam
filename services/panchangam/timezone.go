package panchangam

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimezoneParser handles parsing of various timezone formats
type TimezoneParser struct{}

// NewTimezoneParser creates a new timezone parser
func NewTimezoneParser() *TimezoneParser {
	return &TimezoneParser{}
}

// ParseTimezone parses a timezone string and returns a *time.Location
// Supports:
// 1. IANA timezone identifiers (e.g., "Asia/Kolkata", "America/New_York")
// 2. UTC offset formats (e.g., "+05:30", "-08:00", "UTC+5:30")
// 3. "UTC" or empty string (defaults to UTC)
func (tp *TimezoneParser) ParseTimezone(tz string) (*time.Location, error) {
	tz = strings.TrimSpace(tz)

	if tz == "" {
		return time.UTC, nil
	}

	if tz == "UTC" || tz == "GMT" {
		return time.UTC, nil
	}

	loc, err := time.LoadLocation(tz)
	if err == nil {
		return loc, nil
	}

	loc, err = tp.parseUTCOffset(tz)
	if err == nil {
		return loc, nil
	}

	return nil, fmt.Errorf("invalid timezone '%s': must be a valid IANA timezone identifier (e.g., 'Asia/Kolkata') or UTC offset (e.g., '+05:30', '-08:00')", tz)
}

// parseUTCOffset parses UTC offset formats like "+05:30", "-08:00", "UTC+5:30"
func (tp *TimezoneParser) parseUTCOffset(offset string) (*time.Location, error) {
	// Regular expression to match UTC offset formats
	// Matches: +05:30, -08:00, UTC+5:30, GMT-8, etc.
	re := regexp.MustCompile(`^(UTC|GMT)?([+-])(\d{1,2}):?(\d{2})?$`)
	matches := re.FindStringSubmatch(offset)

	if matches == nil {
		return nil, fmt.Errorf("invalid UTC offset format")
	}

	sign := matches[2]
	hoursStr := matches[3]
	minutesStr := matches[4]

	// Parse hours
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 0 || hours > 14 {
		return nil, fmt.Errorf("invalid hours in UTC offset: must be between 0 and 14")
	}

	// Parse minutes (default to 0 if not provided)
	minutes := 0
	if minutesStr != "" {
		minutes, err = strconv.Atoi(minutesStr)
		if err != nil || minutes < 0 || minutes > 59 {
			return nil, fmt.Errorf("invalid minutes in UTC offset: must be between 0 and 59")
		}
	}

	// Calculate total offset in seconds
	totalSeconds := hours*3600 + minutes*60
	if sign == "-" {
		totalSeconds = -totalSeconds
	}

	// Validate total offset is within reasonable range (-14h to +14h)
	if totalSeconds < -14*3600 || totalSeconds > 14*3600 {
		return nil, fmt.Errorf("UTC offset out of range: must be between -14:00 and +14:00")
	}

	// Create a fixed timezone with the offset
	name := fmt.Sprintf("UTC%s%02d:%02d", sign, hours, minutes)
	return time.FixedZone(name, totalSeconds), nil
}

// GetTimezoneInfo returns information about a timezone including DST status
func (tp *TimezoneParser) GetTimezoneInfo(loc *time.Location, t time.Time) TimezoneInfo {
	name, offset := t.In(loc).Zone()

	// Check if DST is active by comparing with 6 months earlier/later
	_, winterOffset := t.AddDate(0, -6, 0).In(loc).Zone()
	_, summerOffset := t.AddDate(0, 6, 0).In(loc).Zone()

	isDST := offset != winterOffset || offset != summerOffset

	return TimezoneInfo{
		Name:      name,
		Offset:    offset,
		IsDST:     isDST,
		Location:  loc,
		Formatted: formatTimezoneOffset(offset),
	}
}

// TimezoneInfo contains information about a timezone at a specific time
type TimezoneInfo struct {
	Name      string         // Short name (e.g., "IST", "PST", "PDT")
	Offset    int            // Offset from UTC in seconds
	IsDST     bool           // Whether DST is active
	Location  *time.Location // The location object
	Formatted string         // Formatted offset (e.g., "+05:30", "-08:00")
}

// formatTimezoneOffset formats a timezone offset in seconds to ±HH:MM format
func formatTimezoneOffset(offsetSeconds int) string {
	sign := "+"
	if offsetSeconds < 0 {
		sign = "-"
		offsetSeconds = -offsetSeconds
	}

	hours := offsetSeconds / 3600
	minutes := (offsetSeconds % 3600) / 60

	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

// ValidateTimezoneForLocation validates that a timezone offset is plausible for a location.
// It uses the request date so DST-sensitive locations are checked consistently.
func (tp *TimezoneParser) ValidateTimezoneForLocation(loc *time.Location, latitude, longitude float64, date time.Time) (bool, string) {
	if date.IsZero() {
		return false, "timezone validation date is required"
	}

	// Longitude roughly corresponds to UTC offset: 15 degrees per hour
	expectedOffset := int(longitude / 15.0 * 3600)

	testTime := date.In(loc)
	_, actualOffset := testTime.Zone()

	// Allow for some variance (±3 hours or ±45 degrees longitude)
	variance := 3 * 3600
	if actualOffset < expectedOffset-variance || actualOffset > expectedOffset+variance {
		warning := fmt.Sprintf("timezone %s (offset %s) may not match location coordinates (%.2f, %.2f)",
			loc.String(), formatTimezoneOffset(actualOffset), latitude, longitude)
		return false, warning
	}

	return true, ""
}
