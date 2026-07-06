package application

import (
	"context"
	"errors"
	"strings"
	"time"

	authdomain "github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

var (
	ErrCaseNotFound      = errors.New("case not found")
	ErrAccessDenied      = errors.New("access denied")
	ErrInvalidRole       = errors.New("invalid role for this operation")
	ErrCaseNotOpen       = errors.New("case is not in open status")
	ErrCaseNotAssigned   = errors.New("case is not in assigned status")
	ErrCaseNotDiagnosed  = errors.New("case is not in diagnosed status")
	ErrNotCasePatient    = errors.New("user is not the patient of this case")
	ErrNotAssignedDoctor = errors.New("user is not the assigned doctor for this case")
	ErrNoSymptoms        = errors.New("at least one symptom is required")
	ErrInvalidSymptom    = errors.New("invalid symptom severity")
	ErrInvalidDiagnosis  = errors.New("invalid diagnosis notes")
	ErrInvalidPageSize   = errors.New("page size must be between 1 and 100")
	ErrInvalidStatus     = errors.New("invalid case status")
)

type CreateSymptomRequest struct {
	Description string
	Severity    string
}

type CreateCaseCommand struct {
	repo domain.Repository
}

func NewCreateCaseCommand(repo domain.Repository) *CreateCaseCommand {
	return &CreateCaseCommand{repo: repo}
}

func (c *CreateCaseCommand) Execute(ctx context.Context, principal sharedauth.Principal, symptoms []CreateSymptomRequest) (*domain.Case, error) {
	if principal.Role != sharedauth.RolePatient {
		return nil, ErrInvalidRole
	}
	if len(symptoms) == 0 {
		return nil, ErrNoSymptoms
	}
	for _, s := range symptoms {
		if strings.TrimSpace(s.Description) == "" || !isValidSeverity(s.Severity) {
			return nil, ErrInvalidSymptom
		}
	}

	now := time.Now().UTC()
	cs := &domain.Case{
		ID:        uuid.New(),
		TenantID:  principal.TenantID,
		PatientID: principal.UserID,
		Status:    domain.CaseStatusOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, s := range symptoms {
		cs.Symptoms = append(cs.Symptoms, domain.Symptom{
			ID:          uuid.New(),
			Description: s.Description,
			Severity:    domain.Severity(s.Severity),
			ReportedAt:  now,
		})
	}

	if err := c.repo.Create(ctx, cs); err != nil {
		return nil, err
	}

	return cs, nil
}

type AddSymptomCommand struct {
	repo domain.Repository
}

func NewAddSymptomCommand(repo domain.Repository) *AddSymptomCommand {
	return &AddSymptomCommand{repo: repo}
}

func (c *AddSymptomCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID, description, severity string) (*domain.Symptom, error) {
	if principal.Role != sharedauth.RolePatient {
		return nil, ErrInvalidRole
	}
	if !isValidSeverity(severity) {
		return nil, ErrInvalidSymptom
	}

	cs, err := c.repo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}
	if cs.PatientID != principal.UserID {
		return nil, ErrNotCasePatient
	}
	if cs.Status != domain.CaseStatusOpen {
		return nil, ErrCaseNotOpen
	}

	now := time.Now().UTC()
	s := &domain.Symptom{
		ID:          uuid.New(),
		Description: description,
		Severity:    domain.Severity(severity),
		ReportedAt:  now,
	}

	if err := c.repo.AddSymptom(ctx, principal.TenantID, caseID, s); err != nil {
		return nil, err
	}

	return s, nil
}

type AssignDoctorCommand struct {
	repo    domain.Repository
	tenants authdomain.TenantRepository
}

func NewAssignDoctorCommand(repo domain.Repository, tenants authdomain.TenantRepository) *AssignDoctorCommand {
	return &AssignDoctorCommand{repo: repo, tenants: tenants}
}

func (c *AssignDoctorCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID, doctorID uuid.UUID) (*domain.Case, error) {
	if principal.Role != sharedauth.RoleAdministrator {
		return nil, ErrInvalidRole
	}

	cs, err := c.repo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}
	if cs.Status != domain.CaseStatusOpen {
		return nil, ErrCaseNotOpen
	}
	if len(cs.Symptoms) == 0 {
		return nil, ErrNoSymptoms
	}
	if c.tenants != nil {
		membership, err := c.tenants.FindUserTenant(doctorID, principal.TenantID)
		if err != nil || membership == nil || membership.Role != "doctor" {
			return nil, ErrAccessDenied
		}
	}

	cs.DoctorID = &doctorID
	if err := domain.ValidateTransition(cs.Status, domain.CaseStatusAssigned); err != nil {
		return nil, err
	}
	cs.Status = domain.CaseStatusAssigned
	cs.UpdatedAt = time.Now().UTC()

	if err := c.repo.Update(ctx, cs); err != nil {
		return nil, err
	}

	return cs, nil
}

