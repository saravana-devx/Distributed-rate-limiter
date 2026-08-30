package client

type CreateClientRequest struct {
	ClientID      string        `json:"clientId" binding:"required"`
	Algorithm     algorithmType `json:"algorithm" binding:"required,oneof=fixed_window sliding_window token_bucket"`
	Limit         int           `json:"limit" binding:"required,min=1"`
	WindowSeconds int           `json:"window_seconds" binding:"required,min=1"`
}

type UpdateClientRequest struct {
	Algorithm     algorithmType `json:"algorithm" binding:"required,oneof=fixed_window sliding_window token_bucket"`
	Limit         int           `json:"limit" binding:"required,min=1"`
	WindowSeconds int           `json:"window_seconds" binding:"required,min=1"`
}
