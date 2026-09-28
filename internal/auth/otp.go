package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1_000_000))

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", number.Int64()), nil
}

func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return fmt.Sprintf("%x", hash)
}
