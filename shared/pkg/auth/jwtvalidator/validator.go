package jwtvalidator

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserEmail string
	Extra     map[string]any
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenID   string
	UserEmail string
	Extra     map[string]any
	jwt.RegisteredClaims
}

// decodePEM cleans and decodes a PEM formatted string.
func decodePEM(value string) ([]byte, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return nil, errors.New("key is empty")
	}

	if !strings.Contains(cleaned, "-----BEGIN") {
		return nil, errors.New("key is not a valid PEM block")
	}

	return []byte(cleaned), nil
}

// ParseRSAPrivateKeyFromString parses an RSA private key from a PEM formatted string.
func ParseRSAPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	pemBytes, err := decodePEM(key)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA private key: %w", err)
	}

	return privateKey, nil
}

// ParseRSAPublicKeyFromString parses an RSA public key from a PEM formatted string.
func ParseRSAPublicKeyFromString(key string) (*rsa.PublicKey, error) {
	pemBytes, err := decodePEM(key)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA public key: %w", err)
	}

	return publicKey, nil
}

// validateAccessTokenWithKey validates an access token using the provided RSA public key.
func validateAccessTokenWithKey(publicKey *rsa.PublicKey, tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// ValidateAccessToken validates an access token using the provided PEM formatted RSA public key string.
func ValidateAccessToken(publicKeyPEM string, tokenStr string) (*AccessClaims, error) {
	publicKey, err := ParseRSAPublicKeyFromString(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid access token public key: %w", err)
	}

	return validateAccessTokenWithKey(publicKey, tokenStr)
}

// validateRefreshTokenWithKey validates a refresh token using the provided RSA public key.
func validateRefreshTokenWithKey(publicKey *rsa.PublicKey, tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token using the provided PEM formatted RSA public key string.
func ValidateRefreshToken(publicKeyPEM string, tokenStr string) (*RefreshClaims, error) {
	publicKey, err := ParseRSAPublicKeyFromString(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token public key: %w", err)
	}

	return validateRefreshTokenWithKey(publicKey, tokenStr)
}
