package pgx

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("case not found")

type CaseRepository struct {
	pool *pgxpool.Pool
}

func NewCaseRepository(pool *pgxpool.Pool) *CaseRepository {
	return &CaseRepository{pool: pool}
}

func (r *CaseRepository) Create(ctx context.Context, c *domain.Case) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx,
		`INSERT INTO cases (id, tenant_id, patient_id, doctor_id, status, created_at, updated_at, closed_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.TenantID, c.PatientID, c.DoctorID, c.Status, c.CreatedAt, c.UpdatedAt, c.ClosedAt)
	if err != nil {
		return err
	}

	for _, s := range c.Symptoms {
		_, err = tx.Exec(ctx,
			`INSERT INTO symptoms (id, case_id, description, severity, reported_at)
			 VALUES ($1, $2, $3, $4, $5)`,
			s.ID, c.ID, s.Description, s.Severity, s.ReportedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *CaseRepository) GetByID(ctx context.Context, tenantID, caseID uuid.UUID) (*domain.Case, error) {
	var c domain.Case
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, tenant_id, patient_id, doctor_id, status, created_at, updated_at, closed_at
		 FROM cases WHERE id = $1 AND tenant_id = $2`,
		caseID, tenantID,
	).Scan(&c.ID, &c.TenantID, &c.PatientID, &c.DoctorID, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.ClosedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, description, severity, reported_at
		 FROM symptoms WHERE case_id = $1`,
		caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.Symptom
		if err := rows.Scan(&s.ID, &s.Description, &s.Severity, &s.ReportedAt); err != nil {
			return nil, err
		}
		c.Symptoms = append(c.Symptoms, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var diag domain.Diagnosis
	err = r.pool.QueryRow(
		ctx,
		`SELECT id, doctor_id, notes, written_at
		 FROM diagnoses WHERE case_id = $1`,
		caseID,
	).Scan(&diag.ID, &diag.DoctorID, &diag.Notes, &diag.WrittenAt)
	if err == nil {
		c.Diagnosis = &diag
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &c, nil
}

func (r *CaseRepository) ListByPatient(ctx context.Context, tenantID, patientID uuid.UUID, status string, offset, limit int) ([]domain.Case, int, error) {
	return r.listCases(ctx,
		`patient_id = $1 AND tenant_id = $2`, []any{patientID, tenantID},
		status, offset, limit)
}

func (r *CaseRepository) ListByDoctor(ctx context.Context, tenantID, doctorID uuid.UUID, status string, offset, limit int) ([]domain.Case, int, error) {
	return r.listCases(ctx,
		`doctor_id = $1 AND tenant_id = $2`, []any{doctorID, tenantID},
		status, offset, limit)
}

func (r *CaseRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, status string, offset, limit int) ([]domain.Case, int, error) {
	return r.listCases(ctx,
		`tenant_id = $1`, []any{tenantID},
		status, offset, limit)
}

func (r *CaseRepository) listCases(ctx context.Context, whereClause string, baseArgs []any, status string, offset, limit int) ([]domain.Case, int, error) {
	args := make([]any, len(baseArgs))
	copy(args, baseArgs)

	statusFilter := ""
	if status != "" {
		n := len(args) + 1
		statusFilter = " AND status = $" + strconv.Itoa(n)
		args = append(args, status)
	}

	var total int
	err := r.pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM cases WHERE `+whereClause+statusFilter, args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	n := len(args)
	lim := strconv.Itoa(n - 1)
	off := strconv.Itoa(n)

	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, patient_id, doctor_id, status, created_at, updated_at, closed_at
		 FROM cases WHERE `+whereClause+statusFilter+` ORDER BY created_at DESC LIMIT $`+lim+` OFFSET $`+off,
		args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	cases, err := scanCases(rows)
	return cases, total, err
}

func (r *CaseRepository) Update(ctx context.Context, c *domain.Case) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE cases
		 SET doctor_id = $1, status = $2, closed_at = $3, updated_at = $4
		 WHERE id = $5 AND tenant_id = $6`,
		c.DoctorID, c.Status, c.ClosedAt, c.UpdatedAt, c.ID, c.TenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CaseRepository) AddSymptom(ctx context.Context, tenantID, caseID uuid.UUID, s *domain.Symptom) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO symptoms (id, case_id, description, severity, reported_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		s.ID, caseID, s.Description, s.Severity, s.ReportedAt)
	return err
}

func (r *CaseRepository) WriteDiagnosis(ctx context.Context, tenantID, caseID uuid.UUID, d *domain.Diagnosis) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO diagnoses (id, case_id, doctor_id, notes, written_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		d.ID, caseID, d.DoctorID, d.Notes, d.WrittenAt)
	return err
}

func scanCases(rows pgx.Rows) ([]domain.Case, error) {
	var cases []domain.Case
	for rows.Next() {
		var c domain.Case
		if err := rows.Scan(&c.ID, &c.TenantID, &c.PatientID, &c.DoctorID, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.ClosedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, rows.Err()
}

func init() {
	_ = fmt.Sprintf
}
