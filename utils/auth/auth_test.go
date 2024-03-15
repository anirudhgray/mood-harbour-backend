package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		password   string
		shouldPass bool
	}{
		{"Abcdef1!", true},
		{"short1!", false},
		{"NoNumber@", false},
		{"NoSymbol123", false},
	}

	for _, test := range tests {
		result := CheckPasswordStrength(test.password)
		if result != test.shouldPass {
			t.Errorf("CheckPasswordStrength(%s) returned %v; expected %v", test.password, result, test.shouldPass)
		}
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "mySecretPassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := VerifyPassword(password, string(hashedPassword))
	if err != nil {
		t.Errorf("VerifyPassword failed: %v", err)
	}

	wrongPassword := "wrongPassword"
	err = VerifyPassword(wrongPassword, string(hashedPassword))
	if err == nil {
		t.Errorf("VerifyPassword should have failed for the wrong password")
	}
}

// Note: The `LoginCheck` function involves database interactions and token generation.
// It's recommended to use an integration test framework or mock the database interactions for unit tests.
