package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func TestGetSalesOrders(t *testing.T) {
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
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				createTestSalesOrder(t, custID, wID, item.ID)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name: "filter by status",
			setup: func(t *testing.T) {
				setupTest()
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				createTestSalesOrder(t, custID, wID, item.ID)
			},
			query:      "?status=DRAFT",
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/sales-orders"+tt.query, "", GetSalesOrders)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.SalesOrder]
			decodeJSON(t, w, &resp)
			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d orders, got %d", tt.wantCount, len(resp.Data))
			}
		})
	}
}

func TestCreateSalesOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				return fmt.Sprintf(`{"customer_id":%d,"warehouse_id":%d,"order_date":"2026-07-01","notes":"Test order","created_by":"test","items":[{"item_id":%d,"quantity":10,"unit_price":25.50}]}`, custID, wID, item.ID)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "validation error - missing customer_id",
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
				return `{"customer_id":1,"warehouse_id":1,"order_date":"invalid","items":[]}`
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.setup(t)
			w := executeRequest(t, http.MethodPost, "/sales-orders", body, CreateSalesOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestConfirmSalesOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				order := createTestSalesOrder(t, custID, wID, item.ID)
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
			w := executeRequestWithPathValue(t, http.MethodPost, "/sales-orders/"+id+"/confirm", "id", id, "", ConfirmSalesOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestCancelSalesOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				order := createTestSalesOrder(t, custID, wID, item.ID)
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
			w := executeRequestWithPathValue(t, http.MethodPost, "/sales-orders/"+id+"/cancel", "id", id, "", CancelSalesOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestDeleteSalesOrder(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				custID := createTestCustomerUint(t, "CUST001", "Customer 1")
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 100, wID)
				order := createTestSalesOrder(t, custID, wID, item.ID)
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
			w := executeRequestWithPathValue(t, http.MethodDelete, "/sales-orders/"+id, "id", id, "", DeleteSalesOrder)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func createTestCustomerUint(t *testing.T, code, name string) uint {
	t.Helper()
	customer := model.Customer{Code: code, Name: name}
	result := database.DB.Create(&customer)
	if result.Error != nil {
		t.Fatalf("failed to create test customer: %v", result.Error)
	}
	return customer.ID
}

func createTestSalesOrder(t *testing.T, customerID, warehouseID, itemID uint) model.SalesOrder {
	t.Helper()
	order := model.SalesOrder{
		CustomerID:  customerID,
		WarehouseID: warehouseID,
		OrderDate:   time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		Status:      model.SalesOrderStatusDRAFT,
		CreatedBy:   "test",
		Items: []model.SalesOrderItem{
			{ItemID: itemID, Quantity: 10, UnitPrice: 25.50, Subtotal: 255.00},
		},
	}
	result := database.DB.Create(&order)
	if result.Error != nil {
		t.Fatalf("failed to create test sales order: %v", result.Error)
	}
	return order
}
