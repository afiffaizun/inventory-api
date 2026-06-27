package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/afiffazun/inventory-api/internal/config"
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/handler"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/version", handler.Version)

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetItems(w, r)
		case http.MethodPost:
			handler.CreateItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetItem(w, r)
		case http.MethodPut:
			handler.UpdateItem(w, r)
		case http.MethodDelete:
			handler.DeleteItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Inventory API running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
