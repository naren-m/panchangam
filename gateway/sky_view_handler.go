package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/naren-m/panchangam/services/skyview"
)

// handleSkyView handles requests for sky visualization data.
func (g *GatewayServer) handleSkyView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if g.ephemerisProvider == nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "EPHEMERIS_UNAVAILABLE",
				"Ephemeris data provider is not available", nil)
			return
		}

		query := r.URL.Query()

		dateStr := queryValue(query, "date")
		timeStr := queryValue(query, "time")
		observationTime, err := parseSkyViewObservationTime(dateStr, timeStr, time.Now().UTC())
		if err != nil {
			if timeStr != "" {
				writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_DATETIME",
					"Invalid date/time format. Expected: date=YYYY-MM-DD&time=HH:MM:SS", map[string]interface{}{
						"date": dateStr,
						"time": timeStr,
					})
				return
			}
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_DATE",
				"Invalid date format. Expected: YYYY-MM-DD", map[string]interface{}{
					"date": dateStr,
				})
			return
		}

		lat, latValue, ok, err := parseQueryFloat(query, "lat")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER",
				"Missing required parameter: lat", nil)
			return
		}

		if err != nil || lat < -90 || lat > 90 {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER",
				"Invalid latitude. Must be between -90 and 90", map[string]interface{}{
					"parameter": "lat",
					"value":     latValue,
				})
			return
		}

		lng, lngValue, ok, err := parseQueryFloat(query, "lng")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER",
				"Missing required parameter: lng", nil)
			return
		}
		if err != nil || lng < -180 || lng > 180 {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER",
				"Invalid longitude. Must be between -180 and 180", map[string]interface{}{
					"parameter": "lng",
					"value":     lngValue,
				})
			return
		}

		altitude := 0.0
		var altitudeValue string
		if altitude, altitudeValue, ok, err = parseQueryFloat(query, "alt"); ok {
			if err != nil {
				writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER",
					"Invalid altitude format", map[string]interface{}{
						"parameter": "alt",
						"value":     altitudeValue,
					})
				return
			}
		}

		timezone := queryValue(query, "tz")
		if timezone == "" {
			timezone = "UTC"
		}

		observer := skyview.Observer{
			Latitude:  lat,
			Longitude: lng,
			Altitude:  altitude,
			Timezone:  timezone,
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		service := skyview.NewSkyViewService(g.ephemerisProvider)
		skyViewData, err := service.GetSkyView(ctx, observer, observationTime)
		if err != nil {
			logger.Error("Failed to get sky view", "error", err,
				"lat", lat, "lng", lng, "time", observationTime)
			writeErrorResponse(w, r, http.StatusInternalServerError, "SKYVIEW_ERROR",
				"Failed to calculate sky view data", map[string]interface{}{
					"error": err.Error(),
				})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")

		if err := json.NewEncoder(w).Encode(skyViewData); err != nil {
			logger.Error("Failed to encode sky view response", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR",
				"Failed to encode response", nil)
			return
		}

		logger.Debug("Sky view request completed",
			"lat", lat, "lng", lng, "time", observationTime,
			"visible_bodies", len(skyViewData.VisibleBodies))
	}
}

func parseSkyViewObservationTime(dateStr, timeStr string, fallback time.Time) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	timeStr = strings.TrimSpace(timeStr)

	if dateStr == "" {
		return fallback, nil
	}
	if timeStr != "" {
		return time.Parse("2006-01-02T15:04:05", dateStr+"T"+timeStr)
	}
	return time.Parse("2006-01-02", dateStr)
}
