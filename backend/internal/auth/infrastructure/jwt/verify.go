package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/Acauhi99/med-vault/internal/auth/application"
)

func (g *Generator) Verify(rawToken string) (application.JWTClaims, error) {
	claims, err := g.verify(rawToken)
	if err != nil {
		return application.JWTClaims{}, err
	}
	return application.JWTClaims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Role:     claims.Role,
		Type:     claims.Type,
	}, nil
}

func (g *Generator) verify(rawToken string) (Claims, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 2 {
		return Claims{}, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	expectedSig := g.sign(payload)
	if !hmac.Equal(sig, expectedSig) {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if time.Now().Unix() > claims.Expiry {
		return Claims{}, ErrTokenExpired
	}

	return claims, nil
}
