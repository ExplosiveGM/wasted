package auth

type RequestCodeInput struct {
	Login string `json:"login" binding:"required"`
}

type VerifyInput struct {
	Login string `json:"login" binding:"required"`
	Code  int    `json:"code" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
