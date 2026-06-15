package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/naren-m/panchangam/cache"
)

// convertResponseToCache converts gRPC response to cache format.
func convertResponseToCache(data interface{}) *cache.PanchangamCacheData {
	responseMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	cacheData := &cache.PanchangamCacheData{}

	if date, ok := responseMap["date"].(string); ok {
		cacheData.Date = date
	}
	if tithi, ok := responseMap["tithi"].(string); ok {
		cacheData.Tithi = tithi
	}
	if nakshatra, ok := responseMap["nakshatra"].(string); ok {
		cacheData.Nakshatra = nakshatra
	}
	if yoga, ok := responseMap["yoga"].(string); ok {
		cacheData.Yoga = yoga
	}
	if karana, ok := responseMap["karana"].(string); ok {
		cacheData.Karana = karana
	}
	if sunriseTime, ok := responseMap["sunrise_time"].(string); ok {
		cacheData.SunriseTime = sunriseTime
	}
	if sunsetTime, ok := responseMap["sunset_time"].(string); ok {
		cacheData.SunsetTime = sunsetTime
	}

	if events, ok := responseMap["events"].([]interface{}); ok {
		cacheData.Events = make([]cache.Event, len(events))
		for i, event := range events {
			if eventMap, ok := event.(map[string]interface{}); ok {
				cacheEvent := cache.Event{}
				if name, ok := eventMap["name"].(string); ok {
					cacheEvent.Name = name
				}
				if time, ok := eventMap["time"].(string); ok {
					cacheEvent.Time = time
				}
				if eventType, ok := eventMap["event_type"].(string); ok {
					cacheEvent.EventType = eventType
				}
				cacheData.Events[i] = cacheEvent
			}
		}
	}

	return cacheData
}

// convertCacheToResponse converts cache format to response format.
func convertCacheToResponse(data *cache.PanchangamCacheData) interface{} {
	events := make([]map[string]interface{}, len(data.Events))
	for i, event := range data.Events {
		events[i] = map[string]interface{}{
			"name":       event.Name,
			"time":       event.Time,
			"event_type": event.EventType,
		}
	}

	return map[string]interface{}{
		"date":         data.Date,
		"tithi":        data.Tithi,
		"nakshatra":    data.Nakshatra,
		"yoga":         data.Yoga,
		"karana":       data.Karana,
		"sunrise_time": data.SunriseTime,
		"sunset_time":  data.SunsetTime,
		"events":       events,
	}
}

// handleCacheHealth handles cache health check requests.
func (g *GatewayServer) handleCacheHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if g.cache == nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_DISABLED", "Cache is not enabled", nil)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := g.cache.HealthCheck(ctx); err != nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_UNHEALTHY", "Cache health check failed", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if err := writeCacheHealthResponse(w, time.Now().UTC()); err != nil {
			logger.Error("Failed to encode cache health response", "error", err)
		}
	}
}

type cacheHealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
}

func writeCacheHealthResponse(w http.ResponseWriter, now time.Time) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(cacheHealthResponse{
		Status:    "healthy",
		Timestamp: now.Format(time.RFC3339),
		Service:   "redis-cache",
	})
}

// handleCacheStats handles cache statistics requests.
func (g *GatewayServer) handleCacheStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if g.cache == nil {
			writeErrorResponse(w, r, http.StatusServiceUnavailable, "CACHE_DISABLED", "Cache is not enabled", nil)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		stats, err := g.cache.GetStats(ctx)
		if err != nil {
			writeErrorResponse(w, r, http.StatusInternalServerError, "STATS_ERROR", "Failed to get cache statistics", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			logger.Error("Failed to encode cache stats", "error", err)
			writeErrorResponse(w, r, http.StatusInternalServerError, "ENCODING_ERROR", "Failed to encode stats", nil)
		}
	}
}
