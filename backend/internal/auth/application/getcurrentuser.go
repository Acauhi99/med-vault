package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type GetCurrentUserQuery struct {
	users   domain.UserRepository
	tenants domain.TenantRepository
}

func NewGetCurrentUserQuery(users domain.UserRepository, tenants domain.TenantRepository) *GetCurrentUserQuery {
	return &GetCurrentUserQuery{users: users, tenants: tenants}
}

func (q *GetCurrentUserQuery) Execute(ctx context.Context, userID, tenantID uuid.UUID) (domain.User, string, error) {
	user, err := q.users.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, "", err
	}

	membership, err := q.tenants.FindUserTenant(userID, tenantID)
	if err != nil {
		return domain.User{}, "", err
	}

	return *user, membership.Role, nil
}
