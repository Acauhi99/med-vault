package audit

import (
	"net/http"

	"github.com/Acauhi99/med-vault/internal/audit/application"
	"github.com/Acauhi99/med-vault/internal/audit/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/Acauhi99/med-vault/internal/shared/httpx"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type API struct {
	listAuditLogs *application.ListAuditLogsQuery
}

func NewAPI(repo domain.Repository) *API {
	return &API{
		listAuditLogs: application.NewListAuditLogsQuery(repo),
	}
}

func (a *API) ListAuditLogs(w http.ResponseWriter, r *http.Request, params generated.ListAuditLogsParams) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	page, pageSize := 1, 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	if pageSize > 100 {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "page size must be between 1 and 100")
		return
	}

	var resourceType string
	if params.ResourceType != nil {
		resourceType = *params.ResourceType
	}
	var resourceID *uuid.UUID
	if params.ResourceId != nil {
		id := uuid.UUID(*params.ResourceId)
		resourceID = &id
	}

	logs, total, err := a.listAuditLogs.Execute(r.Context(), principal, page, pageSize, resourceType, resourceID)
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only administrators can list audit logs")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list audit logs")
		}
		return
	}

	summaries := make([]generated.AuditLogResponse, len(logs))
	for i, log := range logs {
		action := log.Action
		resourceType := log.ResourceType
		ip := log.IPAddress
		meta := map[string]interface{}(log.Details)
		tenantID := openapi_types.UUID(log.TenantID)
		userID := openapi_types.UUID(log.UserID)
		resourceUUID := openapi_types.UUID(log.ResourceID)
		summaries[i] = generated.AuditLogResponse{
			Id:           &log.ID,
			TenantId:     &tenantID,
			UserId:       &userID,
			Action:       &action,
			ResourceType: &resourceType,
			ResourceId:   &resourceUUID,
			Metadata:     &meta,
			IpAddress:    &ip,
			CreatedAt:    &log.CreatedAt,
		}
	}

	httpx.WriteJSONWithMeta(w, r, http.StatusOK, summaries, httpx.Meta{Page: &page, PageSize: &pageSize, Total: &total})
}
