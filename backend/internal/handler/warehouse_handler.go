package handler

import (
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

func GetWarehouses(w http.ResponseWriter, r *http.Request) {
	filter := repository.WarehouseFilter{
		Search: r.URL.Query().Get("search"),
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	warehouses, total, err := service.GetWarehouses(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouses")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.Warehouse]{
		Data:       warehouses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetAllWarehouses(w http.ResponseWriter, r *http.Request) {
	warehouses, err := service.GetAllWarehouses()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouses")
		return
	}
	respondJSON(w, http.StatusOK, warehouses)
}

func GetWarehouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid warehouse ID")
		return
	}

	warehouse, err := service.GetWarehouseByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Warehouse not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouse")
		return
	}
	respondJSON(w, http.StatusOK, warehouse)
}

func CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var warehouse model.Warehouse

	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := warehouse.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.CreateWarehouse(&warehouse); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create warehouse")
		return
	}
	respondJSON(w, http.StatusCreated, warehouse)
}

func UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid warehouse ID")
		return
	}

	var warehouse model.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	existing, err := service.GetWarehouseByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Warehouse not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouse")
		return
	}

	if warehouse.Code == "" {
		warehouse.Code = existing.Code
	}
	if warehouse.Name == "" {
		warehouse.Name = existing.Name
	}

	if errs := warehouse.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.UpdateWarehouse(uint(id), &warehouse); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update warehouse")
		return
	}
	respondJSON(w, http.StatusOK, warehouse)
}

func DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid warehouse ID")
		return
	}

	_, err = service.GetWarehouseByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Warehouse not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouse")
		return
	}

	if err := service.DeleteWarehouse(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete warehouse")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func SetDefaultWarehouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid warehouse ID")
		return
	}

	_, err = service.GetWarehouseByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Warehouse not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch warehouse")
		return
	}

	if err := service.SetDefaultWarehouse(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to set default warehouse")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Default warehouse updated"})
}
