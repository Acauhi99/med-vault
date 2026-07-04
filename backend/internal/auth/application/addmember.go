package application

import (
	"context"
	"errors"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type AddMemberInput struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	Role     string
}

type AddMemberCommand struct {
	tenants domain.TenantRepository
}

func NewAddMemberCommand(tenants domain.TenantRepository) *AddMemberCommand {
	return &AddMemberCommand{tenants: tenants}
}

func (c *AddMemberCommand) Execute(ctx context.Context, principal Principal, input AddMemberInput) (domain.UserTenant, error) {
	if principal.Role != "administrator" || principal.TenantID != input.TenantID {
		return domain.UserTenant{}, ErrNotAdmin
	}

	if !validRole(input.Role) {
		return domain.UserTenant{}, ErrInvalidRole
	}

	if err := c.tenants.AddMember(ctx, input.TenantID, input.UserID, input.Role); err != nil {
		return domain.UserTenant{}, err
	}

	membership, err := c.tenants.FindUserTenant(input.UserID, input.TenantID)
	if err != nil {
		return domain.UserTenant{}, err
	}

	return *membership, nil
}

var (
	ErrNotAdmin    = errors.New("principal is not an administrator of this tenant")
	ErrInvalidRole = errors.New("invalid role: must be patient, doctor, or administrator")
)

func validRole(role string) bool {
	switch role {
	case "patient", "doctor", "administrator":
		return true
	}
	return false
}
