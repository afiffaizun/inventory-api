package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/afiffazun/inventory-api/internal/config"
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func TestMain(m *testing.M) {
	cfg := config.Load()
	database.Connect(cfg)
	database.Migrate()

	os.Exit(m.Run())
}

func setupTest() {
	database.DB.Exec("DELETE FROM items")
}

func TestHome(t *testing.T) {
	setupTest()

	w := executeRequest(t, http.MethodGet, "/", "", Home)
	assertStatus(t, w, http.StatusOK)

	var resp model.Response
	decodeJSON(t, w, &resp)

	if resp.Application != "Inventory API" {
		t.Errorf("expected application 'Inventory API', got '%s'", resp.Application)
	}
}

func TestHealth(t *testing.T) {
	setupTest()

	w := executeRequest(t, http.MethodGet, "/health", "", Health)
	assertStatus(t, w, http.StatusOK)

	var resp map[string]string
	decodeJSON(t, w, &resp)

	if resp["status"] != "UP" {
		t.Errorf("expected status 'UP', got '%s'", resp["status"])
	}
}

func TestVersion(t *testing.T) {
	setupTest()

	w := executeRequest(t, http.MethodGet, "/version", "", Version)
	assertStatus(t, w, http.StatusOK)

	var resp model.Response
	decodeJSON(t, w, &resp)

	if resp.Version != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got '%s'", resp.Version)
	}
}

func TestGetItems(t *testing.T) {
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
			name: "returns items",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			w := executeRequest(t, http.MethodGet, "/items", "", GetItems)
			assertStatus(t, w, tt.wantStatus)

			var resp []model.Item
			decodeJSON(t, w, &resp)

			if len(resp) != tt.wantCount {
				t.Errorf("expected %d items, got %d", tt.wantCount, len(resp))
			}
		})
	}
}

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			body:       `{"code":"ITEM001","name":"Laptop ACER","stock":10,"location":"Warehouse A"}`,
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Item
				decodeJSON(t, w, &resp)
				if resp.ID == 0 {
					t.Error("expected non-zero ID")
				}
				if resp.Code != "ITEM001" {
					t.Errorf("expected code 'ITEM001', got '%s'", resp.Code)
				}
			},
		},
		{
			name:       "validation error - empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
		{
			name:       "validation error - invalid JSON",
			body:       `invalid json`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()

			w := executeRequest(t, http.MethodPost, "/items", tt.body, CreateItem)
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

func TestGetItem(t *testing.T) {
	tests := []struct {
		name       string
		setupID    func(t *testing.T) string
		pathID     string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupID: func(t *testing.T) string {
				return createTestItem(t, "ITEM001", "Laptop", 5, "Gudang A")
			},
			pathID:     "1",
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Item
				decodeJSON(t, w, &resp)
				if resp.Code != "ITEM001" {
					t.Errorf("expected code 'ITEM001', got '%s'", resp.Code)
				}
			},
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "invalid ID",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()

			pathValue := tt.pathID
			if tt.setupID != nil {
				if id := tt.setupID(t); id != "" {
					pathValue = id
				}
			}

			w := executeRequestWithPathValue(t, http.MethodGet, "/items/"+pathValue, "id", pathValue, "", GetItem)
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

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name       string
		setupID    func(t *testing.T) string
		pathID     string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupID: func(t *testing.T) string {
				return createTestItem(t, "ITEM001", "Laptop", 5, "Gudang A")
			},
			body:       `{"code":"ITEM001","name":"Laptop Updated","stock":10,"location":"Gudang B"}`,
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Item
				decodeJSON(t, w, &resp)
				if resp.Name != "Laptop Updated" {
					t.Errorf("expected name 'Laptop Updated', got '%s'", resp.Name)
				}
			},
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			body:       `{"code":"ITEM001","name":"Laptop Updated","stock":10,"location":"Gudang B"}`,
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "invalid ID",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "abc",
			body:       `{"code":"ITEM001","name":"Laptop Updated","stock":10,"location":"Gudang B"}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()

			pathValue := tt.pathID
			if tt.setupID != nil {
				if id := tt.setupID(t); id != "" {
					pathValue = id
				}
			}

			w := executeRequestWithPathValue(t, http.MethodPut, "/items/"+pathValue, "id", pathValue, tt.body, UpdateItem)
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

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name       string
		setupID    func(t *testing.T) string
		pathID     string
		wantStatus int
		wantCode   string
	}{
		{
			name: "success",
			setupID: func(t *testing.T) string {
				return createTestItem(t, "ITEM001", "Laptop", 5, "Gudang A")
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "invalid ID",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()

			pathValue := tt.pathID
			if tt.setupID != nil {
				if id := tt.setupID(t); id != "" {
					pathValue = id
				}
			}

			w := executeRequestWithPathValue(t, http.MethodDelete, "/items/"+pathValue, "id", pathValue, "", DeleteItem)
			assertStatus(t, w, tt.wantStatus)

			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}

			if tt.name == "success" {
				var count int64
				database.DB.Model(&model.Item{}).Count(&count)
				if count != 0 {
					t.Errorf("expected 0 items after delete, got %d", count)
				}
			}
		})
	}
}
