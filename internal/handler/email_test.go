package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSendBrevoEmail(t *testing.T) {
	for _, tt := range []struct {
		name       string
		statusCode int
		wantError  bool
	}{
		{"accepted", http.StatusCreated, false},
		{"provider rejected", http.StatusUnauthorized, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != brevoEmailURL {
					t.Errorf("unexpected email provider URL: %s", r.URL)
				}
				if r.Method != http.MethodPost || r.Header.Get("api-key") != "test-provider-key" {
					t.Errorf("unexpected request method or authorization")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing JSON content type")
				}
				var body struct {
					Sender struct {
						Email string `json:"email"`
						Name  string `json:"name"`
					} `json:"sender"`
					To []struct {
						Email string `json:"email"`
					} `json:"to"`
					Subject     string `json:"subject"`
					TextContent string `json:"textContent"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Sender.Email != "sender@example.com" || body.Sender.Name != "Team Margin" ||
					len(body.To) != 1 || body.To[0].Email != "reader@example.com" ||
					body.Subject != "Margin password reset code" || !strings.Contains(body.TextContent, "482913") {
					t.Errorf("unexpected email payload: %+v", body)
				}
				return &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(strings.NewReader(`{"messageId":"test-email"}`)),
					Header:     make(http.Header),
				}, nil
			})}

			err := sendBrevoEmail(context.Background(), client, brevoEmailURL,
				"test-provider-key", brevoMessage{
					Sender:      brevoAddress{Email: "sender@example.com", Name: "Team Margin"},
					To:          []brevoAddress{{Email: "reader@example.com"}},
					Subject:     "Margin password reset code",
					TextContent: "Your password reset code is 482913.\nIt expires in 10 minutes.\n",
				})
			if (err != nil) != tt.wantError {
				t.Fatalf("sendBrevoEmail error = %v, want error = %v", err, tt.wantError)
			}
			if err != nil && strings.Contains(err.Error(), "test-provider-key") {
				t.Error("provider error exposed API key")
			}
		})
	}
}
