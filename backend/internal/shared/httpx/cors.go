package httpx

import (
	"net/http"
	"strings"
)

func CORSMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, origin := range strings.Split(allowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-Id")
					w.Header().Set("Access-Control-Max-Age", "600")
					w.Header().Add("Vary", "Origin")
				} else if r.Method == http.MethodOptions {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
