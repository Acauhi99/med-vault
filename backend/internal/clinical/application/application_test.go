package application

import (
	"context"
	"testing"
	"time"

	authdomain "github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type mockCaseRepo struct {
	cases map[uuid.UUID]*domain.Case
}

type mockTenantRepo struct {
	memberships map[uuid.UUID]map[uuid.UUID]authdomain.UserTenant
}

func newMockTenantRepo() *mockTenantRepo {
	return &mockTenantRepo{memberships: make(map[uuid.UUID]map[uuid.UUID]authdomain.UserTenant)}
}

func (m *mockTenantRepo) FindUserTenants(userID uuid.UUID) ([]authdomain.UserTenant, error) {
	var out []authdomain.UserTenant
	for _, tenants := range m.memberships {
		if ut, ok := tenants[userID]; ok {
			out = append(out, ut)
		}
	}
	return out, nil
}

func (m *mockTenantRepo) FindUserTenant(userID, tenantID uuid.UUID) (*authdomain.UserTenant, error) {
	if tenants, ok := m.memberships[tenantID]; ok {
		if ut, ok := tenants[userID]; ok {
			return &ut, nil
		}
	}
	return nil, ErrCaseNotFound
}

func (m *mockTenantRepo) AddMember(context.Context, uuid.UUID, uuid.UUID, string) error { return nil }

func (m *mockTenantRepo) RemoveMember(context.Context, uuid.UUID, uuid.UUID) error { return nil }

func (m *mockTenantRepo) ListMembers(context.Context, uuid.UUID) ([]authdomain.UserTenant, error) {
	return nil, nil
}

func (m *mockTenantRepo) Reactivate(context.Context, uuid.UUID) (*authdomain.Tenant, error) {
	return nil, nil
}

func (m *mockTenantRepo) Create(context.Context, string) (*authdomain.Tenant, error) { return nil, nil }

func (m *mockTenantRepo) Suspend(context.Context, uuid.UUID) (*authdomain.Tenant, error) {
	return nil, nil
}

func newMockRepo() *mockCaseRepo {
	return &mockCaseRepo{cases: make(map[uuid.UUID]*domain.Case)}
}

func (m *mockCaseRepo) Create(_ context.Context, c *domain.Case) error {
	m.cases[c.ID] = c
	return nil
}

func (m *mockCaseRepo) GetByID(_ context.Context, tenantID, caseID uuid.UUID) (*domain.Case, error) {
	c, ok := m.cases[caseID]
	if !ok || c.TenantID != tenantID {
		return nil, ErrCaseNotFound
	}
	return c, nil
}

func (m *mockCaseRepo) ListByPatient(_ context.Context, tenantID, patientID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && c.PatientID == patientID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockCaseRepo) ListByDoctor(_ context.Context, tenantID, doctorID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && c.DoctorID != nil && *c.DoctorID == doctorID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockCaseRepo) ListByTenant(_ context.Context, tenantID uuid.UUID, status string, _, _ int) ([]domain.Case, int, error) {
	var out []domain.Case
	for _, c := range m.cases {
		if c.TenantID == tenantID && (status == "" || string(c.Status) == status) {
			out = append(out, *c)
		}
	}
	return out, len(out), nil
}

func (m *mockCaseRepo) Update(_ context.Context, c *domain.Case) error {
	m.cases[c.ID] = c
	return nil
}

func (m *mockCaseRepo) AddSymptom(_ context.Context, _, _ uuid.UUID, s *domain.Symptom) error {
	_ = s
	return nil
}

func (m *mockCaseRepo) WriteDiagnosis(_ context.Context, _, _ uuid.UUID, d *domain.Diagnosis) error {
	_ = d
	return nil
}

func newPrincipal(userID, tenantID uuid.UUID, role sharedauth.Role) sharedauth.Principal {
	return sharedauth.Principal{UserID: userID, TenantID: tenantID, Role: role}
}

