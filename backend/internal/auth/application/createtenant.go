package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
)

type CreateTenantCommand struct {
	tenantRepo domain.TenantRepository
}

func NewCreateTenantCommand(tenantRepo domain.TenantRepository) *CreateTenantCommand {
	return &CreateTenantCommand{tenantRepo: tenantRepo}
}

func (c *CreateTenantCommand) Execute(ctx context.Context, principal Principal, name string) (*domain.Tenant, error) {
	if principal.Role != "administrator" {
		return nil, ErrNotAdmin
	}
	return c.tenantRepo.Create(ctx, name)
}
