package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Acauhi99/med-vault/internal/audit/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type auditAPIResponse struct {
	Data []struct {
		Action       string `json:"action"`
		ResourceType string `json:"resource_type"`
		TenantID     string `json:"tenant_id"`
	} `json:"data"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error"`
	Meta struct {
		Page     *int `json:"page"`
		PageSize *int `json:"page_size"`
		Total    *int `json:"total"`
	} `json:"meta"`
}

type mockAuditRepo struct {
	logs []domain.AuditLog
}

func (m *mockAuditRepo) Create(_ context.Context, log *domain.AuditLog) error {
	m.logs = append(m.logs, *log)
	return nil
}

func (m *mockAuditRepo) ListByTenant(_ context.Context, tenantID uuid.UUID, offset, limit int, action string, userID *uuid.UUID, resourceType string, resourceID *uuid.UUID) ([]domain.AuditLog, int, error) {
	var filtered []domain.AuditLog
	for _, l := range m.logs {
		if l.TenantID != tenantID {
			continue
		}
		if action != "" && l.Action != action {
			continue
		}
		if userID != nil && l.UserID != *userID {
			continue
		}
		if resourceType != "" && l.ResourceType != resourceType {
			continue
		}
		if resourceID != nil && l.ResourceID != *resourceID {
			continue
		}
		filtered = append(filtered, l)
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

func TestAPI_ListAuditLogs(t *testing.T) {
	tenantID := uuid.New()
	repo := &mockAuditRepo{logs: []domain.AuditLog{{ID: uuid.New(), TenantID: tenantID, UserID: uuid.New(), Action: "case.created", ResourceType: "case", ResourceID: uuid.New(), CreatedAt: time.Now().UTC()}}}
	api := NewAPI(repo)

	tests := []struct {
		name       string
		principal  *sharedauth.Principal
		wantStatus int
		wantCode   string
	}{
		{name: "patient forbidden", principal: &sharedauth.Principal{UserID: uuid.New(), TenantID: tenantID, Role: sharedauth.RolePatient}, wantStatus: http.StatusForbidden, wantCode: "FORBIDDEN"},
		{name: "admin gets logs", principal: &sharedauth.Principal{UserID: uuid.New(), TenantID: tenantID, Role: sharedauth.RoleAdministrator}, wantStatus: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/audit-logs", bytes.NewReader(nil))
			req = req.WithContext(sharedauth.ContextWithPrincipal(req.Context(), *tc.principal))

			page := generated.PageParam(1)
			pageSize := generated.PageSizeParam(20)
			api.ListAuditLogs(rec, req, generated.ListAuditLogsParams{Page: &page, PageSize: &pageSize})

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var resp auditAPIResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if tc.wantCode != "" {
				if resp.Error == nil || resp.Error.Code != tc.wantCode {
					t.Fatalf("error code = %#v, want %q", resp.Error, tc.wantCode)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %#v", resp.Error)
			}
			if len(resp.Data) != 1 || resp.Data[0].Action != "case.created" || resp.Data[0].ResourceType != "case" {
				t.Fatalf("unexpected data: %#v", resp.Data)
			}
			if resp.Meta.Page == nil || *resp.Meta.Page != 1 || resp.Meta.PageSize == nil || *resp.Meta.PageSize != 20 || resp.Meta.Total == nil || *resp.Meta.Total != 1 {
				t.Fatalf("unexpected meta: %#v", resp.Meta)
			}
		})
	}
}

var _ domain.Repository = (*mockAuditRepo)(nil)
