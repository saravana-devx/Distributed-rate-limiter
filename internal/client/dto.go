package client

type CreateClientRequest struct {
	ClientID      string        `json:"clientId"`
	Algorithm     algorithmType `json:"algorithm"`
	Limit         int           `json:"limit"`
	WindowSeconds int           `json:"window_seconds"`
}

type UpdateClientRequest struct {
	Algorithm     algorithmType `json:"algorithm"`
	Limit         int           `json:"limit"`
	WindowSeconds int           `json:"window_seconds"`
}
