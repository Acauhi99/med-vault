package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	UserID       uuid.UUID
	Action       string
	ResourceType string
	ResourceID   uuid.UUID
	Details      map[string]any
	IPAddress    string
	UserAgent    string
	CreatedAt    time.Time
}

type Repository interface {
	Create(ctx context.Context, log *AuditLog) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID, offset, limit int, action string, userID *uuid.UUID, resourceType string, resourceID *uuid.UUID) ([]AuditLog, int, error)
}
