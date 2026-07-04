package auditlog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Acauhi99/med-vault/internal/audit/application"
	"github.com/Acauhi99/med-vault/internal/audit/domain"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type stubRepo struct{}

func (s *stubRepo) Create(_ context.Context, _ *domain.AuditLog) error { return nil }
func (s *stubRepo) ListByTenant(_ context.Context, _ uuid.UUID, _, _ int, _ string, _ *uuid.UUID, _ string, _ *uuid.UUID) ([]domain.AuditLog, int, error) {
	return nil, 0, nil
}

func TestLogWithPrincipal(t *testing.T) {
	logger := NewLogger(application.NewLogActionCommand(&stubRepo{}))
	principal := auth.Principal{UserID: uuid.New(), TenantID: uuid.New(), Role: auth.RolePatient}
	ctx := auth.ContextWithPrincipal(context.Background(), principal)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	logger.Log(ctx, req, "user.registered", "user", uuid.New(), nil)
}

func TestLogWithoutPrincipal(t *testing.T) {
	logger := NewLogger(application.NewLogActionCommand(&stubRepo{}))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	logger.Log(context.Background(), req, "user.registered", "user", uuid.New(), nil)
}
