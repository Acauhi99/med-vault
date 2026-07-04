package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CaseStatus string

const (
	CaseStatusOpen      CaseStatus = "open"
	CaseStatusAssigned  CaseStatus = "assigned"
	CaseStatusDiagnosed CaseStatus = "diagnosed"
	CaseStatusClosed    CaseStatus = "closed"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Symptom struct {
	ID          uuid.UUID
	Description string
	Severity    Severity
	ReportedAt  time.Time
}

type Diagnosis struct {
	ID        uuid.UUID
	DoctorID  uuid.UUID
	Notes     string
	WrittenAt time.Time
}

type Case struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	PatientID uuid.UUID
	DoctorID  *uuid.UUID
	Status    CaseStatus
	Symptoms  []Symptom
	Diagnosis *Diagnosis
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time
}

func ValidateTransition(from, to CaseStatus) error {
	switch from {
	case CaseStatusOpen:
		if to == CaseStatusAssigned || to == CaseStatusOpen {
			return nil
		}
	case CaseStatusAssigned:
		if to == CaseStatusDiagnosed || to == CaseStatusAssigned {
			return nil
		}
	case CaseStatusDiagnosed:
		if to == CaseStatusClosed {
			return nil
		}
	}
	return fmt.Errorf("invalid transition: %s → %s", from, to)
}

type Repository interface {
	Create(ctx context.Context, c *Case) error
	GetByID(ctx context.Context, tenantID, caseID uuid.UUID) (*Case, error)
	ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID, status string, offset, limit int) ([]Case, int, error)
	ListByDoctor(ctx context.Context, tenantID, doctorID uuid.UUID, status string, offset, limit int) ([]Case, int, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, status string, offset, limit int) ([]Case, int, error)
	Update(ctx context.Context, c *Case) error
	AddSymptom(ctx context.Context, tenantID, caseID uuid.UUID, s *Symptom) error
	WriteDiagnosis(ctx context.Context, tenantID, caseID uuid.UUID, d *Diagnosis) error
}
