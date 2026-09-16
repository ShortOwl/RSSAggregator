package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	password := "my-test-password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if hash == password {
		t.Fatal("hash must not contain the plain-text password")
	}
	for _, tt := range []struct {
		name, password, hash string
		wantErr              bool
	}{
		{"correct password", password, hash, false},
		{"wrong password", "wrong-password", hash, true},
		{"malformed hash", password, "not-a-bcrypt-hash", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswordHash(tt.password, tt.hash); (err != nil) != tt.wantErr {
				t.Fatalf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPasswordLength(t *testing.T) {
	for _, size := range []int{72, 73} {
		_, err := HashPassword(strings.Repeat("a", size))
		if (err != nil) != (size > 72) {
			t.Errorf("HashPassword(%d bytes) error = %v", size, err)
		}
	}
}

func TestGenerateJWT(t *testing.T) {
	id := uuid.New()
	secret := "test-secret"
	before := time.Now().UTC().Truncate(time.Second)
	tokenString, err := GenerateJWT(id, secret)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ValidateJWT(tokenString, secret)
	if err != nil || got != id {
		t.Fatalf("ValidateJWT() = %v, %v; want %v", got, err, id)
	}
	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		t.Fatal(err)
	}
	if claims.Issuer != TokenIssuer || claims.Subject != id.String() {
		t.Fatalf("unexpected issuer or subject: %+v", claims)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatal("token must include issued-at and expiration")
	}
	if claims.IssuedAt.Time.Before(before) || claims.IssuedAt.Time.After(time.Now()) {
		t.Errorf("unexpected issued-at: %v", claims.IssuedAt)
	}
	if got := claims.ExpiresAt.Sub(claims.IssuedAt.Time); got != 24*time.Hour {
		t.Errorf("token lifetime = %v, want 24h", got)
	}
}

func TestValidateJWTRejectsInvalidTokens(t *testing.T) {
	for _, name := range []string{
		"expired", "wrong secret", "wrong issuer", "missing expiration",
		"future issued-at", "invalid subject", "wrong signing method", "malformed",
	} {
		t.Run(name, func(t *testing.T) {
			claims := jwt.RegisteredClaims{
				Issuer: TokenIssuer, Subject: uuid.NewString(),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			}
			secret := "test-secret"
			method := jwt.SigningMethodHS256
			switch name {
			case "expired":
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
			case "wrong secret":
				secret = "another-secret"
			case "wrong issuer":
				claims.Issuer = "another-app"
			case "missing expiration":
				claims.ExpiresAt = nil
			case "future issued-at":
				claims.IssuedAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
			case "invalid subject":
				claims.Subject = "not-a-uuid"
			case "wrong signing method":
				method = jwt.SigningMethodHS384
			}
			token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
			if err != nil {
				t.Fatal(err)
			}
			if name == "malformed" {
				token = "not-a-jwt"
			}
			if id, err := ValidateJWT(token, "test-secret"); err == nil || id != uuid.Nil {
				t.Fatalf("ValidateJWT() = %v, %v; want nil UUID and error", id, err)
			}
		})
	}
}

func TestJWTRejectsEmptySecret(t *testing.T) {
	token, err := GenerateJWT(uuid.New(), "test-secret")
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"", " \t "} {
		if _, err := GenerateJWT(uuid.New(), secret); err == nil {
			t.Error("GenerateJWT accepted an empty secret")
		}
		if _, err := ValidateJWT(token, secret); err == nil {
			t.Error("ValidateJWT accepted an empty secret")
		}
	}
}

func TestGetAPIKey(t *testing.T) {
	for _, tt := range []struct {
		header, want string
	}{
		{"ApiKey example-key", "example-key"},
		{"  apikey   example-key  ", "example-key"},
		{"", ""}, {"ApiKey", ""}, {"Bearer token", ""}, {"ApiKey one two", ""},
	} {
		t.Run(tt.header, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tt.header)
			got, err := GetAPIKey(headers)
			if got != tt.want || (err != nil) != (tt.want == "") {
				t.Fatalf("GetAPIKey() = %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
