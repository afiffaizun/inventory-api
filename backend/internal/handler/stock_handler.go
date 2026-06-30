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

// Stock Movement handlers
func GetStockMovements(w http.ResponseWriter, r *http.Request) {
	filter := repository.StockMovementFilter{
		Type: r.URL.Query().Get("type"),
	}

	if itemID := r.URL.Query().Get("item_id"); itemID != "" {
		if v, err := strconv.ParseUint(itemID, 10, 32); err == nil {
			id := uint(v)
			filter.ItemID = &id
		}
	}

	if warehouseID := r.URL.Query().Get("warehouse_id"); warehouseID != "" {
		if v, err := strconv.ParseUint(warehouseID, 10, 32); err == nil {
			id := uint(v)
			filter.WarehouseID = &id
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

	movements, total, err := service.GetStockMovements(filter, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock movements")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.StockMovement]{
		Data:       movements,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetStockMovement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid stock movement ID")
		return
	}

	movement, err := service.GetStockMovementByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Stock movement not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock movement")
		return
	}
	respondJSON(w, http.StatusOK, movement)
}

func CreateStockMovement(w http.ResponseWriter, r *http.Request) {
	var movement model.StockMovement

	if err := json.NewDecoder(r.Body).Decode(&movement); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := movement.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	if err := service.CreateStockMovement(&movement); err != nil {
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			respondError(w, http.StatusBadRequest, "INSUFFICIENT_STOCK", "Insufficient stock for this operation")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create stock movement")
		return
	}
	respondJSON(w, http.StatusCreated, movement)
}

func TransferStock(w http.ResponseWriter, r *http.Request) {
	var req model.TransferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		respondJSON(w, http.StatusBadRequest, model.ValidationErrors{Errors: errs})
		return
	}

	outMovement := model.StockMovement{
		ItemID:        req.ItemID,
		WarehouseID:   req.FromWarehouseID,
		Type:          model.MovementTypeOUT,
		Quantity:      req.Quantity,
		ReferenceType: model.ReferenceTypeTRANSFER,
		Notes:         req.Notes,
		CreatedBy:     req.CreatedBy,
	}

	if err := service.CreateStockMovement(&outMovement); err != nil {
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			respondError(w, http.StatusBadRequest, "INSUFFICIENT_STOCK", "Insufficient stock for this transfer")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process transfer (OUT)")
		return
	}

	inMovement := model.StockMovement{
		ItemID:        req.ItemID,
		WarehouseID:   req.ToWarehouseID,
		Type:          model.MovementTypeIN,
		Quantity:      req.Quantity,
		ReferenceType: model.ReferenceTypeTRANSFER,
		Notes:         req.Notes,
		CreatedBy:     req.CreatedBy,
	}

	if err := service.CreateStockMovement(&inMovement); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process transfer (IN)")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Transfer completed successfully"})
}

func GetStockHistory(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid item ID")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movements, total, err := service.GetStockHistory(uint(itemID), page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock history")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.StockMovement]{
		Data:       movements,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetStockByWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid warehouse ID")
		return
	}

	items, err := service.GetStockByWarehouse(uint(warehouseID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock by warehouse")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func GetStockSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := service.GetStockSummary()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock summary")
		return
	}
	respondJSON(w, http.StatusOK, summary)
}

// Stock Opname handlers
func GetStockOpnames(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	opnames, total, err := service.GetStockOpnames(page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock opnames")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	respondJSON(w, http.StatusOK, model.PaginatedResponse[model.StockOpname]{
		Data:       opnames,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func GetStockOpname(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid stock opname ID")
		return
	}

	opname, err := service.GetStockOpnameByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Stock opname not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock opname")
		return
	}
	respondJSON(w, http.StatusOK, opname)
}

func CreateStockOpname(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOpnameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if req.WarehouseID == 0 {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "WarehouseID is required")
		return
	}

	opname := model.StockOpname{
		WarehouseID: req.WarehouseID,
		Status:      model.OpnameStatusDRAFT,
		Notes:       req.Notes,
		CreatedBy:   req.CreatedBy,
	}

	for _, item := range req.Items {
		opnameItem := model.StockOpnameItem{
			ItemID:         item.ItemID,
			ActualQuantity: item.ActualQuantity,
		}

		var dbItem *model.Item
		dbItem, err := service.GetItemByID(item.ItemID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Item not found")
			return
		}
		opnameItem.SystemQuantity = dbItem.Stock
		opnameItem.CalculateDifference()

		opname.Items = append(opname.Items, opnameItem)
	}

	if err := service.CreateStockOpname(&opname); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create stock opname")
		return
	}
	respondJSON(w, http.StatusCreated, opname)
}

func CompleteStockOpname(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid stock opname ID")
		return
	}

	_, err = service.GetStockOpnameByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Stock opname not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock opname")
		return
	}

	if err := service.CompleteStockOpname(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			respondError(w, http.StatusBadRequest, "INVALID_STATUS", "Stock opname cannot be completed in current status")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to complete stock opname")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Stock opname completed successfully"})
}
