package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewAuthHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RequestCode(c *gin.Context) {
	var jsonSchema struct {
		Login string `json:"login" binding:"required"`
	}

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if jsonSchema.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login required"})
		return
	}

	err := h.service.RequestCode(c.Request.Context(), jsonSchema.Login)

	if err != nil {
		if errors.Is(err, ErrLoginValidation) {
			c.JSON(422, gin.H{"error": "Не удается определить тип логина"})
			return
		} else if errors.Is(err, ErrDatabaseUnexpectedError) {
			c.JSON(500, gin.H{"error": "Ошибка с БД"})
			return
		} else if errors.Is(err, ErrPhoneParsing) {
			c.JSON(422, gin.H{"error": "Невалидный номер телефона"})
			return
		} else {
			c.JSON(500, gin.H{"error": "Неизвестная ошибка"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code sent"})
}

func (h *Handler) Verify(c *gin.Context) {
	var jsonSchema struct {
		Login string `json:"login" binding:"required"`
		Code  int    `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.service.Verify(c.Request.Context(), jsonSchema.Login, jsonSchema.Code)

	if err != nil {
		if errors.Is(err, ErrLoginValidation) {
			c.JSON(422, gin.H{"error": "Не удается определить тип логина"})
			return
		} else if errors.Is(err, ErrDatabaseUnexpectedError) {
			c.JSON(500, gin.H{"error": "Ошибка с БД"})
			return
		} else if errors.Is(err, ErrPhoneParsing) {
			c.JSON(422, gin.H{"error": "Невалидный номер телефона"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenPair})
}

func (h *Handler) Refresh(c *gin.Context) {
	var jsonSchema struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), jsonSchema.RefreshToken)

	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			c.JSON(404, gin.H{"error": "Не удается найти refresh токен или он истек"})
		} else if errors.Is(err, ErrGeneratingAccessToken) {
			c.JSON(500, gin.H{"error": "Ошибка при генерации access_token"})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
