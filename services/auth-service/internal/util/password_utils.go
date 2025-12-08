package internalutil

import "golang.org/x/crypto/bcrypt"

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
		if r <= 32 { // disallow control chars and spaces
			return false
		}
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			// treat other printable runes as special characters
			if r > 32 {
				hasSpecial = true
			}
		}
	}

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

	// For very short passwords (>=4), require at least 2 different categories.
	// For longer passwords (>=8), require at least 3 categories.
	if len(runes) < 8 {
		return categories >= 2
	}
	return categories >= 3
}