package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(serviceError error, c *gin.Context) {
	switch {
	case errors.Is(serviceError, ErrUserWithCodeNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Неизвестная ошибка"})
	case errors.Is(serviceError, ErrLoginValidation):
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: "Не удается определить тип логина"})
	case errors.Is(serviceError, ErrPhoneParsing):
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: "Невалидный номер телефона"})
	case errors.Is(serviceError, ErrRefreshTokenNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Не удается найти refresh токен или он истек"})
	case errors.Is(serviceError, ErrGeneratingAccessToken) || errors.Is(serviceError, ErrGeneratingRefreshToken):
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Error: "Неизвестная ошибка"})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Неизвестная ошибка"})
	}
}