type WriteDiagnosisCommand struct {
	repo domain.Repository
}

func NewWriteDiagnosisCommand(repo domain.Repository) *WriteDiagnosisCommand {
	return &WriteDiagnosisCommand{repo: repo}
}

func (c *WriteDiagnosisCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID, notes string) (*domain.Diagnosis, error) {
	if principal.Role != sharedauth.RoleDoctor {
		return nil, ErrInvalidRole
	}

	cs, err := c.repo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}
	if cs.Status != domain.CaseStatusAssigned {
		return nil, ErrCaseNotAssigned
	}
	if cs.DoctorID == nil || *cs.DoctorID != principal.UserID {
		return nil, ErrNotAssignedDoctor
	}
	if strings.TrimSpace(notes) == "" {
		return nil, ErrInvalidDiagnosis
	}

	now := time.Now().UTC()
	d := &domain.Diagnosis{
		ID:        uuid.New(),
		DoctorID:  principal.UserID,
		Notes:     notes,
		WrittenAt: now,
	}

	if err := c.repo.WriteDiagnosis(ctx, principal.TenantID, caseID, d); err != nil {
		return nil, err
	}

	if err := domain.ValidateTransition(cs.Status, domain.CaseStatusDiagnosed); err != nil {
		return nil, err
	}
	cs.Status = domain.CaseStatusDiagnosed
	cs.Diagnosis = d
	cs.UpdatedAt = now

	if err := c.repo.Update(ctx, cs); err != nil {
		return nil, err
	}

	return d, nil
}

type CloseCaseCommand struct {
	repo domain.Repository
}

func NewCloseCaseCommand(repo domain.Repository) *CloseCaseCommand {
	return &CloseCaseCommand{repo: repo}
}

func (c *CloseCaseCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID) (*domain.Case, error) {
	if principal.Role != sharedauth.RoleAdministrator {
		return nil, ErrInvalidRole
	}

	cs, err := c.repo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}
	if cs.Status != domain.CaseStatusDiagnosed {
		return nil, ErrCaseNotDiagnosed
	}

	now := time.Now().UTC()
	if err := domain.ValidateTransition(cs.Status, domain.CaseStatusClosed); err != nil {
		return nil, err
	}
	cs.Status = domain.CaseStatusClosed
	cs.ClosedAt = &now
	cs.UpdatedAt = now

	if err := c.repo.Update(ctx, cs); err != nil {
		return nil, err
	}

	return cs, nil
}

func isValidSeverity(s string) bool {
	switch domain.Severity(s) {
	case domain.SeverityLow, domain.SeverityMedium, domain.SeverityHigh, domain.SeverityCritical:
		return true
	}
	return false
}

func isValidCaseStatus(status string) bool {
	switch domain.CaseStatus(status) {
	case domain.CaseStatusOpen, domain.CaseStatusAssigned, domain.CaseStatusDiagnosed, domain.CaseStatusClosed:
		return true
	}
	return false
}
