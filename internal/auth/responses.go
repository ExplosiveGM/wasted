package auth

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse[T any] struct {
	Data T `json:"data,omitempty"`
}

type CodeSentData struct {
	Message string `json:"message" example:"Code sent"`
}

type CodeSentResponse = DataResponse[CodeSentData]

type TokenPairResponse struct {
	Data TokenPair `json:"data" example: "{access_token:111, refresh_token: 222}"`
}

type RefreshResponse struct {
	Data RefreshResult `json:"data" example: "{refresh_token: 222}"`
}
