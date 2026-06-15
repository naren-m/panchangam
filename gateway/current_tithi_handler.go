package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
)

type currentTithiResponse struct {
	Date          string                    `json:"date"`
	Tithi         currentTithiDetails       `json:"tithi"`
	PanchaAnga    currentPanchaAngaSummary  `json:"pancha_anga"`
	Raasi         string                    `json:"raasi,omitempty"`
	Day           currentDaySummary         `json:"day"`
	Calculation   currentCalculationSummary `json:"calculation"`
	GeneratedAt   string                    `json:"generated_at"`
	NextRefreshAt string                    `json:"next_refresh_at"`
}

type currentTithiDetails struct {
	Name            string `json:"name"`
	TraditionalName string `json:"traditional_name"`
	Number          int    `json:"number"`
	Paksha          string `json:"paksha"`
	PakshaDay       int    `json:"paksha_day"`
	Type            string `json:"type"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

type currentPanchaAngaSummary struct {
	Nakshatra string `json:"nakshatra"`
	Yoga      string `json:"yoga"`
	Karana    string `json:"karana"`
	Vara      string `json:"vara"`
}

type currentDaySummary struct {
	SunriseTime    string            `json:"sunrise_time"`
	SunsetTime     string            `json:"sunset_time"`
	AbhijitMuhurta currentTimeWindow `json:"abhijit_muhurta"`
}

type currentTimeWindow struct {
	Name       string `json:"name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Auspicious bool   `json:"auspicious"`
}

type currentCalculationSummary struct {
	Timezone       string `json:"timezone"`
	Region         string `json:"region"`
	CalendarSystem string `json:"calendar_system"`
	Method         string `json:"method"`
	Locale         string `json:"locale"`
}

// handleCurrentTithi returns the compact tithi summary expected by the phone and watch app.
func (g *GatewayServer) handleCurrentTithi(client ppb.PanchangamClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		lat, latValue, ok, err := parseQueryFloat(query, "latitude", "lat")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: latitude", nil)
			return
		}
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid latitude format", map[string]interface{}{
				"parameter": "latitude",
				"value":     latValue,
				"expected":  "float64",
			})
			return
		}

		lng, lngValue, ok, err := parseQueryFloat(query, "longitude", "lng")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: longitude", nil)
			return
		}
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid longitude format", map[string]interface{}{
				"parameter": "longitude",
				"value":     lngValue,
				"expected":  "float64",
			})
			return
		}

		timezone := queryValue(query, "timezone", "tz")
		if timezone == "" {
			timezone = "UTC"
		}
		date, err := normalizeCurrentTithiDate(queryValue(query, "date"), timezone)
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid date format", map[string]interface{}{
				"parameter": "date",
				"expected":  "YYYY-MM-DD or RFC3339",
			})
			return
		}

		region := queryValue(query, "region")
		if region == "" {
			region = "global"
		}
		method := queryValue(query, "method", "calculation_method")
		if method == "" {
			method = "traditional"
		}
		locale := queryValue(query, "locale")
		if locale == "" {
			locale = "en"
		}
		calendarSystem := queryValue(query, "calendar_system", "calendarSystem")
		if calendarSystem == "" {
			calendarSystem = "Purnimanta"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		resp, err := client.Get(ctx, &ppb.GetPanchangamRequest{
			Date:              date,
			Latitude:          lat,
			Longitude:         lng,
			Timezone:          timezone,
			Region:            region,
			CalculationMethod: method,
			Locale:            locale,
		})
		if err != nil {
			handleGRPCError(w, r, err)
			return
		}
		if resp == nil || resp.PanchangamData == nil {
			writeErrorResponse(w, r, http.StatusInternalServerError, "EMPTY_RESPONSE", "Panchangam service returned no data", nil)
			return
		}

		generatedAt := time.Now()
		referenceAt := currentTithiReferenceTime(date, timezone, generatedAt)
		summary := makeCurrentTithiResponse(resp.PanchangamData, timezone, region, calendarSystem, method, locale, generatedAt, referenceAt)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		if err := json.NewEncoder(w).Encode(summary); err != nil {
			logger.Error("Failed to encode current tithi response", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR", "Failed to encode response", nil)
			return
		}
	}
}

