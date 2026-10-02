package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
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
	TextContent string         `json:"textContent,omitempty"`
	HTMLContent string         `json:"htmlContent,omitempty"`
}

func (h *Handler) sendPasswordResetEmail(ctx context.Context, to, otp string) error {
	client := &http.Client{Timeout: 15 * time.Second}

	message := brevoMessage{
		Sender:      brevoAddress{Email: h.BrevoSenderEmail, Name: "Team Margin"},
		To:          []brevoAddress{{Email: to}},
		Subject:     "Margin password reset code",
		HTMLContent: passwordResetHTML(otp),
	}

	return sendBrevoEmail(ctx, client, brevoEmailURL, h.BrevoAPIKey, message)
}

func (h *Handler) sendWelcomeEmail(ctx context.Context, to string) error {
	client := &http.Client{Timeout: 15 * time.Second}

	message := brevoMessage{
		Sender:      brevoAddress{Email: h.BrevoSenderEmail, Name: "Team Margin"},
		To:          []brevoAddress{{Email: to}},
		Subject:     "Welcome to Margin",
		HTMLContent: welcomeHTML(),
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

// passwordResetHTML uses inline styles and presentation tables for email clients.
func passwordResetHTML(otp string) string {
	return `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Your Margin recovery code</title></head>
<body style="margin:0;padding:0;background-color:#f6f5f4;color:#111111;font-family:Arial,Helvetica,sans-serif;">
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Your Margin password reset code is ready. It expires in 10 minutes.</div>
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="background-color:#f6f5f4;">
<tr><td align="center" style="padding:40px 16px;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width:560px;">
<tr><td align="center" style="padding:0 0 24px;font-size:28px;font-weight:700;letter-spacing:-1px;color:#111111;">margin<span style="color:#0075de;">.</span></td></tr>
<tr><td align="center" style="padding:0 0 24px;"><h1 style="margin:0;font-size:26px;line-height:1.3;letter-spacing:-0.5px;font-weight:600;color:#111111;">Back to your reading space</h1></td></tr>
<tr><td style="padding:28px 24px;background-color:#ffffff;border:1px solid #e5e3e0;border-radius:12px;">
<p style="margin:0 0 24px;font-size:16px;line-height:1.6;color:#615d59;">Use this code in Margin to reset your password.</p>
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0"><tr><td align="center" style="padding:20px 8px;background-color:#f6f5f4;border-radius:8px;">
<span style="font-family:Consolas,'Courier New',monospace;font-size:32px;font-weight:700;line-height:1.4;letter-spacing:6px;color:#111111;">` + html.EscapeString(otp) + `</span>
</td></tr></table>
<p style="margin:24px 0 12px;font-size:15px;line-height:1.6;color:#615d59;">This code expires in <strong style="color:#111111;">10 minutes</strong> and can only be used once.</p>
<p style="margin:0 0 24px;font-size:15px;line-height:1.6;color:#615d59;">Keep this code private. Do not share it with anyone.</p>
<p style="margin:0;font-size:15px;line-height:1.6;color:#111111;">Happy reading,<br>Team Margin</p>
</td></tr>
<tr><td style="padding:24px 8px 0;font-size:13px;line-height:1.6;color:#615d59;">If you didn’t request a password reset, you can ignore this email. Your password will stay the same.</td></tr>
<tr><td align="center" style="padding:24px 8px 0;font-size:12px;line-height:1.6;color:#757575;">A reading ritual, made your own.</td></tr>
</table>
</td></tr></table>
</body>
</html>`
}

// welcomeHTML matches the password recovery email’s layout and colors.
func welcomeHTML() string {
	return `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Welcome to Margin</title></head>
<body style="margin:0;padding:0;background-color:#f6f5f4;color:#111111;font-family:Arial,Helvetica,sans-serif;">
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Your reading space is ready. Make room for curiosity.</div>
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="background-color:#f6f5f4;">
<tr><td align="center" style="padding:40px 16px;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width:560px;">
<tr><td align="center" style="padding:0 0 24px;font-size:28px;font-weight:700;letter-spacing:-1px;color:#111111;">margin<span style="color:#0075de;">.</span></td></tr>
<tr><td align="center" style="padding:0 0 24px;"><h1 style="margin:0;font-size:26px;line-height:1.3;letter-spacing:-0.5px;font-weight:600;color:#111111;">Make room for curiosity.</h1></td></tr>
<tr><td style="padding:28px 24px;background-color:#ffffff;border:1px solid #e5e3e0;border-radius:12px;">
<p style="margin:0 0 20px;font-size:16px;line-height:1.6;color:#111111;">Welcome to Margin. We’re glad you’re here.</p>
<p style="margin:0 0 24px;font-size:16px;line-height:1.6;color:#615d59;">A little less noise. A lot more perspective. Your favorite voices now have a quiet place to come together.</p>
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0"><tr><td style="padding:20px;background-color:#f6f5f4;border-radius:8px;">
<p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:#615d59;"><strong style="color:#111111;">Follow what matters.</strong><br>Add feeds from the writers and publications you enjoy.</p>
<p style="margin:0;font-size:15px;line-height:1.6;color:#615d59;"><strong style="color:#111111;">Save what stays with you.</strong><br>Bookmark articles to return to when you have a little more time.</p>
</td></tr></table>
<p style="margin:24px 0;font-size:15px;line-height:1.6;color:#615d59;">Start with one feed. Make this space your own, one good read at a time.</p>
<p style="margin:0;font-size:15px;line-height:1.6;color:#111111;">Happy reading,<br>Team Margin</p>
</td></tr>
<tr><td style="padding:24px 8px 0;font-size:13px;line-height:1.6;color:#615d59;">You’re receiving this email because this address was used to create a Margin account. If that wasn’t you, you can ignore this email.</td></tr>
<tr><td align="center" style="padding:24px 8px 0;font-size:12px;line-height:1.6;color:#757575;">A reading ritual, made your own.</td></tr>
</table>
</td></tr></table>
</body>
</html>`
}
