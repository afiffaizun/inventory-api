package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/afiffazun/inventory-api/internal/config"
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/handler"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	http.HandleFunc("/", corsMiddleware(handler.Home))
	http.HandleFunc("/health", corsMiddleware(handler.Health))
	http.HandleFunc("/version", corsMiddleware(handler.Version))

	http.HandleFunc("/items", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetItems(w, r)
		case http.MethodPost:
			handler.CreateItem(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/items/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Inventory API running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
