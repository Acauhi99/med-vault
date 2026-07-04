package application

import (
	"context"
	"errors"
	"time"

	"github.com/Acauhi99/med-vault/internal/audit/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

var (
	ErrInvalidRole = errors.New("invalid role for this operation")
)

type LogActionCommand struct {
	repo domain.Repository
}

func NewLogActionCommand(repo domain.Repository) *LogActionCommand {
	return &LogActionCommand{repo: repo}
}

type LogActionInput struct {
	TenantID     uuid.UUID
	UserID       uuid.UUID
	Action       string
	ResourceType string
	ResourceID   uuid.UUID
	Details      map[string]any
	IPAddress    string
	UserAgent    string
}

func (c *LogActionCommand) Execute(ctx context.Context, input LogActionInput) (*domain.AuditLog, error) {
	log := &domain.AuditLog{
		ID:           uuid.New(),
		TenantID:     input.TenantID,
		UserID:       input.UserID,
		Action:       input.Action,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Details:      input.Details,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
		CreatedAt:    time.Now().UTC(),
	}

	if err := c.repo.Create(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}

type ListAuditLogsQuery struct {
	repo domain.Repository
}

func NewListAuditLogsQuery(repo domain.Repository) *ListAuditLogsQuery {
	return &ListAuditLogsQuery{repo: repo}
}

func (q *ListAuditLogsQuery) Execute(ctx context.Context, principal sharedauth.Principal, page, pageSize int, resourceType string, resourceID *uuid.UUID) ([]domain.AuditLog, int, error) {
	if principal.Role != sharedauth.RoleAdministrator {
		return nil, 0, ErrInvalidRole
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return q.repo.ListByTenant(ctx, principal.TenantID, offset, pageSize, resourceType, resourceID)
}
