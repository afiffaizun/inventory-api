package model

type Response struct {
	Application string `json:"application,omitempty"`
	Author      string `json:"author,omitempty"`
	Status      string `json:"status,omitempty"`
	Version     string `json:"version,omitempty"`
}