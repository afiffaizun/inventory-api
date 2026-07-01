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

func GetSuppliers(w http.ResponseWriter, r *http.Request) {
	filter := repository.SupplierFilter{
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

	suppliers, total, err := service.GetSuppliers(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch suppliers")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.Supplier]{
		Data:       suppliers,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid supplier ID")
		return
	}

	supplier, err := service.GetSupplierByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Supplier not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch supplier")
		return
	}
	respondJSON(w, http.StatusOK, supplier)
}

func CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var supplier model.Supplier

	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := supplier.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.CreateSupplier(&supplier); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create supplier")
		return
	}
	respondJSON(w, http.StatusCreated, supplier)
}

func UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid supplier ID")
		return
	}

	existing, err := service.GetSupplierByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Supplier not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch supplier")
		return
	}

	var input model.Supplier
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if input.Code != "" {
		existing.Code = input.Code
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.Email != "" {
		existing.Email = input.Email
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}
	if input.Address != "" {
		existing.Address = input.Address
	}
	if input.City != "" {
		existing.City = input.City
	}
	if input.Country != "" {
		existing.Country = input.Country
	}

	if errs := existing.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.UpdateSupplier(uint(id), existing); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update supplier")
		return
	}
	respondJSON(w, http.StatusOK, existing)
}

func DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid supplier ID")
		return
	}

	if err := service.DeleteSupplier(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Supplier not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete supplier")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
