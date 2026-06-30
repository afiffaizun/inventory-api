package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/afiffazun/inventory-api/internal/model"
	"github.com/afiffazun/inventory-api/internal/repository"
	"github.com/afiffazun/inventory-api/internal/service"
	"gorm.io/gorm"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	respondJSON(w, status, model.ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{Code: code, Message: message},
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, service.GetHome())
}

func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

func Version(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, service.GetVersion())
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	filter := repository.SearchFilter{
		Search:   r.URL.Query().Get("search"),
		Location: r.URL.Query().Get("location"),
	}

	if minStock := r.URL.Query().Get("min_stock"); minStock != "" {
		if v, err := strconv.Atoi(minStock); err == nil {
			filter.MinStock = v
		}
	}

	if maxStock := r.URL.Query().Get("max_stock"); maxStock != "" {
		if v, err := strconv.Atoi(maxStock); err == nil {
			filter.MaxStock = v
		}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	items, total, err := service.GetItemsWithFilter(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch items")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.Item]{
		Data:       items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid item ID")
		return
	}

	item, err := service.GetItemByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Item not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch item")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if item.Code == "" || item.Name == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Code and Name are required")
		return
	}

	if err := service.CreateItem(&item); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create item")
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid item ID")
		return
	}

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	existing, err := service.GetItemByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Item not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch item")
		return
	}

	if item.Code == "" {
		item.Code = existing.Code
	}
	if item.Name == "" {
		item.Name = existing.Name
	}

	if err := service.UpdateItem(uint(id), &item); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update item")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid item ID")
		return
	}

	_, err = service.GetItemByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Item not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch item")
		return
	}

	if err := service.DeleteItem(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ExportItems(w http.ResponseWriter, r *http.Request) {
	filter := repository.SearchFilter{
		Search:   r.URL.Query().Get("search"),
		Location: r.URL.Query().Get("location"),
	}

	if minStock := r.URL.Query().Get("min_stock"); minStock != "" {
		if v, err := strconv.Atoi(minStock); err == nil {
			filter.MinStock = v
		}
	}

	if maxStock := r.URL.Query().Get("max_stock"); maxStock != "" {
		if v, err := strconv.Atoi(maxStock); err == nil {
			filter.MaxStock = v
		}
	}

	items, err := service.ExportItems(filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to export items")
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=items.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Code", "Name", "Stock", "Location"})

	for _, item := range items {
		writer.Write([]string{
			item.Code,
			item.Name,
			strconv.Itoa(item.Stock),
			item.Location,
		})
	}
}
