package main

import (
	"log"
	"net/http"

	"github.com/afiffazun/inventory-api/internal/handler"
)

func main() {

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/version", handler.Version)

	log.Println("Inventory API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}