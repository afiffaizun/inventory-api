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

func GetSalesOrders(w http.ResponseWriter, r *http.Request) {
	filter := repository.SalesOrderFilter{
		Search: r.URL.Query().Get("search"),
		Status: r.URL.Query().Get("status"),
	}

	if customerID := r.URL.Query().Get("customer_id"); customerID != "" {
		if v, err := strconv.ParseUint(customerID, 10, 32); err == nil {
			id := uint(v)
			filter.CustomerID = &id
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

	orders, total, err := service.GetSalesOrders(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch sales orders")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.SalesOrder]{
		Data:       orders,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid sales order ID")
		return
	}

	order, err := service.GetSalesOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Sales order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch sales order")
		return
	}
	respondJSON(w, http.StatusOK, order)
}

func CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSalesOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid order_date format (use YYYY-MM-DD)")
		return
	}

	order := model.SalesOrder{
		CustomerID:  req.CustomerID,
		WarehouseID: req.WarehouseID,
		OrderDate:   orderDate,
		Status:      model.SalesOrderStatusDRAFT,
		Notes:       req.Notes,
		CreatedBy:   req.CreatedBy,
	}

	if errs := order.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	for _, itemReq := range req.Items {
		item := model.SalesOrderItem{
			ItemID:    itemReq.ItemID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
		}

		if itemReq.UnitPrice == 0 {
			dbItem, err := service.GetItemByID(itemReq.ItemID)
			if err != nil {
				respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Item %d not found", itemReq.ItemID))
				return
			}
			item.UnitPrice = dbItem.SellPrice
		}

		item.CalculateSubtotal()

		if errs := item.Validate(); len(errs) > 0 {
			respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
			return
		}

		order.Items = append(order.Items, item)
	}

	if err := service.CreateSalesOrder(&order); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create sales order")
		return
	}
	respondJSON(w, http.StatusCreated, order)
}

func UpdateSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid sales order ID")
		return
	}

	existing, err := service.GetSalesOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Sales order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch sales order")
		return
	}

	if existing.Status != model.SalesOrderStatusDRAFT {
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", "Only draft orders can be updated")
		return
	}

	var input model.SalesOrder
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if input.Notes != "" {
		existing.Notes = input.Notes
	}

	if err := service.UpdateSalesOrder(uint(id), existing); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update sales order")
		return
	}
	respondJSON(w, http.StatusOK, existing)
}

func DeleteSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid sales order ID")
		return
	}

	existing, err := service.GetSalesOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Sales order not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch sales order")
		return
	}

	if existing.Status != model.SalesOrderStatusDRAFT {
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", "Only draft orders can be deleted")
		return
	}

	if err := service.DeleteSalesOrder(uint(id)); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete sales order")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ConfirmSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid sales order ID")
		return
	}

	if err := service.ConfirmSalesOrder(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Sales order not found")
			return
		}
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Sales order confirmed successfully"})
}

func CancelSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid sales order ID")
		return
	}

	if err := service.CancelSalesOrder(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Sales order not found")
			return
		}
		respondError(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Sales order cancelled successfully"})
}
