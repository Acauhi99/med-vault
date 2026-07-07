package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserTenant struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Role     string
	Name     string // tenant name
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
	Create(user *User) error
}

type TenantRepository interface {
	FindUserTenants(userID uuid.UUID) ([]UserTenant, error)
	FindUserTenant(userID, tenantID uuid.UUID) (*UserTenant, error)
	FindByName(name string) (*Tenant, error)
	AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error
	ListMembers(ctx context.Context, tenantID uuid.UUID) ([]UserTenant, error)
	Reactivate(ctx context.Context, tenantID uuid.UUID) (*Tenant, error)
	Create(ctx context.Context, name string) (*Tenant, error)
	Suspend(ctx context.Context, tenantID uuid.UUID) (*Tenant, error)
}
