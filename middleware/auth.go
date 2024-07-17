package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func Auth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			secret := os.Getenv("AUTH_SECRET")

			if authHeader == "" {
				notAuthError(w)
			}

			if !strings.HasPrefix(authHeader, "Basic ") {
				notAuthError(w)
			}

			given := strings.TrimPrefix(authHeader, "Basic ")

			if subtle.ConstantTimeCompare([]byte(secret), []byte(given)) != 1 {
				notAuthError(w)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func notAuthError(w http.ResponseWriter) {
	http.Error(w, "the authentication secret is incorrect", http.StatusUnauthorized)
}
