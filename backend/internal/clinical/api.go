package clinical

import (
	"encoding/json"
	"net/http"
	"strings"

	authdomain "github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/Acauhi99/med-vault/internal/clinical/application"
	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	"github.com/Acauhi99/med-vault/internal/generated"
	"github.com/Acauhi99/med-vault/internal/shared/auditlog"
	"github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/Acauhi99/med-vault/internal/shared/httpx"
)

type API struct {
	createCase     *application.CreateCaseCommand
	addSymptom     *application.AddSymptomCommand
	assignDoctor   *application.AssignDoctorCommand
	writeDiagnosis *application.WriteDiagnosisCommand
	closeCase      *application.CloseCaseCommand
	getCase        *application.GetCaseQuery
	listCases      *application.ListCasesQuery
	AuditLogger    *auditlog.Logger
}

func NewAPI(repo domain.Repository, tenants authdomain.TenantRepository) *API {
	return &API{
		createCase:     application.NewCreateCaseCommand(repo),
		addSymptom:     application.NewAddSymptomCommand(repo),
		assignDoctor:   application.NewAssignDoctorCommand(repo, tenants),
		writeDiagnosis: application.NewWriteDiagnosisCommand(repo),
		closeCase:      application.NewCloseCaseCommand(repo),
		getCase:        application.NewGetCaseQuery(repo),
		listCases:      application.NewListCasesQuery(repo),
	}
}

func (a *API) ListCases(w http.ResponseWriter, r *http.Request, params generated.ListCasesParams) {
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

	var status string
	if params.Status != nil {
		status = string(*params.Status)
	}

	cases, total, err := a.listCases.Execute(r.Context(), principal, status, page, pageSize)
	if err != nil {
		switch err {
		case application.ErrInvalidPageSize, application.ErrInvalidStatus:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", err.Error())
		}
		return
	}

	summaries := make([]generated.CaseSummary, len(cases))
	for i, c := range cases {
		status := generated.CaseSummaryStatus(c.Status)
		createdAt := c.CreatedAt
		summaries[i] = generated.CaseSummary{
			Id:        &c.ID,
			PatientId: &c.PatientID,
			DoctorId:  c.DoctorID,
			Status:    &status,
			CreatedAt: &createdAt,
		}
	}

	httpx.WriteJSONWithMeta(w, r, http.StatusOK, summaries, httpx.Meta{Page: &page, PageSize: &pageSize, Total: &total})
}

func (a *API) CreateCase(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.CreateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	for _, symptom := range input.Symptoms {
		if !symptom.Severity.Valid() || strings.TrimSpace(symptom.Description) == "" {
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid symptom")
			return
		}
	}

	symptoms := make([]application.CreateSymptomRequest, len(input.Symptoms))
	for i, s := range input.Symptoms {
		symptoms[i] = application.CreateSymptomRequest{
			Description: s.Description,
			Severity:    string(s.Severity),
		}
	}

	cs, err := a.createCase.Execute(r.Context(), principal, symptoms)
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only patients can create cases")
		case application.ErrNoSymptoms:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "at least one symptom is required")
		case application.ErrInvalidSymptom:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid symptom severity")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create case")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"id":         cs.ID,
		"patient_id": cs.PatientID,
		"status":     cs.Status,
		"created_at": cs.CreatedAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "case.created", "case", cs.ID, nil)
	}
}

func (a *API) GetCase(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	cs, err := a.getCase.Execute(r.Context(), principal, id)
	if err != nil {
		switch err {
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "access denied")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get case")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"id":         cs.ID,
		"tenant_id":  cs.TenantID,
		"patient_id": cs.PatientID,
		"doctor_id":  cs.DoctorID,
		"status":     cs.Status,
		"symptoms":   cs.Symptoms,
		"diagnosis":  cs.Diagnosis,
		"created_at": cs.CreatedAt,
		"updated_at": cs.UpdatedAt,
		"closed_at":  cs.ClosedAt,
	})
}

func (a *API) AddSymptom(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.AddSymptomRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if strings.TrimSpace(input.Description) == "" || !input.Severity.Valid() {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid symptom")
		return
	}

	s, err := a.addSymptom.Execute(r.Context(), principal, id, input.Description, string(input.Severity))
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only patients can add symptoms")
		case application.ErrCaseNotFound, application.ErrNotCasePatient:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrCaseNotOpen:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_OPEN", "case is not in open status")
		case application.ErrInvalidSymptom:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid symptom severity")
		case application.ErrInvalidDiagnosis:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add symptom")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"id":          s.ID,
		"description": s.Description,
		"severity":    s.Severity,
		"reported_at": s.ReportedAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "case.symptom_added", "case", id, nil)
	}
}

func (a *API) AssignDoctor(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.AssignDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	cs, err := a.assignDoctor.Execute(r.Context(), principal, id, input.DoctorId)
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only administrators can assign doctors")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrCaseNotOpen:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_OPEN", "case is not in open status")
		case application.ErrAccessDenied:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "doctor must belong to this tenant")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to assign doctor")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"id":         cs.ID,
		"doctor_id":  cs.DoctorID,
		"status":     cs.Status,
		"updated_at": cs.UpdatedAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "case.doctor_assigned", "case", cs.ID, nil)
	}
}

func (a *API) WriteDiagnosis(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	var input generated.WriteDiagnosisRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if strings.TrimSpace(input.Notes) == "" {
		httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "diagnosis notes are required")
		return
	}

	d, err := a.writeDiagnosis.Execute(r.Context(), principal, id, input.Notes)
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only doctors can write diagnoses")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrCaseNotAssigned:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_ASSIGNED", "case is not in assigned status")
		case application.ErrNotAssignedDoctor:
			httpx.WriteError(w, r, http.StatusForbidden, "NOT_ASSIGNED", "user is not the assigned doctor")
		case application.ErrInvalidDiagnosis:
			httpx.WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to write diagnosis")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, map[string]any{
		"id":         d.ID,
		"doctor_id":  d.DoctorID,
		"notes":      d.Notes,
		"written_at": d.WrittenAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "case.diagnosis_written", "case", id, nil)
	}
}

func (a *API) CloseCase(w http.ResponseWriter, r *http.Request, id generated.CaseIdParam) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
		return
	}

	cs, err := a.closeCase.Execute(r.Context(), principal, id)
	if err != nil {
		switch err {
		case application.ErrInvalidRole:
			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "only administrators can close cases")
		case application.ErrCaseNotFound:
			httpx.WriteError(w, r, http.StatusNotFound, "NOT_FOUND", "case not found")
		case application.ErrCaseNotDiagnosed:
			httpx.WriteError(w, r, http.StatusBadRequest, "CASE_NOT_DIAGNOSED", "case is not in diagnosed status")
		default:
			httpx.WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to close case")
		}
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, map[string]any{
		"id":         cs.ID,
		"status":     cs.Status,
		"closed_at":  cs.ClosedAt,
		"updated_at": cs.UpdatedAt,
	})

	if a.AuditLogger != nil {
		a.AuditLogger.Log(r.Context(), r, "case.closed", "case", cs.ID, nil)
	}
}
