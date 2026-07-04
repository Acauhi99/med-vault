package application

import (
	"context"
	"testing"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type selectTenantMockJWTGen struct {
	accessCalls  int
	refreshCalls int
}

func (m *selectTenantMockJWTGen) GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	m.accessCalls++
	return "access-token", nil
}

func (m *selectTenantMockJWTGen) GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	m.refreshCalls++
	return "refresh-token", nil
}

func (m *selectTenantMockJWTGen) GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	return "temp-token", nil
}

func (m *selectTenantMockJWTGen) Verify(rawToken string) (JWTClaims, error) {
	return JWTClaims{}, nil
}

type selectTenantMockRepo struct {
	membership *domain.UserTenant
}

func (m *selectTenantMockRepo) FindUserTenants(userID uuid.UUID) ([]domain.UserTenant, error) {
	return nil, nil
}

func (m *selectTenantMockRepo) FindUserTenant(userID, tenantID uuid.UUID) (*domain.UserTenant, error) {
	if m.membership == nil {
		return nil, ErrUserNotMember
	}

	return m.membership, nil
}

func (m *selectTenantMockRepo) AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
	return nil
}

func (m *selectTenantMockRepo) RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error {
	return nil
}

func (m *selectTenantMockRepo) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]domain.UserTenant, error) {
	return nil, nil
}

func (m *selectTenantMockRepo) Reactivate(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return nil, nil
}

func TestSelectTenantUsesRefreshToken(t *testing.T) {
	userID := uuid.New()
	tenantID := uuid.New()
	repo := &selectTenantMockRepo{membership: &domain.UserTenant{UserID: userID, TenantID: tenantID, Role: "doctor"}}
	jwtGen := &selectTenantMockJWTGen{}
	cmd := NewSelectTenantCommand(repo, jwtGen, 15*time.Minute, 168*time.Hour)

	out, err := cmd.Execute(SelectTenantInput{UserID: userID, TenantID: tenantID})
	if err != nil {
		t.Fatalf("select tenant: %v", err)
	}
	if out.RefreshToken != "refresh-token" {
		t.Fatalf("refresh token = %q, want refresh-token", out.RefreshToken)
	}
	if jwtGen.accessCalls != 1 {
		t.Fatalf("access calls = %d, want 1", jwtGen.accessCalls)
	}
	if jwtGen.refreshCalls != 1 {
		t.Fatalf("refresh calls = %d, want 1", jwtGen.refreshCalls)
	}
}
