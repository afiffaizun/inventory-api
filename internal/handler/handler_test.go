package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	Home(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp model.Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Application != "Inventory API" {
		t.Errorf("expected application 'Inventory API', got '%s'", resp.Application)
	}
}

func TestHealth(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("expected status 'UP', got '%s'", resp["status"])
	}
}

func TestVersion(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	Version(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp model.Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Version != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got '%s'", resp.Version)
	}
}

func TestGetItems_Empty(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	GetItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp []model.Item
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp) != 0 {
		t.Errorf("expected empty list, got %d items", len(resp))
	}
}

func TestCreateItem_Success(t *testing.T) {
	setupTest()

	body := `{"code":"ITEM001","name":"Laptop ACER","stock":10,"location":"Warehouse A"}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateItem(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp model.Item
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if resp.Code != "ITEM001" {
		t.Errorf("expected code 'ITEM001', got '%s'", resp.Code)
	}
	if resp.Name != "Laptop ACER" {
		t.Errorf("expected name 'Laptop ACER', got '%s'", resp.Name)
	}
}

func TestCreateItem_ValidationError(t *testing.T) {
	setupTest()

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp model.ErrorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected error code 'VALIDATION_ERROR', got '%s'", resp.Error.Code)
	}
}

func TestCreateItem_InvalidJSON(t *testing.T) {
	setupTest()

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetItem_Success(t *testing.T) {
	setupTest()

	item := model.Item{Code: "ITEM001", Name: "Laptop", Stock: 5, Location: "Gudang A"}
	result := database.DB.Create(&item)
	if result.Error != nil {
		t.Fatalf("failed to create item: %v", result.Error)
	}
	idStr := strconv.FormatUint(uint64(item.ID), 10)

	req := httptest.NewRequest(http.MethodGet, "/items/"+idStr, nil)
	req.SetPathValue("id", idStr)
	w := httptest.NewRecorder()

	GetItem(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp model.Item
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Code != "ITEM001" {
		t.Errorf("expected code 'ITEM001', got '%s'", resp.Code)
	}
}

func TestGetItem_NotFound(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodGet, "/items/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	GetItem(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp model.ErrorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected error code 'NOT_FOUND', got '%s'", resp.Error.Code)
	}
}

func TestGetItem_InvalidID(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodGet, "/items/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	GetItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateItem_Success(t *testing.T) {
	setupTest()

	item := model.Item{Code: "ITEM001", Name: "Laptop", Stock: 5, Location: "Gudang A"}
	result := database.DB.Create(&item)
	if result.Error != nil {
		t.Fatalf("failed to create item: %v", result.Error)
	}
	idStr := strconv.FormatUint(uint64(item.ID), 10)

	body := `{"code":"ITEM001","name":"Laptop Updated","stock":10,"location":"Gudang B"}`
	req := httptest.NewRequest(http.MethodPut, "/items/"+idStr, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", idStr)
	w := httptest.NewRecorder()

	UpdateItem(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp model.Item
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Name != "Laptop Updated" {
		t.Errorf("expected name 'Laptop Updated', got '%s'", resp.Name)
	}
}

func TestUpdateItem_NotFound(t *testing.T) {
	setupTest()

	body := `{"code":"ITEM001","name":"Laptop Updated","stock":10,"location":"Gudang B"}`
	req := httptest.NewRequest(http.MethodPut, "/items/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	UpdateItem(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteItem_Success(t *testing.T) {
	setupTest()

	item := model.Item{Code: "ITEM001", Name: "Laptop", Stock: 5, Location: "Gudang A"}
	result := database.DB.Create(&item)
	if result.Error != nil {
		t.Fatalf("failed to create item: %v", result.Error)
	}
	idStr := strconv.FormatUint(uint64(item.ID), 10)

	req := httptest.NewRequest(http.MethodDelete, "/items/"+idStr, nil)
	req.SetPathValue("id", idStr)
	w := httptest.NewRecorder()

	DeleteItem(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	var count int64
	database.DB.Model(&model.Item{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 items after delete, got %d", count)
	}
}

func TestDeleteItem_NotFound(t *testing.T) {
	setupTest()

	req := httptest.NewRequest(http.MethodDelete, "/items/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	DeleteItem(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
