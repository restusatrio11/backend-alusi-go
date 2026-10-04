package crypto_test

import (
	"testing"

	"backend-alusi-go/pkg/crypto"
)

func TestPasswordHashing(t *testing.T) {
	rawPassword := "AdminBPS1200!"

	hash, err := crypto.HashPassword(rawPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !crypto.CheckPasswordHash(rawPassword, hash) {
		t.Errorf("Expected password check to succeed with correct password")
	}

	if crypto.CheckPasswordHash("WrongPassword!", hash) {
		t.Errorf("Expected password check to fail with incorrect password")
	}
}
