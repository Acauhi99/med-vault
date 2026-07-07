package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/domain"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Create(user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

type mockTenantRepo struct {
	memberships map[uuid.UUID][]domain.UserTenant
}

func (m *mockTenantRepo) FindUserTenants(userID uuid.UUID) ([]domain.UserTenant, error) {
	return m.memberships[userID], nil
}

func (m *mockTenantRepo) FindUserTenant(userID, tenantID uuid.UUID) (*domain.UserTenant, error) {
	for _, ut := range m.memberships[userID] {
		if ut.TenantID == tenantID {
			return &ut, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockTenantRepo) FindByName(name string) (*domain.Tenant, error) {
	return &domain.Tenant{Name: name, Status: "active"}, nil
}

func (m *mockTenantRepo) AddMember(ctx context.Context, tenantID, userID uuid.UUID, role string) error {
	m.memberships[userID] = append(m.memberships[userID], domain.UserTenant{
		UserID: userID, TenantID: tenantID, Role: role, Name: "Test",
	})
	return nil
}

func (m *mockTenantRepo) RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error {
	for i, ut := range m.memberships[userID] {
		if ut.TenantID == tenantID {
			m.memberships[userID] = append(m.memberships[userID][:i], m.memberships[userID][i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockTenantRepo) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]domain.UserTenant, error) {
	var result []domain.UserTenant
	for _, uts := range m.memberships {
		for _, ut := range uts {
			if ut.TenantID == tenantID {
				result = append(result, ut)
			}
		}
	}
	return result, nil
}

func (m *mockTenantRepo) Reactivate(_ context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return &domain.Tenant{ID: tenantID, Name: "Test", Status: "active"}, nil
}

func (m *mockTenantRepo) Create(_ context.Context, name string) (*domain.Tenant, error) {
	return &domain.Tenant{Name: name, Status: "active"}, nil
}

func (m *mockTenantRepo) Suspend(_ context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return &domain.Tenant{ID: tenantID, Name: "Test", Status: "suspended"}, nil
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *mockHasher) Compare(hashed, password string) error {
	if hashed != "hashed-"+password {
		return errors.New("wrong password")
	}
	return nil
}

type mockJWTGen struct{}

func (m *mockJWTGen) GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "token-" + userID.String()[:8], nil
}

func (m *mockJWTGen) GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "refresh-" + userID.String()[:8], nil
}

func (m *mockJWTGen) GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	return "temp-" + userID.String()[:8], nil
}

func (m *mockJWTGen) Verify(rawToken string) (JWTClaims, error) {
	return JWTClaims{}, nil
}

func TestRegisterSuccess(t *testing.T) {
	users := &mockUserRepo{users: make(map[string]*domain.User)}
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewRegisterCommand(users, tenants, &mockHasher{})

	out, err := cmd.Execute(RegisterInput{Email: "test@example.com", Password: "Password12345!"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if out.Email != "test@example.com" {
		t.Errorf("email = %v, want test@example.com", out.Email)
	}
	if out.Status != "active" {
		t.Errorf("status = %v, want active", out.Status)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	users := &mockUserRepo{users: map[string]*domain.User{
		"existing@example.com": {ID: uuid.New(), Email: "existing@example.com"},
	}}
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewRegisterCommand(users, tenants, &mockHasher{})

	_, err := cmd.Execute(RegisterInput{Email: "existing@example.com", Password: "Password12345!"})
	if err != ErrEmailAlreadyExists {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestRegisterWeakPasswordRejected(t *testing.T) {
	users := &mockUserRepo{users: make(map[string]*domain.User)}
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewRegisterCommand(users, tenants, &mockHasher{})

	_, err := cmd.Execute(RegisterInput{Email: "test@example.com", Password: "weakpass"})
	if err != ErrWeakPassword {
		t.Errorf("expected ErrWeakPassword, got %v", err)
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	userID := uuid.New()
	tenantID := uuid.New()
	users := &mockUserRepo{users: map[string]*domain.User{
		"test@example.com": {ID: userID, Email: "test@example.com", PasswordHash: "hashed-password12345"},
	}}
	tenants := &mockTenantRepo{memberships: map[uuid.UUID][]domain.UserTenant{
		userID: {{UserID: userID, TenantID: tenantID, Role: "patient", Name: "Test Hospital"}},
	}}
	cmd := NewAuthenticateCommand(users, tenants, &mockHasher{}, &mockJWTGen{}, 5*time.Minute)

	out, err := cmd.Execute(AuthenticateInput{Email: "test@example.com", Password: "password12345"})
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}
	if out.AccessToken == "" {
		t.Error("empty access token")
	}
	if len(out.Tenants) != 1 {
		t.Errorf("tenants = %d, want 1", len(out.Tenants))
	}
}

func TestAuthenticateWrongPassword(t *testing.T) {
	users := &mockUserRepo{users: map[string]*domain.User{
		"test@example.com": {ID: uuid.New(), Email: "test@example.com", PasswordHash: "hashed-password12345"},
	}}
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewAuthenticateCommand(users, tenants, &mockHasher{}, &mockJWTGen{}, 5*time.Minute)

	_, err := cmd.Execute(AuthenticateInput{Email: "test@example.com", Password: "wrong"})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticateUserNotFound(t *testing.T) {
	users := &mockUserRepo{users: make(map[string]*domain.User)}
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewAuthenticateCommand(users, tenants, &mockHasher{}, &mockJWTGen{}, 5*time.Minute)

	_, err := cmd.Execute(AuthenticateInput{Email: "nobody@example.com", Password: "password12345"})
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSelectTenantSuccess(t *testing.T) {
	userID := uuid.New()
	tenantID := uuid.New()
	tenants := &mockTenantRepo{memberships: map[uuid.UUID][]domain.UserTenant{
		userID: {{UserID: userID, TenantID: tenantID, Role: "doctor", Name: "Test Hospital"}},
	}}
	cmd := NewSelectTenantCommand(tenants, &mockJWTGen{}, 15*time.Minute, 168*time.Hour)

	out, err := cmd.Execute(SelectTenantInput{UserID: userID, TenantID: tenantID})
	if err != nil {
		t.Fatalf("select tenant: %v", err)
	}
	if out.AccessToken == "" {
		t.Error("empty access token")
	}
	if out.RefreshToken == "" {
		t.Error("empty refresh token")
	}
	if out.ExpiresIn != 900 {
		t.Errorf("expires_in = %d, want 900", out.ExpiresIn)
	}
}

func TestSelectTenantNotMember(t *testing.T) {
	userID := uuid.New()
	tenants := &mockTenantRepo{memberships: make(map[uuid.UUID][]domain.UserTenant)}
	cmd := NewSelectTenantCommand(tenants, &mockJWTGen{}, 15*time.Minute, 168*time.Hour)

	_, err := cmd.Execute(SelectTenantInput{UserID: userID, TenantID: uuid.New()})
	if err != ErrUserNotMember {
		t.Errorf("expected ErrUserNotMember, got %v", err)
	}
}
