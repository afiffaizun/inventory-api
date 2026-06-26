package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Application string `json:"application,omitempty"`
	Author      string `json:"author,omitempty"`
	Status      string `json:"status,omitempty"`
	Version     string `json:"version,omitempty"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{
		Application: "Inventory API",
		Author:      "Afif",
		Status:      "Running",
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	})
}

func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{
		Version: "v1.0.0",
	})
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/health", health)
	http.HandleFunc("/version", version)

	log.Println("Inventory API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}