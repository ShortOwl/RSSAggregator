package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const brevoEmailURL = "https://api.brevo.com/v3/smtp/email"

type brevoAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoMessage struct {
	Sender      brevoAddress   `json:"sender"`
	To          []brevoAddress `json:"to"`
	Subject     string         `json:"subject"`
	TextContent string         `json:"textContent"`
}

func (h *Handler) sendPasswordResetEmail(ctx context.Context, to, otp string) error {
	client := &http.Client{Timeout: 15 * time.Second}

	message := brevoMessage{
		Sender:      brevoAddress{Email: h.BrevoSenderEmail, Name: "Team Margin"},
		To:          []brevoAddress{{Email: to}},
		Subject:     "Margin password reset code",
		TextContent: "Your password reset code is " + otp + ".\nIt expires in 10 minutes.\n",
	}

	return sendBrevoEmail(ctx, client, brevoEmailURL, h.BrevoAPIKey, message)
}

func (h *Handler) sendWelcomeEmail(ctx context.Context, to string) error {
	client := &http.Client{Timeout: 15 * time.Second}

	message := brevoMessage{
		Sender:      brevoAddress{Email: h.BrevoSenderEmail, Name: "Team Margin"},
		To:          []brevoAddress{{Email: to}},
		Subject:     "Welcome to Margin",
		TextContent: "Welcome to Margin.\n\nFollow the voices you care about, save what stays with you, and make room for curiosity.\n\nHappy reading,\nTeam Margin\n",
	}

	return sendBrevoEmail(ctx, client, brevoEmailURL, h.BrevoAPIKey, message)

}
func sendBrevoEmail(ctx context.Context, client *http.Client, endpoint, apikey string, message brevoMessage) error {

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Set("api-key", apikey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("email provider returned HTTP %d", resp.StatusCode)
	}

	return nil

}
