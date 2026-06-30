package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afiffazun/inventory-api/internal/model"
)

func TestGetWarehouses(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		query      string
		wantStatus int
		wantCount  int
		wantTotal  int64
	}{
		{
			name:       "empty list",
			setup:      func(t *testing.T) { setupTest() },
			wantStatus: http.StatusOK,
			wantCount:  0,
			wantTotal:  0,
		},
		{
			name: "returns warehouses",
			setup: func(t *testing.T) {
				setupTest()
				createTestWarehouse(t, "WH001", "Warehouse A")
				createTestWarehouse(t, "WH002", "Warehouse B")
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
			wantTotal:  2,
		},
		{
			name: "search by name",
			setup: func(t *testing.T) {
				setupTest()
				createTestWarehouse(t, "WH001", "Warehouse A")
				createTestWarehouse(t, "WH002", "Warehouse B")
			},
			query:      "?search=Warehouse+A",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			path := "/warehouses" + tt.query
			w := executeRequestWithQuery(t, http.MethodGet, path, "", GetWarehouses)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.Warehouse]
			decodeJSON(t, w, &resp)

			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d warehouses, got %d", tt.wantCount, len(resp.Data))
			}

			if resp.Total != tt.wantTotal {
				t.Errorf("expected total %d, got %d", tt.wantTotal, resp.Total)
			}
		})
	}
}

func TestCreateWarehouse(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			body:       `{"code":"WH001","name":"Warehouse A","city":"Jakarta","country":"Indonesia"}`,
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Warehouse
				decodeJSON(t, w, &resp)
				if resp.ID == 0 {
					t.Error("expected non-zero ID")
				}
				if resp.Code != "WH001" {
					t.Errorf("expected code 'WH001', got '%s'", resp.Code)
				}
			},
		},
		{
			name:       "validation error - empty code",
			body:       `{"name":"Warehouse A"}`,
			wantStatus: http.StatusBadRequest,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				if len(resp.Errors) != 1 || resp.Errors[0].Field != "code" {
					t.Errorf("expected code error, got %v", resp.Errors)
				}
			},
		},
		{
			name:       "validation error - empty name",
			body:       `{"code":"WH001"}`,
			wantStatus: http.StatusBadRequest,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				if len(resp.Errors) != 1 || resp.Errors[0].Field != "name" {
					t.Errorf("expected name error, got %v", resp.Errors)
				}
			},
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

			w := executeRequest(t, http.MethodPost, "/warehouses", tt.body, CreateWarehouse)
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

func TestGetWarehouse(t *testing.T) {
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
				return createTestWarehouse(t, "WH001", "Warehouse A")
			},
			pathID:     "1",
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Warehouse
				decodeJSON(t, w, &resp)
				if resp.Code != "WH001" {
					t.Errorf("expected code 'WH001', got '%s'", resp.Code)
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

			w := executeRequestWithPathValue(t, http.MethodGet, "/warehouses/"+pathValue, "id", pathValue, "", GetWarehouse)
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

func TestUpdateWarehouse(t *testing.T) {
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
				return createTestWarehouse(t, "WH001", "Warehouse A")
			},
			body:       `{"code":"WH001","name":"Warehouse Updated","city":"Surabaya"}`,
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Warehouse
				decodeJSON(t, w, &resp)
				if resp.Name != "Warehouse Updated" {
					t.Errorf("expected name 'Warehouse Updated', got '%s'", resp.Name)
				}
			},
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			body:       `{"code":"WH001","name":"Warehouse Updated"}`,
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
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

			w := executeRequestWithPathValue(t, http.MethodPut, "/warehouses/"+pathValue, "id", pathValue, tt.body, UpdateWarehouse)
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

func TestDeleteWarehouse(t *testing.T) {
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
				return createTestWarehouse(t, "WH001", "Warehouse A")
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

			w := executeRequestWithPathValue(t, http.MethodDelete, "/warehouses/"+pathValue, "id", pathValue, "", DeleteWarehouse)
			assertStatus(t, w, tt.wantStatus)

			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
		})
	}
}
