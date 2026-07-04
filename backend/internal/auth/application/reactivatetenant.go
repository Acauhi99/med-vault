package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type ReactivateTenantCommand struct {
	tenantRepo domain.TenantRepository
}

func NewReactivateTenantCommand(tenantRepo domain.TenantRepository) *ReactivateTenantCommand {
	return &ReactivateTenantCommand{tenantRepo: tenantRepo}
}

func (c *ReactivateTenantCommand) Execute(ctx context.Context, principal Principal, tenantID uuid.UUID) (*domain.Tenant, error) {
	if principal.Role != "administrator" || principal.TenantID != tenantID {
		return nil, ErrNotAdmin
	}
	return c.tenantRepo.Reactivate(ctx, tenantID)
}
