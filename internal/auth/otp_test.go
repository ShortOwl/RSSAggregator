package auth

import (
	"strings"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	for i := 0; i < 100; i++ {
		otp, err := GenerateOTP()
		if err != nil {
			t.Fatal(err)
		}
		if len(otp) != 6 || strings.Trim(otp, "0123456789") != "" {
			t.Fatalf("OTP must be exactly six digits, got %q", otp)
		}
	}
}

func TestHashOTP(t *testing.T) {
	const want = "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	if got := HashOTP("123456"); got != want {
		t.Fatalf("SHA-256 hash = %q, want %q", got, want)
	}
}
