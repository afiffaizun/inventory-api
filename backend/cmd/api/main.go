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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
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

	// Items routes
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

	http.HandleFunc("/items/export", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.ExportItems(w, r)
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

	// Warehouse routes
	http.HandleFunc("/warehouses", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetWarehouses(w, r)
		case http.MethodPost:
			handler.CreateWarehouse(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/warehouses/all", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetAllWarehouses(w, r)
	}))

	http.HandleFunc("/warehouses/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetWarehouse(w, r)
		case http.MethodPut:
			handler.UpdateWarehouse(w, r)
		case http.MethodDelete:
			handler.DeleteWarehouse(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/warehouses/{id}/set-default", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.SetDefaultWarehouse(w, r)
	}))

	// Category routes
	http.HandleFunc("/categories", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCategories(w, r)
		case http.MethodPost:
			handler.CreateCategory(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/categories/all", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetAllCategories(w, r)
	}))

	http.HandleFunc("/categories/tree", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetCategoryTree(w, r)
	}))

	http.HandleFunc("/categories/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCategory(w, r)
		case http.MethodPut:
			handler.UpdateCategory(w, r)
		case http.MethodDelete:
			handler.DeleteCategory(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Inventory API running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
