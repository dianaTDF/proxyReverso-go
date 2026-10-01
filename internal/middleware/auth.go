package middleware

import (
	"net/http"
	"strings"
)

// Auth aplica seguridad L7 verificando la existencia de un Bearer token.
func Auth(validToken string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized - Missing Token", http.StatusUnauthorized)
				return
			}
			
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != validToken {
				http.Error(w, "Forbidden - Invalid Token", http.StatusForbidden)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
