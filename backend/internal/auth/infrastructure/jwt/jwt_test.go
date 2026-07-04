package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndVerifyAccessToken(t *testing.T) {
	gen := NewGenerator("test-secret")
	userID := uuid.New()
	tenantID := uuid.New()

	token, err := gen.GenerateAccessToken(userID, tenantID, "doctor", 15*time.Minute)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}

	claims, err := gen.Verify(token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.TenantID != tenantID {
		t.Errorf("TenantID = %v, want %v", claims.TenantID, tenantID)
	}
	if claims.Role != "doctor" {
		t.Errorf("Role = %v, want %v", claims.Role, "doctor")
	}
	if claims.Type != "access" {
		t.Errorf("Type = %v, want %v", claims.Type, "access")
	}
}

func TestGenerateAndVerifyTemporaryToken(t *testing.T) {
	gen := NewGenerator("test-secret")
	userID := uuid.New()

	token, err := gen.GenerateTemporaryToken(userID, 5*time.Minute)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	claims, err := gen.Verify(token)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.TenantID != (uuid.UUID{}) {
		t.Errorf("TenantID should be zero, got %v", claims.TenantID)
	}
	if claims.Type != "temp" {
		t.Errorf("Type = %v, want %v", claims.Type, "temp")
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	gen := NewGenerator("test-secret")
	userID := uuid.New()

	token, err := gen.GenerateAccessToken(userID, uuid.New(), "patient", -1*time.Hour)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	_, err = gen.Verify(token)
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestVerifyRejectsInvalidSignature(t *testing.T) {
	gen1 := NewGenerator("secret-1")
	gen2 := NewGenerator("secret-2")

	token, _ := gen1.GenerateAccessToken(uuid.New(), uuid.New(), "patient", 15*time.Minute)
	_, err := gen2.Verify(token)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestVerifyRejectsMalformedToken(t *testing.T) {
	gen := NewGenerator("secret")
	_, err := gen.Verify("not-a-token")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
