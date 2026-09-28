package handler

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
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

	return validatePassword(password)
}
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be at least 8 bytes")
	}
	if len(password) > 72 {
		return errors.New("Password must be at most 72 bytes")
	}
	return nil
}

func (h *Handler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't decode request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(params.Email))

	user, err := h.DB.GetUserByEmail(r.Context(), email)
	if err == nil {
		otp, err := auth.GenerateOTP()
		if err == nil {
			now := time.Now().UTC()
			err = h.DB.SavePasswordReset(r.Context(), database.SavePasswordResetParams{
				UserID:    user.ID,
				OtpHash:   auth.HashOTP(otp),
				CreatedAt: now,
				ExpiresAt: now.Add(10 * time.Minute),
			})
			if err == nil {
				err = h.sendPasswordResetEmail(email, otp)
			}
		}
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("password recovery request failed", "error", err)
	}
	response.WithJSON(w, http.StatusOK, map[string]string{
		"message": "If that email is registered, a recovery code will be sent.",
	})

}

func (h *Handler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Email    string `json:"email"`
		OTP      string `json:"otp"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		response.WithError(w, http.StatusBadRequest, "Couldn't decode request body")
		return
	}
	// validation steps.
	if err := validatePassword(params.Password); err != nil {
		response.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(params.OTP) != 6 || strings.Trim(params.OTP, "0123456789") != "" {
		response.WithError(w, http.StatusBadRequest, "Invalid or expired code")
		return
	}

	email := strings.ToLower(strings.TrimSpace(params.Email))
	user, err := h.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("password reset lookup failed", "error", err)
		}
		response.WithError(w, http.StatusBadRequest, "Invalid or expired code")
		return
	}

	// user avi jai to verify otp , it successfull update password.
	reset, err := h.DB.GetPasswordResetByUser(r.Context(), user.ID)
	now := time.Now().UTC()
	otpHash := auth.HashOTP(params.OTP)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("password reset record lookup failed", "error", err)
	}
	if err != nil || !now.Before(reset.ExpiresAt) ||
		subtle.ConstantTimeCompare([]byte(otpHash), []byte(reset.OtpHash)) != 1 {
		response.WithError(w, http.StatusBadRequest, "Invalid or expired code")
		return
	}
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "Couldn't secure the password")
		return
	}
	rows, err := h.DB.CompletePasswordReset(r.Context(), database.CompletePasswordResetParams{
		UserID: user.ID, OtpHash: otpHash, ExpiresAt: now,
		PasswordHash: passwordHash, UpdatedAt: now,
	})
	if err != nil {
		slog.Error("password reset failed", "error", err)
		response.WithError(w, http.StatusInternalServerError, "Couldn't reset password")
		return
	}
	if rows != 1 {
		response.WithError(w, http.StatusBadRequest, "Invalid or expired code")
		return
	}
	response.WithJSON(w, http.StatusOK, map[string]string{"message": "Password reset. Please sign in."})
}