func openCase(repo *mockCaseRepo, id, tenantID, patientID uuid.UUID) {
	repo.cases[id] = &domain.Case{
		ID:        id,
		TenantID:  tenantID,
		PatientID: patientID,
		Status:    domain.CaseStatusOpen,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func assignedCase(repo *mockCaseRepo, id, tenantID, patientID, doctorID uuid.UUID) {
	repo.cases[id] = &domain.Case{
		ID:        id,
		TenantID:  tenantID,
		PatientID: patientID,
		DoctorID:  &doctorID,
		Status:    domain.CaseStatusAssigned,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func diagnosedCase(repo *mockCaseRepo, id, tenantID, patientID, doctorID uuid.UUID) {
	repo.cases[id] = &domain.Case{
		ID:        id,
		TenantID:  tenantID,
		PatientID: patientID,
		DoctorID:  &doctorID,
		Status:    domain.CaseStatusDiagnosed,
		Diagnosis: &domain.Diagnosis{ID: uuid.New(), DoctorID: doctorID, Notes: "test"},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func TestCreateCaseSuccess(t *testing.T) {
	repo := newMockRepo()
	cmd := NewCreateCaseCommand(repo)
	p := newPrincipal(uuid.New(), uuid.New(), sharedauth.RolePatient)

	cs, err := cmd.Execute(context.Background(), p, []CreateSymptomRequest{
		{Description: "headache", Severity: "low"},
		{Description: "fever", Severity: "medium"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.PatientID != p.UserID {
		t.Errorf("PatientID = %v, want %v", cs.PatientID, p.UserID)
	}
	if cs.Status != domain.CaseStatusOpen {
		t.Errorf("Status = %v, want open", cs.Status)
	}
	if len(cs.Symptoms) != 2 {
		t.Errorf("Symptoms = %d, want 2", len(cs.Symptoms))
	}
}

func TestCreateCaseDoctorRejected(t *testing.T) {
	repo := newMockRepo()
	cmd := NewCreateCaseCommand(repo)
	p := newPrincipal(uuid.New(), uuid.New(), sharedauth.RoleDoctor)

	_, err := cmd.Execute(context.Background(), p, []CreateSymptomRequest{
		{Description: "pain", Severity: "high"},
	})
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestCreateCaseEmptySymptoms(t *testing.T) {
	repo := newMockRepo()
	cmd := NewCreateCaseCommand(repo)
	p := newPrincipal(uuid.New(), uuid.New(), sharedauth.RolePatient)

	_, err := cmd.Execute(context.Background(), p, nil)
	if err != ErrNoSymptoms {
		t.Errorf("expected ErrNoSymptoms, got %v", err)
	}
}

func TestAddSymptomSuccess(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)

	cmd := NewAddSymptomCommand(repo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	s, err := cmd.Execute(context.Background(), p, caseID, "nausea", "low")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Description != "nausea" {
		t.Errorf("Description = %v, want nausea", s.Description)
	}
}

func TestAddSymptomDiagnosedCaseRejected(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	diagnosedCase(repo, caseID, tenantID, patientID, uuid.New())

	cmd := NewAddSymptomCommand(repo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, err := cmd.Execute(context.Background(), p, caseID, "cough", "low")
	if err != ErrCaseNotOpen {
		t.Errorf("expected ErrCaseNotOpen, got %v", err)
	}
}

func TestAssignDoctorSuccess(t *testing.T) {
	repo := newMockRepo()
	tenants := newMockTenantRepo()
	patientID := uuid.New()
	adminID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	doctorID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)
	tenants.memberships[tenantID] = map[uuid.UUID]authdomain.UserTenant{
		doctorID: {UserID: doctorID, TenantID: tenantID, Role: "doctor", Name: "Doc"},
	}

	cmd := NewAssignDoctorCommand(repo, tenants)
	p := newPrincipal(adminID, tenantID, sharedauth.RoleAdministrator)

	cs, err := cmd.Execute(context.Background(), p, caseID, doctorID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.DoctorID == nil || *cs.DoctorID != doctorID {
		t.Error("DoctorID not set correctly")
	}
	if cs.Status != domain.CaseStatusAssigned {
		t.Errorf("Status = %v, want assigned", cs.Status)
	}
}

func TestAssignDoctorPatientRejected(t *testing.T) {
	repo := newMockRepo()
	tenants := newMockTenantRepo()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)

	cmd := NewAssignDoctorCommand(repo, tenants)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	_, err := cmd.Execute(context.Background(), p, caseID, uuid.New())
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestWriteDiagnosisSuccess(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	doctorID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	assignedCase(repo, caseID, tenantID, patientID, doctorID)

	cmd := NewWriteDiagnosisCommand(repo)
	p := newPrincipal(doctorID, tenantID, sharedauth.RoleDoctor)

	d, err := cmd.Execute(context.Background(), p, caseID, "viral infection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Notes != "viral infection" {
		t.Errorf("Notes = %v, want viral infection", d.Notes)
	}
}

func TestWriteDiagnosisAdminRejected(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	doctorID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	assignedCase(repo, caseID, tenantID, patientID, doctorID)

	cmd := NewWriteDiagnosisCommand(repo)
	p := newPrincipal(uuid.New(), tenantID, sharedauth.RoleAdministrator)

	_, err := cmd.Execute(context.Background(), p, caseID, "notes")
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestCloseCaseSuccess(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	doctorID := uuid.New()
	adminID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	diagnosedCase(repo, caseID, tenantID, patientID, doctorID)

	cmd := NewCloseCaseCommand(repo)
	p := newPrincipal(adminID, tenantID, sharedauth.RoleAdministrator)

	cs, err := cmd.Execute(context.Background(), p, caseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.Status != domain.CaseStatusClosed {
		t.Errorf("Status = %v, want closed", cs.Status)
	}
	if cs.ClosedAt == nil {
		t.Error("ClosedAt not set")
	}
}

func TestCloseCaseOpenRejected(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	adminID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)

	cmd := NewCloseCaseCommand(repo)
	p := newPrincipal(adminID, tenantID, sharedauth.RoleAdministrator)

	_, err := cmd.Execute(context.Background(), p, caseID)
	if err != ErrCaseNotDiagnosed {
		t.Errorf("expected ErrCaseNotDiagnosed, got %v", err)
	}
}

func TestGetCasePatientOwn(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)

	cmd := NewGetCaseQuery(repo)
	p := newPrincipal(patientID, tenantID, sharedauth.RolePatient)

	cs, err := cmd.Execute(context.Background(), p, caseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.ID != caseID {
		t.Errorf("ID = %v, want %v", cs.ID, caseID)
	}
}

func TestGetCaseDoctorUnassigned(t *testing.T) {
	repo := newMockRepo()
	patientID := uuid.New()
	doctorID := uuid.New()
	tenantID := uuid.New()
	caseID := uuid.New()
	openCase(repo, caseID, tenantID, patientID)

	cmd := NewGetCaseQuery(repo)
	p := newPrincipal(doctorID, tenantID, sharedauth.RoleDoctor)

	_, err := cmd.Execute(context.Background(), p, caseID)
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied, got %v", err)
	}
}

func TestListCasesPatientOwnCases(t *testing.T) {
	repo := newMockRepo()
	patientA := uuid.New()
	patientB := uuid.New()
	tenantID := uuid.New()

	openCase(repo, uuid.New(), tenantID, patientA)
	openCase(repo, uuid.New(), tenantID, patientA)
	openCase(repo, uuid.New(), tenantID, patientB)

	cmd := NewListCasesQuery(repo)
	p := newPrincipal(patientA, tenantID, sharedauth.RolePatient)

	cases, total, err := cmd.Execute(context.Background(), p, "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(cases) != 2 {
		t.Errorf("len = %d, want 2", len(cases))
	}
}
