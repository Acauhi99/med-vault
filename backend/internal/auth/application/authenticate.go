package application

import (
	"errors"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type AuthenticateInput struct {
	Email    string
	Password string
}

type TenantInfo struct {
	TenantID   uuid.UUID `json:"tenant_id"`
	TenantName string    `json:"tenant_name"`
	Role       string    `json:"role"`
}

type AuthenticateOutput struct {
	AccessToken string       `json:"access_token"`
	Tenants     []TenantInfo `json:"tenants"`
}

type AuthenticateCommand struct {
	users   domain.UserRepository
	tenants domain.TenantRepository
	hasher  PasswordHasher
	jwtGen  JWTGenerator
	jwtTTL  time.Duration
}

func NewAuthenticateCommand(
	users domain.UserRepository,
	tenants domain.TenantRepository,
	hasher PasswordHasher,
	jwtGen JWTGenerator,
	jwtTTL time.Duration,
) *AuthenticateCommand {
	return &AuthenticateCommand{
		users:   users,
		tenants: tenants,
		hasher:  hasher,
		jwtGen:  jwtGen,
		jwtTTL:  jwtTTL,
	}
}

func (c *AuthenticateCommand) Execute(input AuthenticateInput) (AuthenticateOutput, error) {
	user, err := c.users.FindByEmail(input.Email)
	if err != nil {
		return AuthenticateOutput{}, ErrInvalidCredentials
	}

	if err := c.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return AuthenticateOutput{}, ErrInvalidCredentials
	}

	tenants, err := c.tenants.FindUserTenants(user.ID)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	tempToken, err := c.jwtGen.GenerateTemporaryToken(user.ID, c.jwtTTL)
	if err != nil {
		return AuthenticateOutput{}, err
	}

	tenantInfos := make([]TenantInfo, len(tenants))
	for i, t := range tenants {
		tenantInfos[i] = TenantInfo{
			TenantID:   t.TenantID,
			TenantName: t.Name,
			Role:       t.Role,
		}
	}

	return AuthenticateOutput{
		AccessToken: tempToken,
		Tenants:     tenantInfos,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")
