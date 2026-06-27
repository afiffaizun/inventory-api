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