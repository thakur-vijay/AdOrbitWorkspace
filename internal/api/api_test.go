package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeathReturnsJSONStatus(t *testing.T) {
	response := httptest.NewRecorder()
	health(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("health response is not JSON: %v", err)
	}

	if body["status"] != "ok" || len(body) != 1 {
		t.Fatalf("health response = %#v, want {status: ok}", body)
	}
}

// func TestRouterWiresVersionAuthEndpointsAndIdempotency(t *testing.T) {
// 	router := NewRouter(Dependencies{
// 		IdempotencyMiddleware: idempotency.NewMiddleware(),
// 	})

// }
