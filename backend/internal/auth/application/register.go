package application

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	ID        uuid.UUID
	Email     string
	Status    string
	CreatedAt time.Time
}

type RegisterCommand struct {
	users   domain.UserRepository
	tenants domain.TenantRepository
	hasher  PasswordHasher
}

func NewRegisterCommand(users domain.UserRepository, tenants domain.TenantRepository, hasher PasswordHasher) *RegisterCommand {
	return &RegisterCommand{users: users, tenants: tenants, hasher: hasher}
}

func (c *RegisterCommand) Execute(input RegisterInput) (RegisterOutput, error) {
	if !isStrongPassword(input.Password) {
		return RegisterOutput{}, ErrWeakPassword
	}

	existing, err := c.users.FindByEmail(input.Email)
	if err == nil && existing != nil {
		return RegisterOutput{}, ErrEmailAlreadyExists
	}

	hash, err := c.hasher.Hash(input.Password)
	if err != nil {
		return RegisterOutput{}, err
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: hash,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := c.users.Create(user); err != nil {
		return RegisterOutput{}, err
	}

	// Auto-provision default tenant for new users — find-or-create
	tenant, err := c.tenants.FindByName("MedVault Demo")
	if err != nil {
		tenant, err = c.tenants.Create(context.Background(), "MedVault Demo")
		if err != nil {
			return RegisterOutput{
				ID:        user.ID,
				Email:     user.Email,
				Status:    user.Status,
				CreatedAt: user.CreatedAt,
			}, nil
		}
	}

	_ = c.tenants.AddMember(context.Background(), tenant.ID, user.ID, "administrator")

	return RegisterOutput{
		ID:        user.ID,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

var ErrEmailAlreadyExists = errors.New("email already exists")

var ErrWeakPassword = errors.New("password must be at least 12 characters and include upper, lower, number, and special")

func isStrongPassword(password string) bool {
	if len(password) < 12 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case !unicode.IsLetter(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial && strings.TrimSpace(password) == password
}
