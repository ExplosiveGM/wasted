package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type TokenParams struct {
	signedToken string
	ttl         time.Time
}

func generateAccessToken(login string) (*TokenParams, error) {
	ttl := time.Now().Add(15 * time.Minute)
	secretKey := viper.GetString("JWT_ACCESS_SECRET")

	return generateToken(login, secretKey, ttl)
}

func generateRefreshToken(login string) (*TokenParams, error) {
	ttl := time.Now().Add(7 * 24 * time.Hour)
	secretKey := viper.GetString("JWT_REFRESH_SECRET")

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
