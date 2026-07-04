package application

import (
	"context"
	"testing"

	"github.com/Acauhi99/med-vault/internal/audit/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type mockAuditRepo struct {
	logs []domain.AuditLog
}

func newMockRepo() *mockAuditRepo {
	return &mockAuditRepo{}
}

func (m *mockAuditRepo) Create(_ context.Context, log *domain.AuditLog) error {
	m.logs = append(m.logs, *log)
	return nil
}

func (m *mockAuditRepo) ListByTenant(_ context.Context, tenantID uuid.UUID, offset, limit int, resourceType string, resourceID *uuid.UUID) ([]domain.AuditLog, int, error) {
	var filtered []domain.AuditLog
	for _, l := range m.logs {
		if l.TenantID == tenantID {
			filtered = append(filtered, l)
		}
	}
	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func TestLogActionSuccess(t *testing.T) {
	repo := newMockRepo()
	cmd := NewLogActionCommand(repo)

	input := LogActionInput{
		TenantID:     uuid.New(),
		UserID:       uuid.New(),
		Action:       "patient.created",
		ResourceType: "patient",
		ResourceID:   uuid.New(),
		Details:      map[string]any{"name": "John"},
		IPAddress:    "127.0.0.1",
		UserAgent:    "test",
	}

	log, err := cmd.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.Action != "patient.created" {
		t.Errorf("Action = %v, want patient.created", log.Action)
	}
	if log.TenantID != input.TenantID {
		t.Errorf("TenantID = %v, want %v", log.TenantID, input.TenantID)
	}
	if len(repo.logs) != 1 {
		t.Errorf("repo.logs = %d, want 1", len(repo.logs))
	}
}

func TestListAuditLogsAdminSuccess(t *testing.T) {
	repo := newMockRepo()
	tenantID := uuid.New()
	repo.logs = append(repo.logs,
		domain.AuditLog{ID: uuid.New(), TenantID: tenantID, Action: "a1"},
		domain.AuditLog{ID: uuid.New(), TenantID: tenantID, Action: "a2"},
	)

	q := NewListAuditLogsQuery(repo)
	p := sharedauth.Principal{UserID: uuid.New(), TenantID: tenantID, Role: sharedauth.RoleAdministrator}

	logs, total, err := q.Execute(context.Background(), p, 1, 20, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(logs) != 2 {
		t.Errorf("len = %d, want 2", len(logs))
	}
}

func TestListAuditLogsNonAdminFails(t *testing.T) {
	repo := newMockRepo()
	q := NewListAuditLogsQuery(repo)
	p := sharedauth.Principal{UserID: uuid.New(), TenantID: uuid.New(), Role: sharedauth.RolePatient}

	_, _, err := q.Execute(context.Background(), p, 1, 20, "", nil)
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}
