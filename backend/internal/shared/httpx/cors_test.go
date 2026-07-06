package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	called := false
	handler := CORSMiddleware("https://med-vault.space,https://www.med-vault.space")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req.Header.Set("Origin", "https://med-vault.space")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.True(t, called)
	assert.Equal(t, "https://med-vault.space", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORSMiddlewareRejectsDisallowedPreflight(t *testing.T) {
	handler := CORSMiddleware("https://med-vault.space")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
