package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenIssuer   = "rss-aggregator"
	TokenDuration = 24 * time.Hour
)

//  These are constants shared by JWT generation and validation.
// - The issuer identifies this application.
// - Tokens expire after 24 hours.

// HashPassword turns a plain-text password into a bcrypt hash suitable for storage.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPasswordHash verifies that password produced the stored bcrypt hash.
func CheckPasswordHash(password, passwordHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return fmt.Errorf("check password hash: %w", err)
	}
	return nil
}

// GenerateJWT creates a signed access token whose subject is the user's ID.
func GenerateJWT(userID uuid.UUID, secret string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("JWT secret is empty")
	}

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(TokenDuration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	return signedToken, nil
	// note secret is not passed in the JWT, it remains with the server.
}

// ValidateJWT verifies a token and returns the user ID stored in its subject.
func ValidateJWT(tokenString, secret string) (uuid.UUID, error) {
	if strings.TrimSpace(secret) == "" {
		return uuid.Nil, errors.New("JWT secret is empty")

	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			return []byte(secret), nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("validate JWT: %w", err)
	}
	if !token.Valid {
		return uuid.Nil, errors.New("JWT is invalid")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse JWT subject: %w", err)
	}
	return userID, nil
}

// GetAPIKey extracts an API key from an Authorization header.

// Example: ApiKey <key>
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Fields(val)

	if len(vals) != 2 {
		return "", errors.New("authentication header malformed")
	}

	if !strings.EqualFold(vals[0], "ApiKey") {
		return "", errors.New("authentication header malformed")
	}

	return vals[1], nil
}
