package api

import (
	"net/http/httptest"
	"testing"
)

func TestHeathReturnsJSONStatus(t *testing.T) {
	response := httptest.NewRecorder()
	health(response)
}
