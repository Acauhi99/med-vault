package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id,omitempty"`
	Role     string    `json:"role,omitempty"`
	Expiry   int64     `json:"exp"`
	Type     string    `json:"type"` // "temp", "access", or "refresh"
}

type Generator struct {
	secret []byte
}

func NewGenerator(secret string) *Generator {
	return &Generator{secret: []byte(secret)}
}

func (g *Generator) GenerateAccessToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return g.generate(Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Expiry:   time.Now().Add(ttl).Unix(),
		Type:     "access",
	})
}

func (g *Generator) GenerateRefreshToken(userID, tenantID uuid.UUID, role string, ttl time.Duration) (string, error) {
	return g.generate(Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Expiry:   time.Now().Add(ttl).Unix(),
		Type:     "refresh",
	})
}

func (g *Generator) GenerateTemporaryToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	return g.generate(Claims{
		UserID: userID,
		Expiry: time.Now().Add(ttl).Unix(),
		Type:   "temp",
	})
}

func (g *Generator) generate(claims Claims) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	sig := g.sign(payload)
	encoded := base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)

	return encoded, nil
}

func (g *Generator) sign(data []byte) []byte {
	mac := hmac.New(sha256.New, g.secret)
	mac.Write(data)
	return mac.Sum(nil)
}
