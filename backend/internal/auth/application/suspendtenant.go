package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type SuspendTenantCommand struct {
	tenantRepo domain.TenantRepository
}

func NewSuspendTenantCommand(tenantRepo domain.TenantRepository) *SuspendTenantCommand {
	return &SuspendTenantCommand{tenantRepo: tenantRepo}
}

func (c *SuspendTenantCommand) Execute(ctx context.Context, principal Principal, tenantID uuid.UUID) (*domain.Tenant, error) {
	if principal.Role != "administrator" || principal.TenantID != tenantID {
		return nil, ErrNotAdmin
	}
	return c.tenantRepo.Suspend(ctx, tenantID)
}
