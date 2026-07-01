package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func TestGetSuppliers(t *testing.T) {
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
			name: "returns suppliers",
			setup: func(t *testing.T) {
				setupTest()
				createTestSupplier(t, "SUP001", "Supplier One")
				createTestSupplier(t, "SUP002", "Supplier Two")
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
			wantTotal:  2,
		},
		{
			name: "search by name",
			setup: func(t *testing.T) {
				setupTest()
				createTestSupplier(t, "SUP001", "Acme Corp")
				createTestSupplier(t, "SUP002", "Beta Inc")
			},
			query:      "?search=Acme",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/suppliers"+tt.query, "", GetSuppliers)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.Supplier]
			decodeJSON(t, w, &resp)

			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d suppliers, got %d", tt.wantCount, len(resp.Data))
			}
			if resp.Total != tt.wantTotal {
				t.Errorf("expected total %d, got %d", tt.wantTotal, resp.Total)
			}
		})
	}
}

func TestCreateSupplier(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			body:       `{"code":"SUP001","name":"Supplier One","email":"sup@test.com","phone":"12345"}`,
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Supplier
				decodeJSON(t, w, &resp)
				if resp.ID == 0 {
					t.Error("expected non-zero ID")
				}
				if resp.Code != "SUP001" {
					t.Errorf("expected code 'SUP001', got '%s'", resp.Code)
				}
			},
		},
		{
			name:       "validation error - empty code",
			body:       `{"name":"Supplier One"}`,
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
			body:       `{"code":"SUP001"}`,
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
			w := executeRequest(t, http.MethodPost, "/suppliers", tt.body, CreateSupplier)
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

func TestGetSupplier(t *testing.T) {
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
				return createTestSupplier(t, "SUP001", "Supplier One")
			},
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Supplier
				decodeJSON(t, w, &resp)
				if resp.Code != "SUP001" {
					t.Errorf("expected code 'SUP001', got '%s'", resp.Code)
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
			w := executeRequestWithPathValue(t, http.MethodGet, "/suppliers/"+pathValue, "id", pathValue, "", GetSupplier)
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

func TestUpdateSupplier(t *testing.T) {
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
				return createTestSupplier(t, "SUP001", "Supplier One")
			},
			body:       `{"code":"SUP001","name":"Supplier Updated"}`,
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Supplier
				decodeJSON(t, w, &resp)
				if resp.Name != "Supplier Updated" {
					t.Errorf("expected name 'Supplier Updated', got '%s'", resp.Name)
				}
			},
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			body:       `{"code":"SUP001","name":"Supplier Updated"}`,
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
			w := executeRequestWithPathValue(t, http.MethodPut, "/suppliers/"+pathValue, "id", pathValue, tt.body, UpdateSupplier)
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

func TestDeleteSupplier(t *testing.T) {
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
				return createTestSupplier(t, "SUP001", "Supplier One")
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
			w := executeRequestWithPathValue(t, http.MethodDelete, "/suppliers/"+pathValue, "id", pathValue, "", DeleteSupplier)
			assertStatus(t, w, tt.wantStatus)
			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
		})
	}
}

func createTestSupplier(t *testing.T, code, name string) string {
	t.Helper()
	supplier := model.Supplier{Code: code, Name: name}
	result := database.DB.Create(&supplier)
	if result.Error != nil {
		t.Fatalf("failed to create test supplier: %v", result.Error)
	}
	return fmt.Sprintf("%d", supplier.ID)
}
