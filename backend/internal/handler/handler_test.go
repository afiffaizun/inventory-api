package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
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
	database.DB.Exec("DELETE FROM warehouses")
	database.DB.Exec("DELETE FROM categories")
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
			name: "returns items",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
			wantTotal:  2,
		},
		{
			name: "search by name",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?search=Laptop",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "search by code",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?search=ITEM002",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "filter by location",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?location=Gudang+A",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "filter by min stock",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?min_stock=20",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "filter by max stock",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?max_stock=20",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "pagination page 1",
			setup: func(t *testing.T) {
				setupTest()
				for i := 1; i <= 15; i++ {
					createTestItem(t, "ITEM"+strconv.Itoa(i), "Item "+strconv.Itoa(i), i, "Gudang A")
				}
			},
			query:      "?page=1&limit=10",
			wantStatus: http.StatusOK,
			wantCount:  10,
			wantTotal:  15,
		},
		{
			name: "pagination page 2",
			setup: func(t *testing.T) {
				setupTest()
				for i := 1; i <= 15; i++ {
					createTestItem(t, "ITEM"+strconv.Itoa(i), "Item "+strconv.Itoa(i), i, "Gudang A")
				}
			},
			query:      "?page=2&limit=10",
			wantStatus: http.StatusOK,
			wantCount:  5,
			wantTotal:  15,
		},
		{
			name: "invalid page defaults to 1",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
			},
			query:      "?page=0",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "invalid limit defaults to 10",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
			},
			query:      "?limit=0",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "limit over 100 defaults to 10",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
			},
			query:      "?limit=200",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
		{
			name: "combined filters",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
				createTestItem(t, "ITEM003", "Keyboard", 25, "Gudang A")
			},
			query:      "?search=L&location=Gudang+A&min_stock=5&max_stock=30",
			wantStatus: http.StatusOK,
			wantCount:  1,
			wantTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			path := "/items" + tt.query
			w := executeRequestWithQuery(t, http.MethodGet, path, "", GetItems)
			assertStatus(t, w, tt.wantStatus)

			var resp model.PaginatedResponse[model.Item]
			decodeJSON(t, w, &resp)

			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d items, got %d", tt.wantCount, len(resp.Data))
			}

			if resp.Total != tt.wantTotal {
				t.Errorf("expected total %d, got %d", tt.wantTotal, resp.Total)
			}

			if resp.Page == 0 {
				t.Error("expected non-zero page")
			}

			if resp.Limit == 0 {
				t.Error("expected non-zero limit")
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
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				if len(resp.Errors) != 2 {
					t.Errorf("expected 2 errors, got %d", len(resp.Errors))
				}
			},
		},
		{
			name:       "validation error - invalid JSON",
			body:       `invalid json`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
		},
		{
			name:       "validation error - negative stock",
			body:       `{"code":"ITEM001","name":"Laptop","stock":-5}`,
			wantStatus: http.StatusBadRequest,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				if len(resp.Errors) != 1 || resp.Errors[0].Field != "stock" {
					t.Errorf("expected stock error, got %v", resp.Errors)
				}
			},
		},
		{
			name:       "validation error - code too long",
			body:       `{"code":"THIS_CODE_IS_WAY_TOO_LONG_FOR_THE_DATABASE_FIELD_THAT_ONLY_ALLOWS_50_CHARACTERS","name":"Laptop","stock":10}`,
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
			name:       "validation error - code invalid format",
			body:       `{"code":"ITEM@001","name":"Laptop","stock":10}`,
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
			name:       "validation error - multiple errors",
			body:       `{"code":"","name":"","stock":-1}`,
			wantStatus: http.StatusBadRequest,
			checkResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp model.ValidationErrors
				decodeJSON(t, w, &resp)
				if len(resp.Errors) != 3 {
					t.Errorf("expected 3 errors, got %d", len(resp.Errors))
				}
			},
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

func TestExportItems(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		query      string
		wantStatus int
		wantLines  int
		wantType   string
	}{
		{
			name: "empty export",
			setup: func(t *testing.T) {
				setupTest()
			},
			wantStatus: http.StatusOK,
			wantLines:  1,
			wantType:   "text/csv",
		},
		{
			name: "export all items",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			wantStatus: http.StatusOK,
			wantLines:  3,
			wantType:   "text/csv",
		},
		{
			name: "export with search filter",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?search=Laptop",
			wantStatus: http.StatusOK,
			wantLines:  2,
			wantType:   "text/csv",
		},
		{
			name: "export with location filter",
			setup: func(t *testing.T) {
				setupTest()
				createTestItem(t, "ITEM001", "Laptop", 10, "Gudang A")
				createTestItem(t, "ITEM002", "Mouse", 50, "Gudang B")
			},
			query:      "?location=Gudang+B",
			wantStatus: http.StatusOK,
			wantLines:  2,
			wantType:   "text/csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			path := "/items/export" + tt.query
			w := executeRequestWithQuery(t, http.MethodGet, path, "", ExportItems)
			assertStatus(t, w, tt.wantStatus)

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.wantType {
				t.Errorf("expected Content-Type '%s', got '%s'", tt.wantType, contentType)
			}

			contentDisposition := w.Header().Get("Content-Disposition")
			if !strings.Contains(contentDisposition, "attachment") {
				t.Errorf("expected Content-Disposition to contain 'attachment', got '%s'", contentDisposition)
			}

			lines := strings.Split(strings.TrimSpace(w.Body.String()), "\n")
			if len(lines) != tt.wantLines {
				t.Errorf("expected %d lines (header + data), got %d", tt.wantLines, len(lines))
				t.Logf("CSV output:\n%s", w.Body.String())
			}

			if len(lines) > 0 {
				header := lines[0]
				expectedHeader := "Code,Name,Stock,Location"
				if header != expectedHeader {
					t.Errorf("expected header '%s', got '%s'", expectedHeader, header)
				}
			}
		})
	}
}
