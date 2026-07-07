package server

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"

	auditapi "github.com/Acauhi99/med-vault/internal/audit"
	"github.com/Acauhi99/med-vault/internal/auth/application"
	clinicalapi "github.com/Acauhi99/med-vault/internal/clinical"
	"github.com/Acauhi99/med-vault/internal/generated"
	imagingapi "github.com/Acauhi99/med-vault/internal/imaging"
	"github.com/Acauhi99/med-vault/internal/shared/auditlog"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/Acauhi99/med-vault/internal/shared/httpx"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type API struct {
	register            *application.RegisterCommand
	authenticate        *application.AuthenticateCommand
	selectTenant        *application.SelectTenantCommand
	refreshToken        *application.RefreshTokenCommand
	getCurrentUser      *application.GetCurrentUserQuery
	addMember           *application.AddMemberCommand
	removeMember        *application.RemoveMemberCommand
	listMembers         *application.ListMembersQuery
	reactivateTenantCmd *application.ReactivateTenantCommand
	createTenantCmd     *application.CreateTenantCommand
	suspendTenantCmd    *application.SuspendTenantCommand
	tokenStore          *application.TokenStore
	logger              *slog.Logger
	ClinicalAPI         *clinicalapi.API
	ImagingAPI          *imagingapi.API
	AuditAPI            *auditapi.API
	auditLog            *auditlog.Logger
}

func NewAPI(
	register *application.RegisterCommand,
	authenticate *application.AuthenticateCommand,
	selectTenant *application.SelectTenantCommand,
	refreshToken *application.RefreshTokenCommand,
	getCurrentUser *application.GetCurrentUserQuery,
	addMember *application.AddMemberCommand,
	removeMember *application.RemoveMemberCommand,
	listMembers *application.ListMembersQuery,
	reactivateTenantCmd *application.ReactivateTenantCommand,
	createTenantCmd *application.CreateTenantCommand,
	suspendTenantCmd *application.SuspendTenantCommand,
	tokenStore *application.TokenStore,
	logger *slog.Logger,
) *API {
	return &API{
		register:            register,
		authenticate:        authenticate,
		selectTenant:        selectTenant,
		refreshToken:        refreshToken,
		getCurrentUser:      getCurrentUser,
		addMember:           addMember,
		removeMember:        removeMember,
		listMembers:         listMembers,
		reactivateTenantCmd: reactivateTenantCmd,
		createTenantCmd:     createTenantCmd,
		suspendTenantCmd:    suspendTenantCmd,
		tokenStore:          tokenStore,
		logger:              logger,
	}
}

func (a *API) ListAuditLogs(w http.ResponseWriter, r *http.Request, params generated.ListAuditLogsParams) {
	a.AuditAPI.ListAuditLogs(w, r, params)
}

func (a *API) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var input generated.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	output, err := a.authenticate.Execute(application.AuthenticateInput{
		Email:    string(input.Email),
		Password: input.Password,
	})
	if err != nil {
		a.logLoginFailure(r, string(input.Email), err)
		httpx.WriteError(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "user.authenticated", "user", uuid.UUID{}, nil)
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"access_token": output.AccessToken,
		"tenants":      output.Tenants,
	})
}

func (a *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input generated.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	output, err := a.refreshToken.Execute(application.RefreshTokenInput{
		RefreshToken: input.RefreshToken,
	})
	if err != nil {
		httpx.WriteError(w, r, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "invalid or expired refresh token")
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"expires_in":    output.ExpiresIn,
	})
}

func (a *API) LogoutUser(w http.ResponseWriter, r *http.Request) {
	var input generated.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	if a.tokenStore != nil && input.RefreshToken != "" {
		a.tokenStore.Revoke(input.RefreshToken)
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "user.logged_out", "user", uuid.UUID{}, nil)
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{"data": nil})
}

func (a *API) logLoginFailure(r *http.Request, email string, err error) {
	if a.logger == nil {
		return
	}

	ip := r.RemoteAddr
	if host, _, splitErr := net.SplitHostPort(r.RemoteAddr); splitErr == nil {
		ip = host
	}

	a.logger.Warn(
		"auth.login_failed",
		slog.String("request_id", httpx.RequestID(r)),
		slog.String("email", email),
		slog.String("remote_addr", ip),
		slog.String("user_agent", r.UserAgent()),
		slog.String("reason", err.Error()),
	)
}

func (a *API) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input generated.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	output, err := a.register.Execute(application.RegisterInput{
		Email:    string(input.Email),
		Password: input.Password,
	})
	if err != nil {
		if err == application.ErrEmailAlreadyExists {
			httpx.WriteError(w, r, http.StatusConflict, "EMAIL_EXISTS", "email already registered")
			return
		}
		if err == application.ErrWeakPassword {
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "user.registered", "user", output.ID, nil)
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"id":         output.ID,
		"email":      output.Email,
		"status":     output.Status,
		"created_at": output.CreatedAt,
	})
}

