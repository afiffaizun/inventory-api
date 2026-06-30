package model

type Response struct {
	Application string `json:"application,omitempty"`
	Author      string `json:"author,omitempty"`
	Status      string `json:"status,omitempty"`
	Version     string `json:"version,omitempty"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
