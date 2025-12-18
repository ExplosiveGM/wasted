package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewAuthHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RequestCode godoc
//
//	@Description	Send code to user email or phone
//	@Tags			auth
//	@Param			input	body	RequestCodeInput	true	"Login info"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CodeSentResponse
//	@Summary		Sent code user
//	@Router			/api/v1/auth/request-code [post]
func (h *Handler) RequestCode(c *gin.Context) {
	var jsonSchema RequestCodeInput

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if jsonSchema.Login == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "login required"})
		return
	}

	err := h.service.RequestCode(c.Request.Context(), jsonSchema.Login)

	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, CodeSentResponse{Data: CodeSentData{Message: "Code sent"}})
}

// Verify godoc
//
//	@Description	Verify code for user
//	@Tags			auth
//	@Param			input	body	VerifyInput	true	"Login and code"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	TokenPairResponse
//	@Summary		Verify code from user
//	@Router			/api/v1/auth/verify [post]
func (h *Handler) Verify(c *gin.Context) {
	var jsonSchema VerifyInput

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokenPair, err := h.service.Verify(c.Request.Context(), jsonSchema.Login, jsonSchema.Code)

	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, TokenPairResponse{Data: tokenPair})
}

// Refresh godoc
//
//	@Description	Refresh access token
//	@Tags			auth
//	@Param			input	body	RefreshInput	true	"Refresh code"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RefreshResponse
//	@Summary		Refresh token
//	@Router			/api/v1/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var jsonSchema RefreshInput

	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), jsonSchema.RefreshToken)

	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{Data: result})
}
