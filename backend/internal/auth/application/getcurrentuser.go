package application

import (
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

func (q *GetCurrentUserQuery) Execute(userID, tenantID uuid.UUID) (domain.User, string, error) {
	user, err := q.users.FindByID(userID)
	if err != nil {
		return domain.User{}, "", err
	}

	membership, err := q.tenants.FindUserTenant(userID, tenantID)
	if err != nil {
		return domain.User{}, "", err
	}

	return *user, membership.Role, nil
}
