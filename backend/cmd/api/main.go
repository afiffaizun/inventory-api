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

	// Stock Movement routes
	http.HandleFunc("/stock-movements", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetStockMovements(w, r)
		case http.MethodPost:
			handler.CreateStockMovement(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/stock-movements/transfer", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.TransferStock(w, r)
	}))

	http.HandleFunc("/stock-movements/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStockMovement(w, r)
	}))

	// Stock History route
	http.HandleFunc("/items/{id}/stock-history", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStockHistory(w, r)
	}))

	// Stock by Warehouse route
	http.HandleFunc("/warehouses/{id}/stock", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStockByWarehouse(w, r)
	}))

	// Stock Summary route
	http.HandleFunc("/stock-summary", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStockSummary(w, r)
	}))

	// Stock Opname routes
	http.HandleFunc("/stock-opnames", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetStockOpnames(w, r)
		case http.MethodPost:
			handler.CreateStockOpname(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/stock-opnames/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetStockOpname(w, r)
	}))

	http.HandleFunc("/stock-opnames/{id}/complete", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.CompleteStockOpname(w, r)
	}))

	// Customer routes
	http.HandleFunc("/customers", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCustomers(w, r)
		case http.MethodPost:
			handler.CreateCustomer(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/customers/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCustomer(w, r)
		case http.MethodPut:
			handler.UpdateCustomer(w, r)
		case http.MethodDelete:
			handler.DeleteCustomer(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Supplier routes
	http.HandleFunc("/suppliers", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSuppliers(w, r)
		case http.MethodPost:
			handler.CreateSupplier(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/suppliers/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSupplier(w, r)
		case http.MethodPut:
			handler.UpdateSupplier(w, r)
		case http.MethodDelete:
			handler.DeleteSupplier(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Sales Order routes
	http.HandleFunc("/sales-orders", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSalesOrders(w, r)
		case http.MethodPost:
			handler.CreateSalesOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/sales-orders/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSalesOrder(w, r)
		case http.MethodPut:
			handler.UpdateSalesOrder(w, r)
		case http.MethodDelete:
			handler.DeleteSalesOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/sales-orders/{id}/confirm", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.ConfirmSalesOrder(w, r)
	}))

	http.HandleFunc("/sales-orders/{id}/cancel", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.CancelSalesOrder(w, r)
	}))

	// Purchase Order routes
	http.HandleFunc("/purchase-orders", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetPurchaseOrders(w, r)
		case http.MethodPost:
			handler.CreatePurchaseOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/purchase-orders/{id}", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetPurchaseOrder(w, r)
		case http.MethodPut:
			handler.UpdatePurchaseOrder(w, r)
		case http.MethodDelete:
			handler.DeletePurchaseOrder(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/purchase-orders/{id}/confirm", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.ConfirmPurchaseOrder(w, r)
	}))

	http.HandleFunc("/purchase-orders/{id}/cancel", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.CancelPurchaseOrder(w, r)
	}))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Inventory API running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
