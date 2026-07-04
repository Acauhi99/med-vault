package application

import "github.com/google/uuid"

type Principal struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Role     string
}
