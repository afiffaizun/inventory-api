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

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	filter := repository.CustomerFilter{
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

	customers, total, err := service.GetCustomers(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch customers")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.Customer]{
		Data:       customers,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid customer ID")
		return
	}

	customer, err := service.GetCustomerByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Customer not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch customer")
		return
	}
	respondJSON(w, http.StatusOK, customer)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := customer.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.CreateCustomer(&customer); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create customer")
		return
	}
	respondJSON(w, http.StatusCreated, customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid customer ID")
		return
	}

	existing, err := service.GetCustomerByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Customer not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch customer")
		return
	}

	var input model.Customer
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

	if err := service.UpdateCustomer(uint(id), existing); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update customer")
		return
	}
	respondJSON(w, http.StatusOK, existing)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid customer ID")
		return
	}

	if err := service.DeleteCustomer(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Customer not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete customer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
