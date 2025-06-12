package middleware

import (
	"net/http"
	"strings"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
)

// Auth is a middleware that checks for a valid JWT.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We can get the token from the Authorization header
		// The header will be in the format `Bearer <token>`
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Alternatively, try to get it from a cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Missing authorization token", http.StatusUnauthorized)
				return
			}
			authHeader = "Bearer " + cookie.Value
		}

		// Split the header to get the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate the token
		_, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}
