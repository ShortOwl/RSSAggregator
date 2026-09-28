package handler

import (
	"errors"
	"net/smtp"
	"strings"
)

// this file will send an email with otp to the user.

func (h *Handler) sendPasswordResetEmail(to, otp string) error {
	if strings.ContainsAny(to+h.SMTPEmail, "\r\n") {
		return errors.New("invalid email address")
	}

	message := []byte(
		"From: " + h.SMTPEmail + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: Margin password reset code\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			"Your password reset code is " + otp + ".\r\n" +
			"It expires in 10 minutes.\r\n",
	)

	server := "smtp.gmail.com" // Gmail's outgoing mail server.
	smtpAuth := smtp.PlainAuth("", h.SMTPEmail, h.SMTPPassword, server)
	return smtp.SendMail(server+":587", smtpAuth, h.SMTPEmail, []string{to}, message)

}
