package pgx

import (
	"context"
	"errors"

	"github.com/Acauhi99/med-vault/internal/imaging/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("image not found")

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) Create(ctx context.Context, img *domain.Image) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO images (id, tenant_id, patient_id, case_id, file_name, content_type, s3_key, uploaded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		img.ID, img.TenantID, img.PatientID, img.CaseID, img.FileName, img.ContentType, img.S3Key, img.UploadedAt)
	return err
}

func (r *ImageRepository) ListByCase(ctx context.Context, tenantID, caseID uuid.UUID) ([]domain.Image, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT i.id, i.tenant_id, i.patient_id, i.case_id, i.file_name, i.content_type, i.s3_key, i.uploaded_at
		 FROM images i
		 JOIN cases c ON c.id = i.case_id
		 WHERE i.case_id = $1 AND c.tenant_id = $2
		 ORDER BY i.uploaded_at DESC`,
		caseID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []domain.Image
	for rows.Next() {
		var img domain.Image
		if err := rows.Scan(&img.ID, &img.TenantID, &img.PatientID, &img.CaseID, &img.FileName, &img.ContentType, &img.S3Key, &img.UploadedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *ImageRepository) GetByID(ctx context.Context, tenantID, imageID uuid.UUID) (*domain.Image, error) {
	var img domain.Image
	err := r.pool.QueryRow(ctx,
		`SELECT i.id, i.tenant_id, i.patient_id, i.case_id, i.file_name, i.content_type, i.s3_key, i.uploaded_at
		 FROM images i
		 JOIN cases c ON c.id = i.case_id
		 WHERE i.id = $1 AND c.tenant_id = $2`,
		imageID, tenantID,
	).Scan(&img.ID, &img.TenantID, &img.PatientID, &img.CaseID, &img.FileName, &img.ContentType, &img.S3Key, &img.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepository) FindByS3Key(ctx context.Context, tenantID, caseID uuid.UUID, s3Key string) (*domain.Image, error) {
	var img domain.Image
	err := r.pool.QueryRow(ctx,
		`SELECT i.id, i.tenant_id, i.patient_id, i.case_id, i.file_name, i.content_type, i.s3_key, i.uploaded_at
		 FROM images i
		 JOIN cases c ON c.id = i.case_id
		 WHERE i.s3_key = $1 AND i.case_id = $2 AND c.tenant_id = $3`,
		s3Key, caseID, tenantID,
	).Scan(&img.ID, &img.TenantID, &img.PatientID, &img.CaseID, &img.FileName, &img.ContentType, &img.S3Key, &img.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}
