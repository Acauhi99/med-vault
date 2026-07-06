package auth

import (
	"net/http"

	"github.com/Acauhi99/med-vault/internal/auth/application"
)

func TenantMiddleware(gen application.JWTGenerator, unauthorized func(http.ResponseWriter, *http.Request)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken := bearerToken(r.Header.Get("Authorization"))
			if rawToken == "" || gen == nil {
				unauthorized(w, r)
				return
			}

			claims, err := gen.Verify(rawToken)
			if err != nil {
				unauthorized(w, r)
				return
			}

			if claims.Type != "access" {
				unauthorized(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), Principal{
				UserID:   claims.UserID,
				TenantID: claims.TenantID,
				Role:     Role(claims.Role),
			})))
		})
	}
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return ""
	}
	return header[len(prefix):]
}
