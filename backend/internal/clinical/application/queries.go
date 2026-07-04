package application

import (
	"context"

	"github.com/Acauhi99/med-vault/internal/clinical/domain"
	sharedauth "github.com/Acauhi99/med-vault/internal/shared/auth"
	"github.com/google/uuid"
)

type GetCaseQuery struct {
	repo domain.Repository
}

func NewGetCaseQuery(repo domain.Repository) *GetCaseQuery {
	return &GetCaseQuery{repo: repo}
}

func (q *GetCaseQuery) Execute(ctx context.Context, principal sharedauth.Principal, caseID uuid.UUID) (*domain.Case, error) {
	cs, err := q.repo.GetByID(ctx, principal.TenantID, caseID)
	if err != nil {
		return nil, ErrCaseNotFound
	}

	switch principal.Role {
	case sharedauth.RoleAdministrator:
	case sharedauth.RoleDoctor:
		if cs.DoctorID == nil || *cs.DoctorID != principal.UserID {
			return nil, ErrAccessDenied
		}
	case sharedauth.RolePatient:
		if cs.PatientID != principal.UserID {
			return nil, ErrAccessDenied
		}
	default:
		return nil, ErrAccessDenied
	}

	return cs, nil
}

type ListCasesQuery struct {
	repo domain.Repository
}

func NewListCasesQuery(repo domain.Repository) *ListCasesQuery {
	return &ListCasesQuery{repo: repo}
}

func (q *ListCasesQuery) Execute(ctx context.Context, principal sharedauth.Principal, status string, page, pageSize int) ([]domain.Case, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 100 {
		return nil, 0, ErrInvalidPageSize
	}
	if status != "" && !isValidCaseStatus(status) {
		return nil, 0, ErrInvalidStatus
	}

	offset := (page - 1) * pageSize

	switch principal.Role {
	case sharedauth.RolePatient:
		return q.repo.ListByPatient(ctx, principal.TenantID, principal.UserID, status, offset, pageSize)
	case sharedauth.RoleDoctor:
		return q.repo.ListByDoctor(ctx, principal.TenantID, principal.UserID, status, offset, pageSize)
	case sharedauth.RoleAdministrator:
		return q.repo.ListByTenant(ctx, principal.TenantID, status, offset, pageSize)
	default:
		return nil, 0, ErrAccessDenied
	}
}
