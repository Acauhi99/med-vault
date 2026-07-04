package application

import (
	"errors"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type SelectTenantInput struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
}

type SelectTenantOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type SelectTenantCommand struct {
	tenants    domain.TenantRepository
	jwtGen     JWTGenerator
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewSelectTenantCommand(
	tenants domain.TenantRepository,
	jwtGen JWTGenerator,
	accessTTL, refreshTTL time.Duration,
) *SelectTenantCommand {
	return &SelectTenantCommand{
		tenants:    tenants,
		jwtGen:     jwtGen,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (c *SelectTenantCommand) Execute(input SelectTenantInput) (SelectTenantOutput, error) {
	membership, err := c.tenants.FindUserTenant(input.UserID, input.TenantID)
	if err != nil {
		return SelectTenantOutput{}, ErrUserNotMember
	}

	accessToken, err := c.jwtGen.GenerateAccessToken(membership.UserID, membership.TenantID, membership.Role, c.accessTTL)
	if err != nil {
		return SelectTenantOutput{}, err
	}

	refreshToken, err := c.jwtGen.GenerateAccessToken(membership.UserID, membership.TenantID, membership.Role, c.refreshTTL)
	if err != nil {
		return SelectTenantOutput{}, err
	}

	return SelectTenantOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(c.accessTTL.Seconds()),
	}, nil
}

var ErrUserNotMember = errors.New("user is not a member of this tenant")
