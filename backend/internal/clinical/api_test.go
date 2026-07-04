package clinical

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type clinicalAPIResponse struct {
	Data struct {
		ID        string `json:"id"`
		PatientID string `json:"patient_id"`
		Status    string `json:"status"`
	} `json:"data"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error"`
}

type mockClinicalRepo struct {
	cases map[uuid.UUID]*domain.Case
}

func newMockClinicalRepo() *mockClinicalRepo {
	return &mockClinicalRepo{cases: make(map[uuid.UUID]*domain.Case)}
}

func (m *mockClinicalRepo) Create(_ context.Context, c *domain.Case) error {
	m.cases[c.ID] = c
	return nil
}

func (m *mockClinicalRepo) GetByID(_ context.Context, tenantID, caseID uuid.UUID) (*domain.Case, error) {
	c, ok := m.cases[caseID]
	if !ok || c.TenantID != tenantID {
		return nil, errors.New("case not found")
	}
	return c, nil
}

func (m *mockClinicalRepo) ListByPatient(_ context.Context, tenantID, patientID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && c.PatientID == patientID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockClinicalRepo) ListByDoctor(_ context.Context, tenantID, doctorID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && c.DoctorID != nil && *c.DoctorID == doctorID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockClinicalRepo) ListByTenant(_ context.Context, tenantID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockClinicalRepo) Update(_ context.Context, c *domain.Case) error {
	m.cases[c.ID] = c
	return nil
}

func (m *mockClinicalRepo) AddSymptom(_ context.Context, _, _ uuid.UUID, _ *domain.Symptom) error {
	return nil
}

func (m *mockClinicalRepo) WriteDiagnosis(_ context.Context, _, _ uuid.UUID, _ *domain.Diagnosis) error {
	return nil
}

type mockClinicalTenantRepo struct{}

func (m *mockClinicalTenantRepo) FindUserTenants(uuid.UUID) ([]authdomain.UserTenant, error) {
	return nil, nil
}
func (m *mockClinicalTenantRepo) FindUserTenant(uuid.UUID, uuid.UUID) (*authdomain.UserTenant, error) {
	return nil, errors.New("not found")
}
func (m *mockClinicalTenantRepo) AddMember(context.Context, uuid.UUID, uuid.UUID, string) error {
	return nil
}
func (m *mockClinicalTenantRepo) RemoveMember(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (m *mockClinicalTenantRepo) ListMembers(context.Context, uuid.UUID) ([]authdomain.UserTenant, error) {
	return nil, nil
}
func (m *mockClinicalTenantRepo) Reactivate(context.Context, uuid.UUID) (*authdomain.Tenant, error) {
	return nil, nil
}

func (m *mockClinicalTenantRepo) Create(context.Context, string) (*authdomain.Tenant, error) {
	return nil, nil
}

func (m *mockClinicalTenantRepo) Suspend(context.Context, uuid.UUID) (*authdomain.Tenant, error) {
	return nil, nil
}

func TestAPI_CreateCase(t *testing.T) {
	repo := newMockClinicalRepo()
	api := NewAPI(repo, &mockClinicalTenantRepo{})

	body, err := json.Marshal(generated.CreateCaseRequest{
		Symptoms: []generated.AddSymptomRequest{{Description: "headache", Severity: generated.AddSymptomRequestSeverity("low")}},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	principal := sharedauth.Principal{UserID: uuid.New(), TenantID: uuid.New(), Role: sharedauth.RolePatient}

	tests := []struct {
		name        string
		principal   *sharedauth.Principal
		wantStatus  int
		wantCode    string
		wantStatusS string
	}{
		{name: "unauthenticated", wantStatus: http.StatusUnauthorized, wantCode: "UNAUTHORIZED"},
		{name: "doctor forbidden", principal: &sharedauth.Principal{UserID: uuid.New(), TenantID: principal.TenantID, Role: sharedauth.RoleDoctor}, wantStatus: http.StatusForbidden, wantCode: "FORBIDDEN"},
		{name: "patient creates case", principal: &principal, wantStatus: http.StatusCreated, wantStatusS: string(domain.CaseStatusOpen)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader(body))
			if tc.principal != nil {
				req = req.WithContext(sharedauth.ContextWithPrincipal(req.Context(), *tc.principal))
			}

			api.CreateCase(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			var resp clinicalAPIResponse
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
			if resp.Data.PatientID != principal.UserID.String() {
				t.Fatalf("patient_id = %q, want %q", resp.Data.PatientID, principal.UserID.String())
			}
			if resp.Data.Status != tc.wantStatusS {
				t.Fatalf("status = %q, want %q", resp.Data.Status, tc.wantStatusS)
			}
		})
	}
}

var _ domain.Repository = (*mockClinicalRepo)(nil)
var _ authdomain.TenantRepository = (*mockClinicalTenantRepo)(nil)
