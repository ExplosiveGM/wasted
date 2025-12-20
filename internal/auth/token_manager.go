package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ExplosiveGM/wasted/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenParams struct {
	signedToken string
	ttl         time.Time
}

func generateAccessToken(login string, jwtConfig *config.JWTConfig) (*TokenParams, error) {
	ttl := time.Now().Add(15 * time.Minute)
	secretKey := jwtConfig.AccessSecret

	return generateToken(login, secretKey, ttl)
}

func generateRefreshToken(login string, jwtConfig *config.JWTConfig) (*TokenParams, error) {
	ttl := time.Now().Add(7 * 24 * time.Hour)
	secretKey := jwtConfig.RefreshSecret

	return generateToken(login, secretKey, ttl)
}

func generateToken(login string, secretKey string, ttl time.Time) (*TokenParams, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nickname": login,
		"exp":      ttl.Unix(),
	})

	bytedSecretKey := []byte(secretKey)
	signedToken, err := token.SignedString(bytedSecretKey)

	if err != nil {
		return nil, fmt.Errorf("some error on generating token: %w", err)
	}

	return &TokenParams{signedToken: signedToken, ttl: ttl}, nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
