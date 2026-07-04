package application

import (
	"testing"
	"time"
)

func TestTokenStoreRevokeAndCheck(t *testing.T) {
	store := NewTokenStore(time.Hour)

	if store.IsRevoked("token-abc") {
		t.Error("token should not be revoked before revocation")
	}

	store.Revoke("token-abc")

	if !store.IsRevoked("token-abc") {
		t.Error("token should be revoked after Revoke()")
	}

	store.Revoke("token-abc")
	if !store.IsRevoked("token-abc") {
		t.Error("double-revoke should still be revoked")
	}
}

func TestTokenStoreExpiry(t *testing.T) {
	store := NewTokenStore(0) // immediate expiry

	store.Revoke("expired-token")

	// On next check, the token is expired → not revoked
	if store.IsRevoked("expired-token") {
		t.Error("expired revocation entry should be cleaned up")
	}
}

func TestRefreshTokenRevocation(t *testing.T) {
	tokenStore := NewTokenStore(time.Hour)
	cmd := NewRefreshTokenCommand(&refreshMockJWTGen{}, 15*time.Minute, 168*time.Hour, tokenStore)

	tokenStore.Revoke("revoked-token")

	_, err := cmd.Execute(RefreshTokenInput{RefreshToken: "revoked-token"})
	if err != ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken for revoked token, got %v", err)
	}
}