func makeCurrentTithiResponse(data *ppb.PanchangamData, timezone, region, calendarSystem, method, locale string, generatedAt time.Time, referenceAt time.Time) currentTithiResponse {
	startTime, endTime, hasTithiWindow := currentTithiWindow(data, timezone, referenceAt)
	nextRefreshAt := generatedAt.Add(1 * time.Hour)
	if hasTithiWindow && endTime.After(referenceAt) {
		nextRefreshAt = endTime
	}

	return currentTithiResponse{
		Date:  data.Date,
		Tithi: makeCurrentTithiDetails(data.Tithi, startTime, endTime),
		PanchaAnga: currentPanchaAngaSummary{
			Nakshatra: data.Nakshatra,
			Yoga:      data.Yoga,
			Karana:    data.Karana,
			Vara:      currentTithiEventValue(data.Events, "VARA", "Vara"),
		},
		Raasi: currentTithiEventValue(data.Events, "RAASI", "Raasi"),
		Day: currentDaySummary{
			SunriseTime:    data.SunriseTime,
			SunsetTime:     data.SunsetTime,
			AbhijitMuhurta: makeCurrentAbhijitWindow(data.Events),
		},
		Calculation: currentCalculationSummary{
			Timezone:       timezone,
			Region:         region,
			CalendarSystem: calendarSystem,
			Method:         method,
			Locale:         locale,
		},
		GeneratedAt:   generatedAt.UTC().Format(time.RFC3339),
		NextRefreshAt: nextRefreshAt.UTC().Format(time.RFC3339),
	}
}

func currentTithiWindow(data *ppb.PanchangamData, timezone string, referenceAt time.Time) (time.Time, time.Time, bool) {
	startTime, hasStart := currentTithiEventDateTime(data.Date, timezone, data.Events, "TITHI_START")
	endTime, hasEnd := currentTithiEventDateTime(data.Date, timezone, data.Events, "TITHI_END")
	if hasStart && hasEnd {
		if !endTime.After(startTime) || endTime.Sub(startTime) < 12*time.Hour {
			endTime = endTime.Add(24 * time.Hour)
		}
		referenceTime := currentTithiLocalTime(referenceAt, timezone)
		for !referenceTime.Before(endTime) {
			startTime = startTime.Add(24 * time.Hour)
			endTime = endTime.Add(24 * time.Hour)
		}
		for referenceTime.Before(startTime) {
			startTime = startTime.Add(-24 * time.Hour)
			endTime = endTime.Add(-24 * time.Hour)
		}
		return startTime, endTime, true
	}

	return referenceAt.Add(-1 * time.Hour), referenceAt.Add(24 * time.Hour), false
}

func currentTithiReferenceTime(date, timezone string, fallback time.Time) time.Time {
	if parsed, err := time.Parse(time.RFC3339, date); err == nil {
		return parsed
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	if parsed, err := time.ParseInLocation("2006-01-02", date, loc); err == nil {
		localFallback := fallback.In(loc)
		return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), localFallback.Hour(), localFallback.Minute(), localFallback.Second(), localFallback.Nanosecond(), loc)
	}
	return fallback
}

func currentTithiLocalTime(value time.Time, timezone string) time.Time {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return value.UTC()
	}
	return value.In(loc)
}

