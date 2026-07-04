package auth

import (
	"net/http"

	"github.com/Acauhi99/med-vault/internal/shared/httpx"
)

func RequireRole(allowed ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
				return
			}

			if principal.TenantID == (Principal{}).TenantID && principal.Role == "" {
				httpx.WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "tenant selection required")
				return
			}

			for _, role := range allowed {
				if principal.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			httpx.WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
		})
	}
}
