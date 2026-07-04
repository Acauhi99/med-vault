package application

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type refreshMockJWTGen struct{}

func (m *refreshMockJWTGen) GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "access-" + userID.String()[:8], nil
}

func (m *refreshMockJWTGen) GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "refresh-" + userID.String()[:8], nil
}

func (m *refreshMockJWTGen) GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	return "temp-" + userID.String()[:8], nil
}

func (m *refreshMockJWTGen) Verify(rawToken string) (JWTClaims, error) {
	if rawToken == "bad-token" {
		return JWTClaims{}, ErrInvalidRefreshToken
	}
	if rawToken == "temp-token" {
		return JWTClaims{Type: "temp", UserID: uuid.New(), TenantID: uuid.New(), Role: "patient"}, nil
	}
	return JWTClaims{
		Type:     "refresh",
		UserID:   uuid.New(),
		TenantID: uuid.New(),
		Role:     "doctor",
	}, nil
}

func TestRefreshTokenSuccess(t *testing.T) {
	cmd := NewRefreshTokenCommand(&refreshMockJWTGen{}, 15*time.Minute, 168*time.Hour)

	out, err := cmd.Execute(RefreshTokenInput{RefreshToken: "valid-refresh-token"})
	if err != nil {
		t.Fatalf("refresh: %v", err)
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

func TestRefreshTokenInvalid(t *testing.T) {
	cmd := NewRefreshTokenCommand(&refreshMockJWTGen{}, 15*time.Minute, 168*time.Hour)

	_, err := cmd.Execute(RefreshTokenInput{RefreshToken: "bad-token"})
	if err != ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestRefreshTokenWrongType(t *testing.T) {
	cmd := NewRefreshTokenCommand(&refreshMockJWTGen{}, 15*time.Minute, 168*time.Hour)

	_, err := cmd.Execute(RefreshTokenInput{RefreshToken: "temp-token"})
	if err != ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken for non-refresh token, got %v", err)
	}
}
