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

func GetCategories(w http.ResponseWriter, r *http.Request) {
	filter := repository.CategoryFilter{
		Search: r.URL.Query().Get("search"),
	}

	if parentID := r.URL.Query().Get("parent_id"); parentID != "" {
		if v, err := strconv.ParseUint(parentID, 10, 32); err == nil {
			pid := uint(v)
			filter.ParentID = &pid
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

	categories, total, err := service.GetCategories(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch categories")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.Category]{
		Data:       categories,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := service.GetAllCategories()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch categories")
		return
	}
	respondJSON(w, http.StatusOK, categories)
}

func GetCategoryTree(w http.ResponseWriter, r *http.Request) {
	tree, err := service.GetCategoryTree()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch category tree")
		return
	}
	respondJSON(w, http.StatusOK, tree)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid category ID")
		return
	}

	category, err := service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch category")
		return
	}
	respondJSON(w, http.StatusOK, category)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := category.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.CreateCategory(&category); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create category")
		return
	}
	respondJSON(w, http.StatusCreated, category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid category ID")
		return
	}

	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	existing, err := service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch category")
		return
	}

	if category.Code == "" {
		category.Code = existing.Code
	}
	if category.Name == "" {
		category.Name = existing.Name
	}

	if errs := category.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.UpdateCategory(uint(id), &category); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update category")
		return
	}
	respondJSON(w, http.StatusOK, category)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid category ID")
		return
	}

	_, err = service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch category")
		return
	}

	if err := service.DeleteCategory(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete category")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
