package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

func executeRequest(t *testing.T, method, path, body string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	} else {
		reqBody = bytes.NewBufferString("")
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	handler(w, req)

	return w
}

func executeRequestWithPathValue(t *testing.T, method, path, pathKey, pathValue, body string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	} else {
		reqBody = bytes.NewBufferString("")
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetPathValue(pathKey, pathValue)

	w := httptest.NewRecorder()
	handler(w, req)

	return w
}

func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	if err := json.NewDecoder(w.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func createTestItem(t *testing.T, code, name string, stock int, location string) string {
	t.Helper()

	item := model.Item{
		Code:     code,
		Name:     name,
		Stock:    stock,
		Location: location,
	}
	result := database.DB.Create(&item)
	if result.Error != nil {
		t.Fatalf("failed to create test item: %v", result.Error)
	}
	return strconv.FormatUint(uint64(item.ID), 10)
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if w.Code != expected {
		t.Errorf("expected status %d, got %d", expected, w.Code)
	}
}

func assertErrorCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode string) {
	t.Helper()

	var resp model.ErrorResponse
	decodeJSON(t, w, &resp)

	if resp.Error.Code != expectedCode {
		t.Errorf("expected error code '%s', got '%s'", expectedCode, resp.Error.Code)
	}
}

func executeRequestWithQuery(t *testing.T, method, path, body string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	} else {
		reqBody = bytes.NewBufferString("")
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	handler(w, req)

	return w
}
