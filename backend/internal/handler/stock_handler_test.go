package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"github.com/afiffazun/inventory-api/internal/repository"
)

func TestGetStockMovements(t *testing.T) {
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
			name: "returns movements",
			setup: func(t *testing.T) {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeIN, 100)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeOUT, 20)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "filter by type",
			setup: func(t *testing.T) {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeIN, 100)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeOUT, 20)
			},
			query:      "?type=IN",
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/stock-movements"+tt.query, "", GetStockMovements)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.StockMovement]
			decodeJSON(t, w, &resp)
			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d movements, got %d", tt.wantCount, len(resp.Data))
			}
		})
	}
}

func TestGetStockMovement(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		setup      func(t *testing.T) string
		wantStatus int
		wantCode   string
	}{
		{
			name:   "success",
			pathID: "1",
			setup: func(t *testing.T) string {
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				m := createTestStockMovement(t, item.ID, wID, model.MovementTypeIN, 100)
				return fmt.Sprintf("%d", m.ID)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathID:     "999",
			setup:      func(t *testing.T) string { return "" },
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "invalid ID",
			pathID:     "abc",
			setup:      func(t *testing.T) string { return "" },
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			pathValue := tt.pathID
			if tt.setup != nil {
				if id := tt.setup(t); id != "" {
					pathValue = id
				}
			}
			w := executeRequestWithPathValue(t, http.MethodGet, "/stock-movements/"+pathValue, "id", pathValue, "", GetStockMovement)
			assertStatus(t, w, tt.wantStatus)
			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
		})
	}
}

func TestCreateStockMovement(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success IN",
			setup: func(t *testing.T) string {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				return fmt.Sprintf(`{"item_id":%d,"warehouse_id":%d,"type":"IN","quantity":100,"created_by":"test"}`, item.ID, wID)
			},
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.StockMovement
				decodeJSON(t, w, &resp)
				if resp.ID == 0 {
					t.Error("expected non-zero ID")
				}
				if resp.Type != model.MovementTypeIN {
					t.Errorf("expected type 'IN', got '%s'", resp.Type)
				}
			},
		},
		{
			name: "success OUT",
			setup: func(t *testing.T) string {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 50, wID)
				return fmt.Sprintf(`{"item_id":%d,"warehouse_id":%d,"type":"OUT","quantity":20,"created_by":"test"}`, item.ID, wID)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "validation error - missing item_id",
			setup:      func(t *testing.T) string { setupTest(); return `{"warehouse_id":1,"type":"IN","quantity":100}` },
			wantStatus: http.StatusBadRequest,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				found := false
				for _, e := range resp.Errors {
					if e.Field == "item_id" {
						found = true
					}
				}
				if !found {
					t.Error("expected item_id validation error")
				}
			},
		},
		{
			name: "validation error - invalid type",
			setup: func(t *testing.T) string {
				setupTest()
				return `{"item_id":1,"warehouse_id":1,"type":"INVALID","quantity":100}`
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - zero quantity",
			setup: func(t *testing.T) string {
				setupTest()
				return `{"item_id":1,"warehouse_id":1,"type":"IN","quantity":0}`
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body
			if tt.setup != nil {
				body = tt.setup(t)
			}
			w := executeRequest(t, http.MethodPost, "/stock-movements", body, CreateStockMovement)
			assertStatus(t, w, tt.wantStatus)
			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
			if tt.checkResp != nil {
				tt.checkResp(t, w)
			}
		})
	}
}

func TestGetStockHistory(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		setup      func(t *testing.T) string
		wantStatus int
		wantCount  int
	}{
		{
			name:   "success",
			pathID: "1",
			setup: func(t *testing.T) string {
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 0, wID)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeIN, 100)
				createTestStockMovement(t, item.ID, wID, model.MovementTypeOUT, 20)
				return fmt.Sprintf("%d", item.ID)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "invalid ID",
			pathID:     "abc",
			setup:      func(t *testing.T) string { return "" },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			pathValue := tt.pathID
			if tt.setup != nil {
				if id := tt.setup(t); id != "" {
					pathValue = id
				}
			}
			w := executeRequestWithPathValue(t, http.MethodGet, "/items/"+pathValue+"/stock-history", "id", pathValue, "", GetStockHistory)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestGetStockSummary(t *testing.T) {
	setupTest()

	w := executeRequest(t, http.MethodGet, "/stock-summary", "", GetStockSummary)
	assertStatus(t, w, http.StatusOK)

	var resp []repository.StockSummary
	decodeJSON(t, w, &resp)
	if len(resp) != 0 {
		t.Errorf("expected 0 summary entries, got %d", len(resp))
	}
}

func TestGetStockOpnames(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T)
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
			name: "returns opnames",
			setup: func(t *testing.T) {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 50, wID)
				opname := model.StockOpname{
					WarehouseID: wID,
					Status:      model.OpnameStatusDRAFT,
					Items: []model.StockOpnameItem{
						{ItemID: item.ID, SystemQuantity: 50, ActualQuantity: 48, Difference: -2},
					},
				}
				database.DB.Create(&opname)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/stock-opnames", "", GetStockOpnames)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.StockOpname]
			decodeJSON(t, w, &resp)
			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d opnames, got %d", tt.wantCount, len(resp.Data))
			}
		})
	}
}

func TestCreateStockOpname(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		wantStatus int
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				setupTest()
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 50, wID)
				return fmt.Sprintf(`{"warehouse_id":%d,"notes":"Test opname","created_by":"test","items":[{"item_id":%d,"actual_quantity":48}]}`, wID, item.ID)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "validation error - missing warehouse_id",
			setup:      func(t *testing.T) string { setupTest(); return `{"notes":"Test","items":[]}` },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := ""
			if tt.setup != nil {
				body = tt.setup(t)
			}
			w := executeRequest(t, http.MethodPost, "/stock-opnames", body, CreateStockOpname)
			assertStatus(t, w, tt.wantStatus)
		})
	}
}

func TestCompleteStockOpname(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string
		pathID     string
		wantStatus int
		wantCode   string
	}{
		{
			name: "success",
			setup: func(t *testing.T) string {
				wID := createTestWarehouseUint(t, "WH001", "Warehouse 1")
				item := createTestItemWithWarehouse(t, "ITM001", "Item 1", 50, wID)
				opname := model.StockOpname{
					WarehouseID: wID,
					Status:      model.OpnameStatusDRAFT,
					Items: []model.StockOpnameItem{
						{ItemID: item.ID, SystemQuantity: 50, ActualQuantity: 48, Difference: -2},
					},
				}
				database.DB.Create(&opname)
				return strconv.FormatUint(uint64(opname.ID), 10)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			pathID:     "999",
			setup:      func(t *testing.T) string { return "" },
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			pathValue := tt.pathID
			if tt.setup != nil {
				if id := tt.setup(t); id != "" {
					pathValue = id
				}
			}
			w := executeRequestWithPathValue(t, http.MethodPost, "/stock-opnames/"+pathValue+"/complete", "id", pathValue, "", CompleteStockOpname)
			assertStatus(t, w, tt.wantStatus)
			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
		})
	}
}
