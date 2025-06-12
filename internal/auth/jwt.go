package auth

import (
	"time"

	"[github.com/your-username/your-repo-name/backend/internal/config](https://github.com/your-username/your-repo-name/backend/internal/config)"
	"github.com/golang-jwt/jwt/v4"
)

// Claims struct will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a given email.
func GenerateJWT(email string) (string, error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}

	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token using the HS256 algorithm and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT checks if the token is valid.
func ValidateJWT(tokenString string) (*Claims, error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in []byte format, which is required by my signing method
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
