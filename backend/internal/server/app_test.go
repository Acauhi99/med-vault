package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Acauhi99/med-vault/internal/shared/config"
)

func TestHealthReturnsOK(t *testing.T) {
	app := newTestApp(t)
	defer app.Close()

	recorder := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != "ok" {
		t.Fatalf("expected health ok, got %q", body.Data.Status)
	}
	if body.Error != nil {
		t.Fatalf("expected nil error, got %#v", body.Error)
	}
}

func TestOpenAPIEndpointsAreMountedUnderAPIV1(t *testing.T) {
	app := newTestApp(t)
	defer app.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	request.Header.Set("X-Request-Id", "test-request")
	app.Server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
	if recorder.Header().Get("X-Request-Id") != "test-request" {
		t.Fatalf("expected request id header to be propagated")
	}
}

func newTestApp(t *testing.T) *App {
	t.Helper()

	app, err := New(context.Background(), config.Config{
		Env:             "test",
		HTTPAddr:        ":0",
		RequestIDHeader: "X-Request-Id",
	}, slog.Default())
	if err != nil {
		t.Fatalf("new app: %v", err)
	}
	return app
}
