package auth

import (
	"context"

	"github.com/google/uuid"
)

type Role string

const (
	RolePatient       Role = "patient"
	RoleDoctor        Role = "doctor"
	RoleAdministrator Role = "administrator"
)

type Principal struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Role     Role
}

type principalKey struct{}

func ContextWithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}

func MustPrincipal(ctx context.Context) Principal {
	p, ok := PrincipalFromContext(ctx)
	if !ok {
		panic("auth: principal not in context")
	}
	return p
}
