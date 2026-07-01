package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func TestGetPurchaseOrders(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		query      string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "empty list",
			setup:      func(t *testing.T) { setupTest() },
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "returns orders",
			setup: func(t *testing.T) {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				createTestPurchaseOrder(t, supID, wID, item.ID)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "filter by status",
			setup: func(t *testing.T) {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				createTestPurchaseOrder(t, supID, wID, item.ID)
			},
			query:      "?status=DRAFT",
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/purchase-orders"+tt.query, "", GetPurchaseOrders)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.PurchaseOrder]
			decodeJSON(t, w, &resp)
			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d orders, got %d", tt.wantCount, len(resp.Data))
			}
		})
	}
}

func TestCreatePurchaseOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				return fmt.Sprintf(`{"supplier_id":%d,"warehouse_id":%d,"order_date":"2026-07-01","notes":"Test order","created_by":"test","items":[{"item_id":%d,"quantity":50,"unit_cost":15.00}]}`, supID, wID, item.ID)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "validation error - missing supplier_id",
			setup: func(t *testing.T) string {
				setupTest()
				return `{"warehouse_id":1,"order_date":"2026-07-01","items":[]}`
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - invalid date",
			setup: func(t *testing.T) string {
				setupTest()
				return `{"supplier_id":1,"warehouse_id":1,"order_date":"invalid","items":[]}`
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.setup(t)
			w := executeRequest(t, http.MethodPost, "/purchase-orders", body, CreatePurchaseOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestConfirmPurchaseOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				order := createTestPurchaseOrder(t, supID, wID, item.ID)
				return fmt.Sprintf("%d", order.ID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			setup:      func(t *testing.T) string { setupTest(); return "999" },
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			id := tt.setup(t)
			w := executeRequestWithPathValue(t, http.MethodPost, "/purchase-orders/"+id+"/confirm", "id", id, "", ConfirmPurchaseOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestCancelPurchaseOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				order := createTestPurchaseOrder(t, supID, wID, item.ID)
				return fmt.Sprintf("%d", order.ID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			setup:      func(t *testing.T) string { setupTest(); return "999" },
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			id := tt.setup(t)
			w := executeRequestWithPathValue(t, http.MethodPost, "/purchase-orders/"+id+"/cancel", "id", id, "", CancelPurchaseOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestDeletePurchaseOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				supID := createTestSupplierUint(t, "SUP001", "Supplier 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				order := createTestPurchaseOrder(t, supID, wID, item.ID)
				return fmt.Sprintf("%d", order.ID)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "not found",
			setup:      func(t *testing.T) string { setupTest(); return "999" },
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			id := tt.setup(t)
			w := executeRequestWithPathValue(t, http.MethodDelete, "/purchase-orders/"+id, "id", id, "", DeletePurchaseOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func createTestSupplierUint(t *testing.T, code, name string) uint {
	t.Helper()
	supplier := model.Supplier{Code: code, Name: name}
	result := database.DB.Create(&supplier)
	if result.Error != nil {
		t.Fatalf("failed to create test supplier: %v", result.Error)
	}
	return supplier.ID
}

func createTestPurchaseOrder(t *testing.T, supplierID, warehouseID, itemID uint) model.PurchaseOrder {
	t.Helper()
	order := model.PurchaseOrder{
		SupplierID:  supplierID,
		WarehouseID: warehouseID,
		OrderDate:   time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		Status:      model.PurchaseOrderStatusDRAFT,
		CreatedBy:   "test",
		Items: []model.PurchaseOrderItem{
			{ItemID: itemID, Quantity: 50, UnitCost: 15.00, Subtotal: 750.00},
		},
	}
	result := database.DB.Create(&order)
	if result.Error != nil {
		t.Fatalf("failed to create test purchase order: %v", result.Error)
	}
	return order
}
