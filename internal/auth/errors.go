package auth

import "errors"

var (
	ErrLoginValidation         = errors.New("login validation problem")
	ErrDatabaseUnexpectedError = errors.New("database unexpected error")
	ErrPhoneParsing            = errors.New("phone parsing problem")
	ErrUserWithCodeNotFound    = errors.New("user with code not found or code expired")
	ErrGeneratingAccessToken   = errors.New("error on generating access token")
	ErrGeneratingRefreshToken  = errors.New("error on generating refresh token")
	ErrRefreshTokenNotFound    = errors.New("refresh token not found or expired")
)
