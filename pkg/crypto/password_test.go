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

	migrationHash := "$2a$10$CtL0O5Pagh8WXov1KJIqVOmf7w0WxX.EoccJgsFoyJ/JcQNMShZHq"
	if !crypto.CheckPasswordHash("AdminBPS1200!", migrationHash) {
		t.Errorf("Migration hash did not match AdminBPS1200!")
	}
}