func (a *API) SelectTenant(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok || principal.UserID == (openapi_types.UUID{}) {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.SelectTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	output, err := a.selectTenant.Execute(application.SelectTenantInput{
		UserID:   principal.UserID,
		TenantID: input.TenantId,
	})
	if err != nil {
		httpx.WriteError(w, r, http.StatusForbidden, "NOT_MEMBER", "user is not a member of this tenant")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "user.tenant_selected", "tenant", input.TenantId, nil)
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"expires_in":    output.ExpiresIn,
	})
}

func (a *API) ListCases(w http.ResponseWriter, r *http.Request, params generated.ListCasesParams) {
	a.ClinicalAPI.ListCases(w, r, params)
}

func (a *API) CreateCase(w http.ResponseWriter, r *http.Request) {
	a.ClinicalAPI.CreateCase(w, r)
}

func (a *API) GetCase(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ClinicalAPI.GetCase(w, r, id)
}

func (a *API) AssignDoctor(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ClinicalAPI.AssignDoctor(w, r, id)
}

func (a *API) CloseCase(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ClinicalAPI.CloseCase(w, r, id)
}

func (a *API) WriteDiagnosis(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ClinicalAPI.WriteDiagnosis(w, r, id)
}

func (a *API) ListImages(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ImagingAPI.ListImages(w, r, id)
}

func (a *API) ConfirmUpload(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ImagingAPI.ConfirmUpload(w, r, id)
}

func (a *API) RequestUploadURL(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ImagingAPI.RequestUploadURL(w, r, id)
}

func (a *API) AddSymptom(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	a.ClinicalAPI.AddSymptom(w, r, id)
}

func (a *API) GetDownloadURL(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	a.ImagingAPI.GetDownloadURL(w, r, id)
}

func (a *API) DeleteImage(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	a.ImagingAPI.DeleteImage(w, r, id)
}

func (a *API) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	user, role, err := a.getCurrentUser.Execute(principal.UserID, principal.TenantID)
	if err != nil {
		httpx.WriteError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"id":         user.ID,
		"tenant_id":  principal.TenantID,
		"email":      user.Email,
		"role":       role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}

func (a *API) ListTenantMembers(w http.ResponseWriter, r *http.Request, tenantId openapi_types.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	members, err := a.listMembers.Execute(r.Context(), appPrincipal, tenantId)
	if err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator of this tenant")
			return
		}
		httpx.WriteError(w, r, http.StatusForbidden, "NOT_MEMBER", "user is not a member of this tenant")
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, members)
}

func (a *API) AddTenantMember(w http.ResponseWriter, r *http.Request, tenantId openapi_types.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.AddTenantMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	member, err := a.addMember.Execute(r.Context(), appPrincipal, application.AddMemberInput{
		TenantID: tenantId,
		UserID:   input.UserId,
		Role:     string(input.Role),
	})
	if err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator of this tenant")
			return
		}
		if err == application.ErrInvalidRole {
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "role must be patient, doctor, or administrator")
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add member")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "tenant.member_added", "tenant", tenantId, nil)
	}

	httpx.WriteJSON(w, r, http.StatusCreated, member)
}

func (a *API) RemoveTenantMember(w http.ResponseWriter, r *http.Request, tenantId openapi_types.UUID, userId openapi_types.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	if err := a.removeMember.Execute(r.Context(), appPrincipal, tenantId, userId); err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator of this tenant")
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to remove member")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "tenant.member_removed", "tenant", tenantId, nil)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) ReactivateTenant(w http.ResponseWriter, r *http.Request, tenantId openapi_types.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	tenant, err := a.reactivateTenantCmd.Execute(r.Context(), appPrincipal, tenantId)
	if err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator of this tenant")
			return
		}
		httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "tenant not found or not suspended")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "tenant.reactivated", "tenant", tenantId, nil)
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"data": map[string]any{
			"id":         tenant.ID,
			"name":       tenant.Name,
			"status":     tenant.Status,
			"created_at": tenant.CreatedAt,
		},
	})
}

func (a *API) CreateTenant(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	if input.Name == "" {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "name is required")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	tenant, err := a.createTenantCmd.Execute(r.Context(), appPrincipal, input.Name)
	if err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator")
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create tenant")
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "tenant.created", "tenant", tenant.ID, nil)
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"data": map[string]any{
			"id":         tenant.ID,
			"name":       tenant.Name,
			"status":     tenant.Status,
			"created_at": tenant.CreatedAt,
		},
	})
}

func (a *API) SuspendTenant(w http.ResponseWriter, r *http.Request, tenantId openapi_types.UUID) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	appPrincipal := application.Principal{
		UserID:   principal.UserID,
		TenantID: principal.TenantID,
		Role:     string(principal.Role),
	}
	tenant, err := a.suspendTenantCmd.Execute(r.Context(), appPrincipal, tenantId)
	if err != nil {
		if err == application.ErrNotAdmin {
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ADMIN", "principal is not an administrator of this tenant")
			return
		}
		httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	if a.auditLog != nil {
		a.auditLog.Log(r.Context(), r, "tenant.suspended", "tenant", tenantId, nil)
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"data": map[string]any{
			"id":         tenant.ID,
			"name":       tenant.Name,
			"status":     tenant.Status,
			"created_at": tenant.CreatedAt,
		},
	})
}
