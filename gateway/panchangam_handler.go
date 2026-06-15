package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	ppb "github.com/naren-m/panchangam/proto"
)

// handlePanchangam handles HTTP requests to the panchangam endpoint with caching.
func (g *GatewayServer) handlePanchangam(client ppb.PanchangamClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		date := query.Get("date")
		if date == "" {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: date", nil)
			return
		}

		lat, latValue, ok, err := parseQueryFloat(query, "lat")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: lat", nil)
			return
		}
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid latitude format", map[string]interface{}{
				"parameter": "lat",
				"value":     latValue,
				"expected":  "float64",
			})
			return
		}

		lng, lngValue, ok, err := parseQueryFloat(query, "lng")
		if !ok {
			writeErrorResponse(w, r, http.StatusBadRequest, "MISSING_PARAMETER", "Missing required parameter: lng", nil)
			return
		}
		if err != nil {
			writeErrorResponse(w, r, http.StatusBadRequest, "INVALID_PARAMETER", "Invalid longitude format", map[string]interface{}{
				"parameter": "lng",
				"value":     lngValue,
				"expected":  "float64",
			})
			return
		}

		timezone := query.Get("tz")
		if timezone == "" {
			timezone = "UTC"
		}

		region := query.Get("region")
		if region == "" {
			region = "global"
		}

		calculationMethod := query.Get("method")
		if calculationMethod == "" {
			calculationMethod = "traditional"
		}

		locale := query.Get("locale")
		if locale == "" {
			locale = "en"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		var responseData interface{}
		cacheHit := false

		if g.cache != nil {
			cacheKey := g.cache.GenerateCacheKey(date, region, calculationMethod, lat, lng)

			cachedData, err := g.cache.Get(ctx, cacheKey)
			if err != nil {
				logger.Error("Cache get error", "error", err, "key", cacheKey)
			} else if cachedData != nil {
				responseData = convertCacheToResponse(cachedData)
				cacheHit = true
				logger.Debug("Cache hit", "key", cacheKey, "date", date)
			}
		}

		if !cacheHit {
			req := &ppb.GetPanchangamRequest{
				Date:              date,
				Latitude:          lat,
				Longitude:         lng,
				Timezone:          timezone,
				Region:            region,
				CalculationMethod: calculationMethod,
				Locale:            locale,
			}

			resp, err := client.Get(ctx, req)
			if err != nil {
				handleGRPCError(w, r, err)
				return
			}

			responseData = resp.PanchangamData

			if g.cache != nil && resp.PanchangamData != nil {
				cacheData := convertResponseToCache(resp.PanchangamData)
				cacheKey := g.cache.GenerateCacheKey(date, region, calculationMethod, lat, lng)

				if err := g.cache.Set(ctx, cacheKey, cacheData); err != nil {
					logger.Error("Cache set error", "error", err, "key", cacheKey)
				} else {
					logger.Debug("Cache set", "key", cacheKey, "date", date)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if cacheHit {
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Cache-Control", "public, max-age=1800")
		} else {
			w.Header().Set("X-Cache", "MISS")
			w.Header().Set("Cache-Control", "public, max-age=300")
		}

		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			logger.Error("Failed to encode response", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR", "Failed to encode response", nil)
			return
		}
	}
}
