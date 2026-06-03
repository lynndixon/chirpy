package auth

import "testing"

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
