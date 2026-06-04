package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashAndCheck(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	// correct password should match
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}
	if !match {
		t.Error("expected the correct password to match the hash")
	}

	// wrong password should NOT match
	match, _ = CheckPasswordHash("wrong-password", hash)
	if match {
		t.Error("expected the wrong password to fail")
	}
}

func TestJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	secret := "my-secret-key"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	// Validate the token
	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}
	if returnedUserID.String() != userID {
		t.Errorf("expected user ID %s, got %s", userID, returnedUserID.String())
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	secret := "my-secret-key"
	expiresIn := -1 * time.Hour // already expired

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	// Validate the token, should fail due to expiration
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected ValidateJWT to return an error for expired token")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	secret := "my-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	// Validate the token with the wrong secret, should fail
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("expected ValidateJWT to return an error for wrong secret")
	}
}
