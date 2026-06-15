package gateway

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	handler := addHealthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach next handler for health check")
	}))

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var health map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &health)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", health["status"])
	assert.Equal(t, "panchangam-gateway", health["service"])
	assert.NotEmpty(t, health["timestamp"])
}

func TestWriteHealthCheckResponseReturnsWriteError(t *testing.T) {
	err := writeHealthCheckResponse(errorResponseWriter{}, time.Unix(0, 0).UTC())

	if err == nil {
		t.Fatal("expected write error")
	}
	assert.Contains(t, err.Error(), "write failed")
}

func TestLoggingMiddleware(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := loggingMiddleware(next)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-Id", "test-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, nextCalled)
	assert.Equal(t, "test-123", w.Header().Get("X-Request-Id"))
	assert.NotEmpty(t, w.Header().Get("X-Response-Time"))
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	if id1 == id2 {
		t.Logf("Generated identical IDs (rare but possible): %s", id1)
	}
	assert.Contains(t, id1, "req_")
	assert.Contains(t, id2, "req_")
}

func TestRequestIDFromRequest(t *testing.T) {
	withHeader := httptest.NewRequest(http.MethodGet, "/test", nil)
	withHeader.Header.Set("X-Request-Id", "caller-request")

	assert.Equal(t, "caller-request", requestIDFromRequest(withHeader))

	withoutHeader := httptest.NewRequest(http.MethodGet, "/test", nil)
	generated := requestIDFromRequest(withoutHeader)

	assert.NotEmpty(t, generated)
	assert.Contains(t, generated, "req_")
}

func TestParseCORSOrigins(t *testing.T) {
	defaults := []string{
		"http://localhost:5173",
		"http://localhost:3000",
	}

	tests := []struct {
		name string
		env  string
		want []string
	}{
		{
			name: "empty env uses defaults",
			env:  "",
			want: defaults,
		},
		{
			name: "trims comma separated origins",
			env:  " https://app.example.com ,https://admin.example.com ",
			want: []string{"https://app.example.com", "https://admin.example.com"},
		},
		{
			name: "blank entries are ignored",
			env:  "https://app.example.com,,  ,https://admin.example.com",
			want: []string{"https://app.example.com", "https://admin.example.com"},
		},
		{
			name: "all blank entries use defaults",
			env:  " , , ",
			want: defaults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseCORSOrigins(tt.env, defaults))
		})
	}
}

type errorResponseWriter struct{}

func (errorResponseWriter) Header() http.Header {
	return http.Header{}
}

func (errorResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func (errorResponseWriter) WriteHeader(int) {}
