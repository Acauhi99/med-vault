package pgx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Acauhi99/med-vault/internal/audit/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	detailsJSON, err := json.Marshal(log.Details)
	if err != nil {
		return err
	}

	// ponytail: bypass RLS for audit writes
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
		return err
	}

	_, err = conn.Exec(ctx,
		`INSERT INTO audit_logs (id, tenant_id, user_id, action, resource_type, resource_id, details, ip_address, user_agent, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		log.ID, log.TenantID, log.UserID, log.Action, log.ResourceType, log.ResourceID, detailsJSON, log.IPAddress, log.UserAgent, log.CreatedAt)
	return err
}

func (r *AuditRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, offset, limit int, action string, userID *uuid.UUID, resourceType string, resourceID *uuid.UUID) ([]domain.AuditLog, int, error) {
	// ponytail: bypass RLS for audit queries — cross-tenant admin read
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
		return nil, 0, err
	}

	where := "tenant_id = $1"
	args := []any{tenantID}
	if action != "" {
		where += fmt.Sprintf(" AND action = $%d", len(args)+1)
		args = append(args, action)
	}
	if userID != nil {
		where += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, *userID)
	}
	if resourceType != "" {
		where += fmt.Sprintf(" AND resource_type = $%d", len(args)+1)
		args = append(args, resourceType)
	}
	if resourceID != nil {
		where += fmt.Sprintf(" AND resource_id = $%d", len(args)+1)
		args = append(args, *resourceID)
	}
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE ` + where

	var total int
	err = conn.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	lim := len(args) - 1
	off := len(args)
	selectQuery := `SELECT id, tenant_id, user_id, action, resource_type, resource_id, details, ip_address, user_agent, created_at
		 FROM audit_logs WHERE ` + where +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprint(lim) + ` OFFSET $` + fmt.Sprint(off)

	rows, err := conn.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		var detailsJSON []byte
		if err := rows.Scan(&l.ID, &l.TenantID, &l.UserID, &l.Action, &l.ResourceType, &l.ResourceID, &detailsJSON, &l.IPAddress, &l.UserAgent, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		if detailsJSON != nil {
			if err := json.Unmarshal(detailsJSON, &l.Details); err != nil {
				return nil, 0, err
			}
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
