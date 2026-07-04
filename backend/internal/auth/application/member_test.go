package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type mockTenantMemberRepo struct {
	memberships map[string]domain.UserTenant
	byTenant    map[uuid.UUID][]domain.UserTenant
}

func newMockTenantMemberRepo() *mockTenantMemberRepo {
	return &mockTenantMemberRepo{
		memberships: make(map[string]domain.UserTenant),
		byTenant:    make(map[uuid.UUID][]domain.UserTenant),
	}
}

func (m *mockTenantMemberRepo) FindUserTenants(userID uuid.UUID) ([]domain.UserTenant, error) {
	var result []domain.UserTenant
	for _, ut := range m.memberships {
		if ut.UserID == userID {
			result = append(result, ut)
		}
	}
	return result, nil
}

func (m *mockTenantMemberRepo) FindUserTenant(userID, tenantID uuid.UUID) (*domain.UserTenant, error) {
	key := tenantID.String() + ":" + userID.String()
	ut, ok := m.memberships[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return &ut, nil
}

func (m *mockTenantMemberRepo) AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
	key := tenantID.String() + ":" + userID.String()
	m.memberships[key] = domain.UserTenant{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Name:     "Test Tenant",
	}
	m.byTenant[tenantID] = append(m.byTenant[tenantID], domain.UserTenant{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Name:     "Test Tenant",
	})
	return nil
}

func (m *mockTenantMemberRepo) RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error {
	key := tenantID.String() + ":" + userID.String()
	delete(m.memberships, key)
	return nil
}

func (m *mockTenantMemberRepo) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]domain.UserTenant, error) {
	return m.byTenant[tenantID], nil
}

func (m *mockTenantMemberRepo) Reactivate(_ context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return &domain.Tenant{ID: tenantID, Name: "Test", Status: "active"}, nil
}

func TestAddMemberAdminOK(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewAddMemberCommand(repo)
	ctx := context.Background()
	adminID := uuid.New()
	tenantID := uuid.New()
	targetUser := uuid.New()

	// Add admin to repo so FindUserTenant works after AddMember
	repo.memberships[tenantID.String()+":"+adminID.String()] = domain.UserTenant{
		UserID: adminID, TenantID: tenantID, Role: "administrator", Name: "Test",
	}

	principal := Principal{UserID: adminID, TenantID: tenantID, Role: "administrator"}
	member, err := cmd.Execute(ctx, principal, AddMemberInput{
		TenantID: tenantID,
		UserID:   targetUser,
		Role:     "patient",
	})
	if err != nil {
		t.Fatalf("add member: %v", err)
	}
	if member.UserID != targetUser {
		t.Errorf("user_id = %v, want %v", member.UserID, targetUser)
	}
	if member.Role != "patient" {
		t.Errorf("role = %v, want patient", member.Role)
	}
}

func TestAddMemberNonAdminFails(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewAddMemberCommand(repo)
	ctx := context.Background()
	tenantID := uuid.New()

	principal := Principal{UserID: uuid.New(), TenantID: tenantID, Role: "doctor"}

	_, err := cmd.Execute(ctx, principal, AddMemberInput{
		TenantID: tenantID,
		UserID:   uuid.New(),
		Role:     "patient",
	})
	if err != ErrNotAdmin {
		t.Errorf("expected ErrNotAdmin, got %v", err)
	}
}

func TestAddMemberInvalidRole(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewAddMemberCommand(repo)
	ctx := context.Background()
	tenantID := uuid.New()
	adminID := uuid.New()

	repo.memberships[tenantID.String()+":"+adminID.String()] = domain.UserTenant{
		UserID: adminID, TenantID: tenantID, Role: "administrator", Name: "Test",
	}
	principal := Principal{UserID: adminID, TenantID: tenantID, Role: "administrator"}

	_, err := cmd.Execute(ctx, principal, AddMemberInput{
		TenantID: tenantID,
		UserID:   uuid.New(),
		Role:     "invalid",
	})
	if err != ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}

func TestRemoveMemberAdminOK(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewRemoveMemberCommand(repo)
	ctx := context.Background()
	adminID := uuid.New()
	tenantID := uuid.New()
	targetUser := uuid.New()

	repo.memberships[tenantID.String()+":"+adminID.String()] = domain.UserTenant{
		UserID: adminID, TenantID: tenantID, Role: "administrator", Name: "Test",
	}
	repo.memberships[tenantID.String()+":"+targetUser.String()] = domain.UserTenant{
		UserID: targetUser, TenantID: tenantID, Role: "patient", Name: "Test",
	}

	principal := Principal{UserID: adminID, TenantID: tenantID, Role: "administrator"}

	err := cmd.Execute(ctx, principal, tenantID, targetUser)
	if err != nil {
		t.Fatalf("remove member: %v", err)
	}

	_, err = repo.FindUserTenant(targetUser, tenantID)
	if err == nil {
		t.Error("expected user to be removed")
	}
}

func TestRemoveMemberNonAdminFails(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewRemoveMemberCommand(repo)
	ctx := context.Background()
	tenantID := uuid.New()

	principal := Principal{UserID: uuid.New(), TenantID: tenantID, Role: "patient"}

	err := cmd.Execute(ctx, principal, tenantID, uuid.New())
	if err != ErrNotAdmin {
		t.Errorf("expected ErrNotAdmin, got %v", err)
	}
}

func TestListMembersIsMemberOK(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewListMembersQuery(repo)
	ctx := context.Background()
	tenantID := uuid.New()
	userID := uuid.New()

	repo.memberships[tenantID.String()+":"+userID.String()] = domain.UserTenant{
		UserID: userID, TenantID: tenantID, Role: "administrator", Name: "Test",
	}
	repo.byTenant[tenantID] = []domain.UserTenant{
		{UserID: userID, TenantID: tenantID, Role: "administrator", Name: "Test"},
	}

	principal := Principal{UserID: userID, TenantID: tenantID, Role: "administrator"}

	members, err := cmd.Execute(ctx, principal, tenantID)
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("members = %d, want 1", len(members))
	}
	if members[0].Role != "administrator" {
		t.Errorf("role = %v, want administrator", members[0].Role)
	}
}

func TestListMembersNotMemberFails(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewListMembersQuery(repo)
	ctx := context.Background()
	tenantID := uuid.New()

	principal := Principal{UserID: uuid.New(), TenantID: tenantID, Role: "administrator"}

	_, err := cmd.Execute(ctx, principal, tenantID)
	if err != ErrUserNotMember {
		t.Errorf("expected ErrUserNotMember, got %v", err)
	}
}

func TestListMembersNonAdminFails(t *testing.T) {
	repo := newMockTenantMemberRepo()
	cmd := NewListMembersQuery(repo)
	ctx := context.Background()
	tenantID := uuid.New()

	principal := Principal{UserID: uuid.New(), TenantID: tenantID, Role: "doctor"}

	_, err := cmd.Execute(ctx, principal, tenantID)
	if err != ErrNotAdmin {
		t.Errorf("expected ErrNotAdmin, got %v", err)
	}
}
