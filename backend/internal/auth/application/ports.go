package application

import (
	"time"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, password string) error
}

type JWTGenerator interface {
	GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error)
	GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error)
	GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error)
	Verify(rawToken string) (JWTClaims, error)
}

type JWTClaims struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Role     string
	Type     string
}
