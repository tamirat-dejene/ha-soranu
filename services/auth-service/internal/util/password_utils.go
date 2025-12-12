package internalutil

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plain text password.
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsPasswordStrong checks if the given password meets strength requirements.
func IsPasswordStrong(password string) bool {
	const minPasswordLength = 4
	runes := []rune(password)

	if len(runes) < minPasswordLength {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range runes {
		// Reject whitespace or control characters
		if unicode.IsSpace(r) || r < 33 {
			return false
		}

		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			// Any other printable character counts as special
			hasSpecial = true
		}
	}

	// Count categories
	categories := 0
	if hasUpper {
		categories++
	}
	if hasLower {
		categories++
	}
	if hasDigit {
		categories++
	}
	if hasSpecial {
		categories++
	}

	// Rules:
	// - Short passwords (4-7 chars): at least 2 categories
	// - Long passwords (8+ chars): at least 3 categories
	if len(runes) < 8 {
		return categories >= 2
	}
	return categories >= 3
}
