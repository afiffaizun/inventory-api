package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func TestGetCustomers(t *testing.T) {
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
			name: "returns customers",
			setup: func(t *testing.T) {
				setupTest()
				createTestCustomer(t, "CUST001", "Customer One")
				createTestCustomer(t, "CUST002", "Customer Two")
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
			wantTotal:  2,
		},
		{
			name: "search by name",
			setup: func(t *testing.T) {
				setupTest()
				createTestCustomer(t, "CUST001", "John Doe")
				createTestCustomer(t, "CUST002", "Jane Smith")
			},
			query:      "?search=John",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			w := executeRequestWithQuery(t, http.MethodGet, "/customers"+tt.query, "", GetCustomers)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.Customer]
			decodeJSON(t, w, &resp)

			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d customers, got %d", tt.wantCount, len(resp.Data))
			}
			if resp.Total != tt.wantTotal {
				t.Errorf("expected total %d, got %d", tt.wantTotal, resp.Total)
			}
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
		checkResp  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "success",
			body:       `{"code":"CUST001","name":"Customer One","email":"cust@test.com","phone":"12345"}`,
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Customer
				decodeJSON(t, w, &resp)
				if resp.ID == 0 {
					t.Error("expected non-zero ID")
				}
				if resp.Code != "CUST001" {
					t.Errorf("expected code 'CUST001', got '%s'", resp.Code)
				}
			},
		},
		{
			name:       "validation error - empty code",
			body:       `{"name":"Customer One"}`,
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
			body:       `{"code":"CUST001"}`,
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
			w := executeRequest(t, http.MethodPost, "/customers", tt.body, CreateCustomer)
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

func TestGetCustomer(t *testing.T) {
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
				return createTestCustomer(t, "CUST001", "Customer One")
			},
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Customer
				decodeJSON(t, w, &resp)
				if resp.Code != "CUST001" {
					t.Errorf("expected code 'CUST001', got '%s'", resp.Code)
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
			w := executeRequestWithPathValue(t, http.MethodGet, "/customers/"+pathValue, "id", pathValue, "", GetCustomer)
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

func TestUpdateCustomer(t *testing.T) {
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
				return createTestCustomer(t, "CUST001", "Customer One")
			},
			body:       `{"code":"CUST001","name":"Customer Updated"}`,
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.Customer
				decodeJSON(t, w, &resp)
				if resp.Name != "Customer Updated" {
					t.Errorf("expected name 'Customer Updated', got '%s'", resp.Name)
				}
			},
		},
		{
			name:       "not found",
			setupID:    func(t *testing.T) string { return "" },
			pathID:     "999",
			body:       `{"code":"CUST001","name":"Customer Updated"}`,
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
			w := executeRequestWithPathValue(t, http.MethodPut, "/customers/"+pathValue, "id", pathValue, tt.body, UpdateCustomer)
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

func TestDeleteCustomer(t *testing.T) {
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
				return createTestCustomer(t, "CUST001", "Customer One")
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
			w := executeRequestWithPathValue(t, http.MethodDelete, "/customers/"+pathValue, "id", pathValue, "", DeleteCustomer)
			assertStatus(t, w, tt.wantStatus)
			if tt.wantCode != "" {
				assertErrorCode(t, w, tt.wantCode)
			}
		})
	}
}

func createTestCustomer(t *testing.T, code, name string) string {
	t.Helper()
	customer := model.Customer{Code: code, Name: name}
	result := database.DB.Create(&customer)
	if result.Error != nil {
		t.Fatalf("failed to create test customer: %v", result.Error)
	}
	return fmt.Sprintf("%d", customer.ID)
}
