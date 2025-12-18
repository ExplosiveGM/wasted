package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ExplosiveGM/wasted/internal/database"
	"github.com/ExplosiveGM/wasted/internal/messaging"
	"github.com/ExplosiveGM/wasted/internal/utils"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog"
)

type Service struct {
	queries *database.Queries
	logger  zerolog.Logger
}

func NewAuthService(queries *database.Queries, logger zerolog.Logger) *Service {
	return &Service{queries: queries, logger: logger}
}

func (s *Service) RequestCode(ctx context.Context, login string) error {
	loginType, err := determineLoginType(login)

	s.logger.Debug().Str("Login", login).Str("LoginType", loginType).Msg("Определение типа логина")

	if err != nil {
		s.logger.Err(err).Str("Login", login).Msg("Ошибка определения типа логина")
		return err
	}

	if loginType == "phone" {
		login, err = normalizePhone(login)

		if err != nil {
			s.logger.Err(err).Str("Login", login).Msg("Ошибка после нормализации телефона")
			return err
		}
	}

	code := utils.GenerateCode(100000, 999999)
	codeStr := strconv.Itoa(code)
	codeExpiresAt := time.Now().Add(10 * time.Minute)

	s.logger.Debug().Str("Code", codeStr).Time("Code expires at", codeExpiresAt).Msg("Сгенерирован код")

	user, err := s.queries.FindUserByLogin(ctx, login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := s.queries.CreateUser(ctx, database.CreateUserParams{
				Code: sql.NullString{String: codeStr}, CodeExpiresAt: sql.NullTime{Time: codeExpiresAt}, Login: login,
			})

			if err != nil {
				s.logger.Err(err).Str("Login", login).Msg("Ошибка при создании пользователя")
			}
		} else {
			return fmt.Errorf("%w: %v", ErrDatabaseUnexpectedError, err)
		}
	}

	err = s.queries.UpdateUser(ctx, database.UpdateUserParams{
		ID: user.ID, Code: sql.NullString{String: codeStr, Valid: true}, CodeExpiresAt: sql.NullTime{Time: codeExpiresAt, Valid: true},
	})

	if err != nil {
		s.logger.Err(err).Str("Login", login).Msg("Ошибка при обновлении пользователя")
	}

	sender := messaging.NewFakeSender()

	switch loginType {
	case "email":
		sender.SendCodeViaEmail(login, codeStr)
	case "phone":
		sender.SendCodeViaSms(login, codeStr)
	}

	return nil
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Verify(ctx context.Context, login string, code int) (TokenPair, error) {
	loginType, err := determineLoginType(login)

	s.logger.Debug().Str("Login", login).Str("LoginType", loginType).Msg("Определение типа логина")

	if err != nil {
		return TokenPair{}, err
	}

	if loginType == "phone" {
		login, err = normalizePhone(login)

		if err != nil {
			s.logger.Err(err).Str("Login", login).Msg("Ошибка после нормализации телефона")
			return TokenPair{}, err
		}
	}

	codeStr := strconv.Itoa(code)
	user, err := s.queries.CheckUserByCode(ctx, database.CheckUserByCodeParams{
		Code: sql.NullString{String: codeStr, Valid: true}, Login: login,
	})

	if err != nil {
		s.logger.Debug().Str("Code", codeStr).Str("Login", login).Msg("User с указанным кодом не найден")
		return TokenPair{}, fmt.Errorf("%w: %v", ErrUserWithCodeNotFound, err)
	}

	accessTokenParams, err := generateAccessToken(user.Login)

	if err != nil {
		s.logger.Err(err).Msg("Ошибка при генерации access-токена")
		return TokenPair{}, fmt.Errorf("%w: %v", ErrGeneratingAccessToken, err)
	}

	refreshTokenParams, err := generateRefreshToken(user.Login)

	if err != nil {
		s.logger.Err(err).Msg("Ошибка при генерации refresh-токена")
		return TokenPair{}, fmt.Errorf("%w: %v", ErrGeneratingRefreshToken, err)
	}

	hashedRefreshToken := HashToken(refreshTokenParams.signedToken)

	err = s.queries.UpdateUserRefreshToken(ctx, database.UpdateUserRefreshTokenParams{
		RefreshToken:          sql.NullString{String: hashedRefreshToken, Valid: true},
		RefreshTokenExpiresAt: sql.NullTime{Time: refreshTokenParams.ttl, Valid: true},
		ID:                    user.ID,
	})

	if err != nil {
		return TokenPair{}, fmt.Errorf("%w: %v", ErrDatabaseUnexpectedError, err)
	}

	s.logger.Debug().Str("Access-Token", accessTokenParams.signedToken).Str("Refresh-Token", refreshTokenParams.signedToken).Msg("Информация по токенам")

	return TokenPair{AccessToken: accessTokenParams.signedToken, RefreshToken: refreshTokenParams.signedToken}, nil
}

type RefreshResult struct {
	AccessToken string `json:"access_token"`
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (RefreshResult, error) {
	hashedRefreshToken := HashToken(refreshToken)

	user, err := s.queries.CheckUserByRefreshCode(ctx, sql.NullString{String: hashedRefreshToken, Valid: true})

	if err != nil {
		return RefreshResult{}, fmt.Errorf("%w: %v", ErrRefreshTokenNotFound, err)
	}

	accessTokenParams, err := generateAccessToken(user.Login)

	if err != nil {
		return RefreshResult{}, fmt.Errorf("%w: %v", ErrGeneratingAccessToken, err)
	}

	return RefreshResult{AccessToken: accessTokenParams.signedToken}, nil
}

func determineLoginType(login string) (string, error) {
	loginType, err := utils.DetermineLoginType(login)

	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrLoginValidation, err)
	}

	return loginType, err
}

func normalizePhone(login string) (string, error) {
	num, err := phonenumbers.Parse(login, "RU")

	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPhoneParsing, err)
	} else {
		login = phonenumbers.Format(num, phonenumbers.E164)

		return login, nil
	}
}
