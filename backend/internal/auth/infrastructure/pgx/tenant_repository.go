package pgx

import (
	"context"
	"errors"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

func (r *TenantRepository) FindUserTenants(userID uuid.UUID) ([]domain.UserTenant, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT ut.user_id, ut.tenant_id, ut.role, t.name
		 FROM user_tenants ut
		 JOIN tenants t ON t.id = ut.tenant_id
		 WHERE ut.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []domain.UserTenant
	for rows.Next() {
		var ut domain.UserTenant
		if err := rows.Scan(&ut.UserID, &ut.TenantID, &ut.Role, &ut.Name); err != nil {
			return nil, err
		}
		tenants = append(tenants, ut)
	}
	return tenants, rows.Err()
}

func (r *TenantRepository) FindUserTenant(userID, tenantID uuid.UUID) (*domain.UserTenant, error) {
	row := r.pool.QueryRow(context.Background(),
		`SELECT ut.user_id, ut.tenant_id, ut.role, t.name
		 FROM user_tenants ut
		 JOIN tenants t ON t.id = ut.tenant_id
		 WHERE ut.user_id = $1 AND ut.tenant_id = $2`, userID, tenantID)

	var ut domain.UserTenant
	if err := row.Scan(&ut.UserID, &ut.TenantID, &ut.Role, &ut.Name); err != nil {
		return nil, err
	}
	return &ut, nil
}

func (r *TenantRepository) AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_tenants (user_id, tenant_id, role)
		 VALUES ($1, $2, $3)`, userID, tenantID, role)
	return err
}

func (r *TenantRepository) RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_tenants
		 WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID)
	return err
}

func (r *TenantRepository) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]domain.UserTenant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ut.user_id, ut.tenant_id, ut.role, t.name
		 FROM user_tenants ut
		 JOIN tenants t ON t.id = ut.tenant_id
		 WHERE ut.tenant_id = $1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.UserTenant
	for rows.Next() {
		var ut domain.UserTenant
		if err := rows.Scan(&ut.UserID, &ut.TenantID, &ut.Role, &ut.Name); err != nil {
			return nil, err
		}
		members = append(members, ut)
	}
	return members, rows.Err()
}

func (r *TenantRepository) Reactivate(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.pool.QueryRow(
		ctx,
		`UPDATE tenants SET status = 'active', updated_at = NOW()
		 WHERE id = $1 AND status = 'suspended'
		 RETURNING id, name, status, created_at, updated_at`, tenantID,
	).Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("tenant not found or not suspended")
		}
		return nil, err
	}
	return &t, nil
}

func (r *TenantRepository) Create(ctx context.Context, name string) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO tenants (name, status)
		 VALUES ($1, 'active')
		 RETURNING id, name, status, created_at, updated_at`, name,
	).Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TenantRepository) Suspend(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	var t domain.Tenant
	err := r.pool.QueryRow(
		ctx,
		`UPDATE tenants SET status = 'suspended', updated_at = NOW()
		 WHERE id = $1 AND status = 'active'
		 RETURNING id, name, status, created_at, updated_at`, tenantID,
	).Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("tenant not found or not active")
		}
		return nil, err
	}
	return &t, nil
}