func currentTithiEventDateTime(date, timezone string, events []*ppb.PanchangamEvent, eventType string) (time.Time, bool) {
	clock := currentTithiEventClock(events, eventType)
	if clock == "" {
		return time.Time{}, false
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+clock, loc)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func currentTithiEventClock(events []*ppb.PanchangamEvent, eventType string) string {
	for _, event := range events {
		if event != nil && event.EventType == eventType && strings.TrimSpace(event.Time) != "" {
			return strings.TrimSpace(event.Time)
		}
	}
	return ""
}

func makeCurrentTithiDetails(display string, startTime, endTime time.Time) currentTithiDetails {
	name := strings.TrimSpace(display)
	if idx := strings.Index(name, " - "); idx >= 0 {
		name = strings.TrimSpace(name[:idx])
	}
	if name == "" {
		name = "Tithi"
	}

	paksha := "Shukla"
	if strings.Contains(display, "Krishna") {
		paksha = "Krishna"
	}

	pakshaDay := currentTithiPakshaDay(display)
	number := pakshaDay
	if paksha == "Krishna" {
		number = 15 + pakshaDay
	}

	return currentTithiDetails{
		Name:            name,
		TraditionalName: name,
		Number:          number,
		Paksha:          paksha,
		PakshaDay:       pakshaDay,
		Type:            currentTithiType(pakshaDay),
		StartTime:       startTime.UTC().Format(time.RFC3339),
		EndTime:         endTime.UTC().Format(time.RFC3339),
	}
}

func makeCurrentAbhijitWindow(events []*ppb.PanchangamEvent) currentTimeWindow {
	start := currentTithiEventValue(events, "ABHIJIT_MUHURTA", "Abhijit Muhurta")
	if start == "" {
		return currentTimeWindow{
			Name:       "Abhijit Muhurta",
			StartTime:  "12:00:00",
			EndTime:    "12:48:00",
			Auspicious: false,
		}
	}

	end := currentTithiAddClockMinutes(start, 48)
	return currentTimeWindow{
		Name:       "Abhijit Muhurta",
		StartTime:  start,
		EndTime:    end,
		Auspicious: true,
	}
}

func currentTithiEventValue(events []*ppb.PanchangamEvent, eventType, prefix string) string {
	for _, event := range events {
		if event == nil || event.EventType != eventType {
			continue
		}
		if event.Time != "" && eventType == "ABHIJIT_MUHURTA" {
			return event.Time
		}
		name := strings.TrimSpace(event.Name)
		name = strings.TrimPrefix(name, prefix+":")
		return strings.TrimSpace(name)
	}
	return ""
}

func currentTithiAddClockMinutes(value string, minutes int) string {
	clock, err := time.Parse("15:04:05", value)
	if err != nil {
		return value
	}
	return clock.Add(time.Duration(minutes) * time.Minute).Format("15:04:05")
}

func currentTithiPakshaDay(display string) int {
	idx := strings.Index(display, "Day ")
	if idx < 0 {
		return 1
	}
	remaining := display[idx+len("Day "):]
	fields := strings.Fields(remaining)
	if len(fields) == 0 {
		return 1
	}
	day, err := strconv.Atoi(fields[0])
	if err != nil || day < 1 || day > 15 {
		return 1
	}
	return day
}

func currentTithiType(pakshaDay int) string {
	switch pakshaDay {
	case 1, 6, 11:
		return "Nanda"
	case 2, 7, 12:
		return "Bhadra"
	case 3, 8, 13:
		return "Jaya"
	case 4, 9, 14:
		return "Rikta"
	case 5, 10, 15:
		return "Purna"
	default:
		return "Nanda"
	}
}

func parseQueryFloat(query map[string][]string, names ...string) (float64, string, bool, error) {
	value := queryValue(query, names...)
	if value == "" {
		return 0, "", false, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	return parsed, value, true, err
}

func queryValue(query map[string][]string, names ...string) string {
	for _, name := range names {
		values := query[name]
		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}

func normalizeCurrentTithiDate(value, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	if value == "" {
		return time.Now().In(loc).Format(time.RFC3339), nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed.In(loc).Format(time.RFC3339), nil
	}

	if _, err := time.Parse("2006-01-02", value); err == nil {
		return value, nil
	}

	return "", fmt.Errorf("invalid current tithi date: %w", err)
}
