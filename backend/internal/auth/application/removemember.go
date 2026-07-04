package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type RemoveMemberCommand struct {
	tenants domain.TenantRepository
}

func NewRemoveMemberCommand(tenants domain.TenantRepository) *RemoveMemberCommand {
	return &RemoveMemberCommand{tenants: tenants}
}

func (c *RemoveMemberCommand) Execute(ctx context.Context, principal Principal, tenantID, userID uuid.UUID) error {
	if principal.Role != "administrator" || principal.TenantID != tenantID {
		return ErrNotAdmin
	}

	return c.tenants.RemoveMember(ctx, tenantID, userID)
}
