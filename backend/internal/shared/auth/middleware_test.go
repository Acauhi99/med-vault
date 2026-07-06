package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/application"
	"github.com/google/uuid"
)

func TestTenantMiddlewareRejectsMissingBearerToken(t *testing.T) {
	handler := TenantMiddleware(nil, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestTenantMiddlewareAddsPrincipalToContext(t *testing.T) {
	expected := Principal{UserID: uuid.New(), TenantID: uuid.New(), Role: RolePatient}

	gen := &stubJWTGen{
		claims: application.JWTClaims{
			UserID:   expected.UserID,
			TenantID: expected.TenantID,
			Role:     string(expected.Role),
			Type:     "access",
		},
	}

	handler := TenantMiddleware(gen, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actual, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("principal missing from context")
		}
		if actual != expected {
			t.Fatalf("expected principal %#v, got %#v", expected, actual)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}

func TestTenantMiddlewareRejectsTempToken(t *testing.T) {
	gen := &stubJWTGen{
		claims: application.JWTClaims{
			UserID: uuid.New(),
			Type:   "temp",
		},
	}

	handler := TenantMiddleware(gen, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/cases", nil)
	request.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestTenantMiddlewareRejectsInvalidToken(t *testing.T) {
	gen := &stubJWTGen{err: errors.New("invalid")}

	handler := TenantMiddleware(gen, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

type stubJWTGen struct {
	claims application.JWTClaims
	err    error
}

func (s *stubJWTGen) GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "", nil
}

func (s *stubJWTGen) GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	return "", nil
}

func (s *stubJWTGen) GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return "", nil
}

func (s *stubJWTGen) Verify(rawToken string) (application.JWTClaims, error) {
	return s.claims, s.err
}
