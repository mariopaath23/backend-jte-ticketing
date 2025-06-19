package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
)

// UserClaimsKey is the type for our context key. Using a custom type
// prevents collisions with other context keys.
type UserClaimsKey string

// ClaimsKey is the exported constant that we will use across our application
// to access the user claims in the context.
const ClaimsKey UserClaimsKey = "userClaims"

// Auth is a middleware that checks for a valid JWT from either a cookie or Authorization header.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""

		// 1. Try to get the token from the HttpOnly cookie first.
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}

		// 2. If no cookie, try to get it from the Authorization header.
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Auth Error: No token found in cookie or Authorization header")
				http.Error(w, "Missing authorization token", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Println("Auth Error: Invalid Authorization header format")
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString = parts[1]
		}

		// 3. Validate the token we found.
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("Auth Error: Token validation failed. Error: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. If the token is valid, add claims to the request context using our exported key.
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)

		// 5. Call the next handler in the chain.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
