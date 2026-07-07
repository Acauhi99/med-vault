package application

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	clinicaldomain "github.com/Acauhi99/med-vault/internal/clinical/domain"
	imagingdomain "github.com/Acauhi99/med-vault/internal/imaging/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

var (
	ErrCaseNotFound       = errors.New("case not found")
	ErrAccessDenied       = errors.New("access denied")
	ErrInvalidRole        = errors.New("invalid role for this operation")
	ErrCaseNotOpen        = errors.New("case is not in open status")
	ErrImageNotFound      = errors.New("image not found")
	ErrInvalidFileExt     = errors.New("unsupported file extension")
	ErrInvalidContentType = errors.New("unsupported content type")
	ErrInvalidFileName    = errors.New("invalid file name")
	ErrFileTooLarge       = errors.New("file size exceeds 50MB limit")
)

const maxUploadSizeBytes int64 = 50 * 1024 * 1024

type RequestUploadURLCommand struct {
	repo     imagingdomain.Repository
	caseRepo clinicaldomain.Repository
	storage  Storage
}

func NewRequestUploadURLCommand(repo imagingdomain.Repository, caseRepo clinicaldomain.Repository, storage Storage) *RequestUploadURLCommand {
	return &RequestUploadURLCommand{repo: repo, caseRepo: caseRepo, storage: storage}
}

func (c *RequestUploadURLCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID, fileName, contentType string, fileSize int64) (uploadURL, s3Key string, expiresIn int, err error) {
	if principal.Role != sharedauth.RolePatient {
		return "", "", 0, ErrInvalidRole
	}
	if fileSize <= 0 || fileSize > maxUploadSizeBytes {
		return "", "", 0, ErrFileTooLarge
	}

	cs, err := c.caseRepo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return "", "", 0, ErrCaseNotFound
	}
	if cs.PatientID != principal.UserID {
		return "", "", 0, ErrAccessDenied
	}
	if cs.Status != clinicaldomain.CaseStatusOpen {
		return "", "", 0, ErrCaseNotOpen
	}
	if !isAllowedContentType(contentType) {
		return "", "", 0, ErrInvalidContentType
	}

	ext := path.Ext(fileName)
	if ext == "" {
		return "", "", 0, ErrInvalidFileExt
	}
	cleanName := path.Base(fileName)
	if cleanName == "." || cleanName == "/" || strings.TrimSpace(cleanName) == "" {
		return "", "", 0, ErrInvalidFileName
	}

	imageID := uuid.New()
	s3Key = principal.TenantID.String() + "/" + caseID.String() + "/" + imageID.String() + "/" + cleanName

	ttl := 15 * time.Minute
	url, err := c.storage.GenerateUploadURL(ctx, s3Key, contentType, ttl)
	if err != nil {
		return "", "", 0, err
	}

	if err := c.repo.Create(ctx, &imagingdomain.Image{
		ID:          imageID,
		TenantID:    principal.TenantID,
		PatientID:   principal.UserID,
		CaseID:      caseID,
		FileName:    cleanName,
		ContentType: contentType,
		S3Key:       s3Key,
		UploadedAt:  time.Now().UTC(),
	}); err != nil {
		return "", "", 0, err
	}

	return url, s3Key, int(ttl.Seconds()), nil
}

type ConfirmUploadCommand struct {
	repo     imagingdomain.Repository
	caseRepo clinicaldomain.Repository
}

func NewConfirmUploadCommand(repo imagingdomain.Repository, caseRepo clinicaldomain.Repository) *ConfirmUploadCommand {
	return &ConfirmUploadCommand{repo: repo, caseRepo: caseRepo}
}

func (c *ConfirmUploadCommand) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID, s3Key, fileName, contentType string) (*imagingdomain.Image, error) {
	if principal.Role != sharedauth.RolePatient {
		return nil, ErrInvalidRole
	}

	cs, err := c.caseRepo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}
	if cs.PatientID != principal.UserID {
		return nil, ErrAccessDenied
	}
	if cs.Status != clinicaldomain.CaseStatusOpen {
		return nil, ErrCaseNotOpen
	}

	img, err := c.repo.FindByS3Key(ctx, principal.TenantID, caseID, s3Key)
	if err != nil {
		return nil, ErrImageNotFound
	}
	if img.FileName != fileName || img.ContentType != contentType {
		return nil, ErrAccessDenied
	}

	return img, nil
}

type ListImagesQuery struct {
	repo     imagingdomain.Repository
	caseRepo clinicaldomain.Repository
}

func NewListImagesQuery(repo imagingdomain.Repository, caseRepo clinicaldomain.Repository) *ListImagesQuery {
	return &ListImagesQuery{repo: repo, caseRepo: caseRepo}
}

func (q *ListImagesQuery) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID) ([]imagingdomain.Image, error) {
	cs, err := q.caseRepo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}

	if !hasCaseAccess(principal, cs) {
		return nil, ErrAccessDenied
	}

	return q.repo.ListByCase(ctx, principal.TenantID, caseID)
}

type GetDownloadURLCommand struct {
	repo     imagingdomain.Repository
	caseRepo clinicaldomain.Repository
	storage  Storage
}

func NewGetDownloadURLCommand(repo imagingdomain.Repository, caseRepo clinicaldomain.Repository, storage Storage) *GetDownloadURLCommand {
	return &GetDownloadURLCommand{repo: repo, caseRepo: caseRepo, storage: storage}
}

func (c *GetDownloadURLCommand) Execute(ctx context.Context, principal sharedauth.Principal, imageID uuid.UUID) (downloadURL string, expiresIn int, err error) {
	img, err := c.repo.GetByID(ctx, principal.TenantID, imageID)
	if err != nil {
		return "", 0, ErrImageNotFound
	}

	cs, err := c.caseRepo.GetByID(ctx, principal.TenantID, img.CaseID)
	if err != nil {
		return "", 0, ErrCaseNotFound
	}

	if !hasCaseAccess(principal, cs) {
		return "", 0, ErrAccessDenied
	}

	ttl := 15 * time.Minute
	url, err := c.storage.GenerateDownloadURL(ctx, img.S3Key, ttl)
	if err != nil {
		return "", 0, err
	}

	return url, int(ttl.Seconds()), nil
}

func hasCaseAccess(principal sharedauth.Principal, cs *clinicaldomain.Case) bool {
	switch principal.Role {
	case sharedauth.RoleAdministrator:
		return true
	case sharedauth.RoleDoctor:
		return cs.DoctorID != nil && *cs.DoctorID == principal.UserID
	case sharedauth.RolePatient:
		return cs.PatientID == principal.UserID
	}
	return false
}

func isAllowedContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/dicom":
		return true
	}
	return false
}

type DeleteImageCommand struct {
	repo    imagingdomain.Repository
	storage Storage
}

func NewDeleteImageCommand(repo imagingdomain.Repository, storage Storage) *DeleteImageCommand {
	return &DeleteImageCommand{repo: repo, storage: storage}
}

func (c *DeleteImageCommand) Execute(ctx context.Context, principal sharedauth.Principal, imageID uuid.UUID) error {
	if principal.Role != sharedauth.RoleAdministrator {
		return ErrInvalidRole
	}

	img, err := c.repo.GetByID(ctx, principal.TenantID, imageID)
	if err != nil {
		return ErrImageNotFound
	}

	if err := c.storage.DeleteObject(ctx, img.S3Key); err != nil {
		return err
	}

	return c.repo.Delete(ctx, principal.TenantID, imageID)
}
