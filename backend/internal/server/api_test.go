package server

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Acauhi99/med-vault/internal/shared/httpx"
)

func TestLogLoginFailure(t *testing.T) {
	var buf bytes.Buffer
	api := &API{logger: slog.New(slog.NewJSONHandler(&buf, nil))}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req = req.WithContext(httpx.ContextWithRequestID(context.Background(), "req-123"))
	req.RemoteAddr = "192.0.2.10:12345"
	req.Header.Set("User-Agent", "test-agent")

	api.logLoginFailure(req, "test@example.com", context.Canceled)

	out := buf.String()
	for _, want := range []string{"auth.login_failed", "req-123", "test@example.com", "192.0.2.10", "test-agent"} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output missing %q: %s", want, out)
		}
	}
}
