package jwtvalidator

import (
	"crypto/rsa"
	"encoding/base64"
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

func decodePEMOrBase64(value string) ([]byte, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return nil, errors.New("key is empty")
	}

	if strings.Contains(cleaned, "-----BEGIN") {
		return []byte(cleaned), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err == nil && len(decoded) > 0 {
		return decoded, nil
	}

	return nil, errors.New("key must be a PEM block or base64-encoded PEM")
}

func ParseRSAPrivateKeyFromString(key string) (*rsa.PrivateKey, error) {
	pemBytes, err := decodePEMOrBase64(key)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA private key: %w", err)
	}

	return privateKey, nil
}

func ParseRSAPublicKeyFromString(key string) (*rsa.PublicKey, error) {
	pemBytes, err := decodePEMOrBase64(key)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid RSA public key: %w", err)
	}

	return publicKey, nil
}

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

func ValidateAccessToken(publicKeyPEM string, tokenStr string) (*AccessClaims, error) {
	publicKey, err := ParseRSAPublicKeyFromString(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid access token public key: %w", err)
	}

	return validateAccessTokenWithKey(publicKey, tokenStr)
}

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

func ValidateRefreshToken(publicKeyPEM string, tokenStr string) (*RefreshClaims, error) {
	publicKey, err := ParseRSAPublicKeyFromString(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token public key: %w", err)
	}

	return validateRefreshTokenWithKey(publicKey, tokenStr)
}
