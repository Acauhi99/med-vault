package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenCommand struct {
	jwtGen     JWTGenerator
	accessTTL  time.Duration
	refreshTTL time.Duration
	tokenStore *TokenStore
}

func NewRefreshTokenCommand(
	jwtGen JWTGenerator,
	accessTTL, refreshTTL time.Duration,
	tokenStore *TokenStore,
) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		jwtGen:     jwtGen,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		tokenStore: tokenStore,
	}
}

func (c *RefreshTokenCommand) Execute(input RefreshTokenInput) (RefreshTokenOutput, error) {
	claims, err := c.jwtGen.Verify(input.RefreshToken)
	if err != nil {
		return RefreshTokenOutput{}, ErrInvalidRefreshToken
	}

	if claims.Type != "refresh" {
		return RefreshTokenOutput{}, ErrInvalidRefreshToken
	}

	if claims.TenantID == (uuid.UUID{}) {
		return RefreshTokenOutput{}, ErrInvalidRefreshToken
	}

	if c.tokenStore != nil && c.tokenStore.IsRevoked(input.RefreshToken) {
		return RefreshTokenOutput{}, ErrInvalidRefreshToken
	}

	accessToken, err := c.jwtGen.GenerateAccessToken(claims.UserID, claims.TenantID, claims.Role, c.accessTTL)
	if err != nil {
		return RefreshTokenOutput{}, err
	}

	refreshToken, err := c.jwtGen.GenerateRefreshToken(claims.UserID, claims.TenantID, claims.Role, c.refreshTTL)
	if err != nil {
		return RefreshTokenOutput{}, err
	}

	return RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(c.accessTTL.Seconds()),
	}, nil
}

var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
