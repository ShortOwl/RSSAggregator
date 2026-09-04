package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ShortOwl/RSSAggregator/internal/auth"
	"github.com/ShortOwl/RSSAggregator/internal/database"
	"github.com/ShortOwl/RSSAggregator/internal/response"
	"github.com/lib/pq"
)

type authResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

// HandleRegister creates a password-protected user and returns an access token.
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Decode json into Local parameters struct.
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't decode request body")
		return
	}

	params.Name = strings.TrimSpace(params.Name)
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))

	if err := validateRegistrationInput(params.Name, params.Email, params.Password); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Hash the password.
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't secure the password")
		return
	}
	// save the user to DB.
	user, err := h.DB.RegisterUser(r.Context(), database.RegisterUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: passwordHash,
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
			response.WithError(w, http.StatusConflict, "Email is already registered")
			return
		}
		response.WithError(w, http.StatusInternalServerError, "Couldn't register user")
		return
	}

	token, err := auth.GenerateJWT(user.ID, h.JWTSecret)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't generate JWT token")
		return
	}

	response.WithJSON(w, http.StatusOK, authResponse{
		Token:     token,
		TokenType: "Bearer",
	})

}

// HandleLogin verifies an email and password and returns a new access token.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// login ma apde email id ane password apiye.
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't decode request body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(params.Email))
	user, err := h.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.WithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		response.WithError(w, http.StatusInternalServerError, "Couldn't log in")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		response.WithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	// Generate JWT token.
	token, err := auth.GenerateJWT(user.ID, h.JWTSecret)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't generate JWT token")
		return
	}
	response.WithJSON(w, http.StatusOK, authResponse{
		Token:     token,
		TokenType: "Bearer",
	})
	// returns a fresh 24 HOUR JWT token.

}
func validateRegistrationInput(name, email, password string) error {
	if name == "" {
		return errors.New("Name is required")
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return errors.New("A valid email is required")
	}

	if len(password) < 8 {
		return errors.New("Password must be at least 8 bytes")
	}
	if len(password) > 72 {
		return errors.New("Password must be at most 72 bytes")
	}
	return nil
}
