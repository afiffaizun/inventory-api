package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/afiffazun/inventory-api/internal/model"
	"github.com/afiffazun/inventory-api/internal/repository"
	"github.com/afiffazun/inventory-api/internal/service"
	"gorm.io/gorm"
)

func GetPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	filter := repository.PurchaseOrderFilter{
		Search: r.URL.Query().Get("search"),
		Status: r.URL.Query().Get("status"),
	}

	if supplierID := r.URL.Query().Get("supplier_id"); supplierID != "" {
		if v, err := strconv.ParseUint(supplierID, 10, 32); err == nil {
			id := uint(v)
			filter.SupplierID = &id
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

	orders, total, err := service.GetPurchaseOrders(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch purchase orders")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.PurchaseOrder]{
		Data:       orders,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid purchase order ID")
		return
	}

	order, err := service.GetPurchaseOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Purchase order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch purchase order")
		return
	}
	respondJSON(w, http.StatusOK, order)
}

func CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePurchaseOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid order_date format (use YYYY-MM-DD)")
		return
	}

	order := model.PurchaseOrder{
		SupplierID:  req.SupplierID,
		WarehouseID: req.WarehouseID,
		OrderDate:   orderDate,
		Status:      model.PurchaseOrderStatusDRAFT,
		Notes:       req.Notes,
		CreatedBy:   req.CreatedBy,
	}

	if errs := order.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	for _, itemReq := range req.Items {
		item := model.PurchaseOrderItem{
			ItemID:   itemReq.ItemID,
			Quantity: itemReq.Quantity,
			UnitCost: itemReq.UnitCost,
		}

		if itemReq.UnitCost == 0 {
			dbItem, err := service.GetItemByID(itemReq.ItemID)
			if err != nil {
				respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Item %d not found", itemReq.ItemID))
				return
			}
			item.UnitCost = dbItem.CostPrice
		}

		item.CalculateSubtotal()

		if errs := item.Validate(); len(errs) > 0 {
			respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
			return
		}

		order.Items = append(order.Items, item)
	}

	if err := service.CreatePurchaseOrder(&order); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create purchase order")
		return
	}
	respondJSON(w, http.StatusCreated, order)
}

func UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid purchase order ID")
		return
	}

	existing, err := service.GetPurchaseOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Purchase order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch purchase order")
		return
	}

	if existing.Status != model.PurchaseOrderStatusDRAFT {
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", "Only draft orders can be updated")
		return
	}

	var input model.PurchaseOrder
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if input.Notes != "" {
		existing.Notes = input.Notes
	}

	if err := service.UpdatePurchaseOrder(uint(id), existing); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update purchase order")
		return
	}
	respondJSON(w, http.StatusOK, existing)
}

func DeletePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid purchase order ID")
		return
	}

	existing, err := service.GetPurchaseOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Purchase order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch purchase order")
		return
	}

	if existing.Status != model.PurchaseOrderStatusDRAFT {
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", "Only draft orders can be deleted")
		return
	}

	if err := service.DeletePurchaseOrder(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete purchase order")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ConfirmPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid purchase order ID")
		return
	}

	if err := service.ConfirmPurchaseOrder(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Purchase order not found")
			return
		}
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Purchase order confirmed successfully"})
}

func CancelPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid purchase order ID")
		return
	}

	if err := service.CancelPurchaseOrder(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Purchase order not found")
			return
		}
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Purchase order cancelled successfully"})
}
