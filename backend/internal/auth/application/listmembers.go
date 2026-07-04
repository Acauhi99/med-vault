package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type ListMembersQuery struct {
	tenants domain.TenantRepository
}

func NewListMembersQuery(tenants domain.TenantRepository) *ListMembersQuery {
	return &ListMembersQuery{tenants: tenants}
}

func (q *ListMembersQuery) Execute(ctx context.Context, principal Principal, tenantID uuid.UUID) ([]domain.UserTenant, error) {
	if principal.Role != "administrator" || principal.TenantID != tenantID {
		return nil, ErrNotAdmin
	}

	membership, err := q.tenants.FindUserTenant(principal.UserID, tenantID)
	if err != nil || membership == nil {
		return nil, ErrUserNotMember
	}
	if membership.Role != "administrator" {
		return nil, ErrNotAdmin
	}

	return q.tenants.ListMembers(ctx, tenantID)
}
