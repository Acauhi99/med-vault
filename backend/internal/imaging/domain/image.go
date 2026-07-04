package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	PatientID   uuid.UUID
	CaseID      uuid.UUID
	FileName    string
	ContentType string
	S3Key       string
	UploadedAt  time.Time
}

type Repository interface {
	Create(ctx context.Context, img *Image) error
	ListByCase(ctx context.Context, tenantID, caseID uuid.UUID) ([]Image, error)
	GetByID(ctx context.Context, tenantID, imageID uuid.UUID) (*Image, error)
	FindByS3Key(ctx context.Context, tenantID, caseID uuid.UUID, s3Key string) (*Image, error)
	Delete(ctx context.Context, tenantID, imageID uuid.UUID) error
}
