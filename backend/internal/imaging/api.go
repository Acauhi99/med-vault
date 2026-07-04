package imaging

import (
	"encoding/json"
	"net/http"
	"strings"

	clinicaldomain "github.com/Acauhi99/med-vault/internal/clinical/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	"github.com/Acauhi99/med-vault/internal/imaging/application"
	imagingdomain "github.com/Acauhi99/med-vault/internal/imaging/domain"
	"github.com/Acauhi99/med-vault/internal/shared/auditlog"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/Acauhi99/med-vault/internal/shared/httpx"
	"github.com/google/uuid"
)

type API struct {
	requestUploadURL *application.RequestUploadURLCommand
	confirmUpload    *application.ConfirmUploadCommand
	listImages       *application.ListImagesQuery
	getDownloadURL   *application.GetDownloadURLCommand
	AuditLogger      *auditlog.Logger
}

func NewAPI(repo imagingdomain.Repository, caseRepo clinicaldomain.Repository, storage application.Storage) *API {
	return &API{
		requestUploadURL: application.NewRequestUploadURLCommand(repo, caseRepo, storage),
		confirmUpload:    application.NewConfirmUploadCommand(repo, caseRepo),
		listImages:       application.NewListImagesQuery(repo, caseRepo),
		getDownloadURL:   application.NewGetDownloadURLCommand(repo, caseRepo, storage),
	}
}

func (a *API) RequestUploadURL(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.UploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if strings.TrimSpace(input.FileName) == "" || !input.ContentType.Valid() {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid upload request")
		return
	}

	url, s3Key, expiresIn, err := a.requestUploadURL.Execute(r.Context(), principal, id, input.FileName, string(input.ContentType))
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only patients can upload images")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "access denied")
		case application.ErrCaseNotOpen:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_OPEN", "case is not in open status")
		case application.ErrInvalidFileExt:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "unsupported file extension")
		case application.ErrInvalidContentType, application.ErrInvalidFileName:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate upload URL")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"upload_url": url,
		"s3_key":     s3Key,
		"expires_in": expiresIn,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "image.upload_requested", "case", id, nil)
	}
}

func (a *API) ConfirmUpload(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.ConfirmUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if strings.TrimSpace(input.S3Key) == "" || strings.TrimSpace(input.FileName) == "" || !input.ContentType.Valid() {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid upload confirmation")
		return
	}

	img, err := a.confirmUpload.Execute(r.Context(), principal, id, input.S3Key, input.FileName, string(input.ContentType))
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only patients can confirm uploads")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "access denied")
		case application.ErrCaseNotOpen:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_OPEN", "case is not in open status")
		case application.ErrImageNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "image not found")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to confirm upload")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"id":           img.ID,
		"case_id":      img.CaseID,
		"file_name":    img.FileName,
		"content_type": img.ContentType,
		"uploaded_at":  img.UploadedAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "image.upload_confirmed", "case", id, nil)
	}
}

func (a *API) ListImages(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	images, err := a.listImages.Execute(r.Context(), principal, id)
	if err != nil {
		switch err {
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "access denied")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list images")
		}
		return
	}

	summaries := make([]generated.ImageResponse, len(images))
	for i, img := range images {
		summaries[i] = generated.ImageResponse{
			Id:          &img.ID,
			CaseId:      &img.CaseID,
			FileName:    &img.FileName,
			ContentType: &img.ContentType,
			UploadedAt:  &img.UploadedAt,
		}
	}

	httpx.WriteJSON(w, r, http.StatusOK, summaries)
}

func (a *API) GetDownloadURL(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	url, expiresIn, err := a.getDownloadURL.Execute(r.Context(), principal, id)
	if err != nil {
		switch err {
		case application.ErrImageNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "image not found")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "access denied")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate download URL")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"download_url": url,
		"expires_in":   expiresIn,
	})
}
