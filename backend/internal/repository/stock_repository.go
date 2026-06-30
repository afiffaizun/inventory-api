package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type StockMovementFilter struct {
	ItemID      *uint
	WarehouseID *uint
	Type        string
}

func GetStockMovements(filter StockMovementFilter, page, limit int) ([]model.StockMovement, int64, error) {
	var movements []model.StockMovement
	var total int64

	query := applyStockMovementFilter(database.DB.Model(&model.StockMovement{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Item").Preload("Warehouse").Offset(offset).Limit(limit).Order("id DESC").Find(&movements).Error; err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

func GetStockMovementByID(id uint) (*model.StockMovement, error) {
	var movement model.StockMovement
	result := database.DB.Preload("Item").Preload("Warehouse").First(&movement, id)
	return &movement, result.Error
}

func CreateStockMovement(movement *model.StockMovement) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		var item model.Item
		if err := tx.First(&item, movement.ItemID).Error; err != nil {
			return err
		}

		switch movement.Type {
		case model.MovementTypeIN:
			item.Stock += movement.Quantity
		case model.MovementTypeOUT:
			if item.Stock < movement.Quantity {
				return gorm.ErrInvalidTransaction
			}
			item.Stock -= movement.Quantity
		case model.MovementTypeADJUSTMENT:
			item.Stock += movement.Quantity
		}

		return tx.Model(&item).Update("stock", item.Stock).Error
	})
}

func GetStockHistory(itemID uint, page, limit int) ([]model.StockMovement, int64, error) {
	var movements []model.StockMovement
	var total int64

	query := database.DB.Model(&model.StockMovement{}).Where("item_id = ?", itemID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Warehouse").Offset(offset).Limit(limit).Order("id DESC").Find(&movements).Error; err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

func GetStockByWarehouse(warehouseID uint) ([]model.Item, error) {
	var items []model.Item
	result := database.DB.Where("warehouse_id = ? AND is_active = ?", warehouseID, true).Order("id ASC").Find(&items)
	return items, result.Error
}

func GetStockSummary() ([]StockSummary, error) {
	type StockSummaryResult struct {
		WarehouseID   uint
		WarehouseName string
		TotalItems    int64
		TotalStock    int64
	}

	var results []StockSummaryResult
	err := database.DB.Model(&model.Item{}).
		Select("items.warehouse_id, warehouses.name as warehouse_name, count(items.id) as total_items, sum(items.stock) as total_stock").
		Joins("LEFT JOIN warehouses ON warehouses.id = items.warehouse_id").
		Where("items.is_active = ?", true).
		Group("items.warehouse_id, warehouses.name").
		Order("items.warehouse_id ASC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	var summaries []StockSummary
	for _, r := range results {
		summaries = append(summaries, StockSummary{
			WarehouseID:   r.WarehouseID,
			WarehouseName: r.WarehouseName,
			TotalItems:    r.TotalItems,
			TotalStock:    r.TotalStock,
		})
	}

	return summaries, nil
}

type StockSummary struct {
	WarehouseID   uint   `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	TotalItems    int64  `json:"total_items"`
	TotalStock    int64  `json:"total_stock"`
}

// StockOpname functions
func GetStockOpnames(page, limit int) ([]model.StockOpname, int64, error) {
	var opnames []model.StockOpname
	var total int64

	query := database.DB.Model(&model.StockOpname{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Warehouse").Preload("Items").Preload("Items.Item").Offset(offset).Limit(limit).Order("id DESC").Find(&opnames).Error; err != nil {
		return nil, 0, err
	}

	return opnames, total, nil
}

func GetStockOpnameByID(id uint) (*model.StockOpname, error) {
	var opname model.StockOpname
	result := database.DB.Preload("Warehouse").Preload("Items").Preload("Items.Item").First(&opname, id)
	return &opname, result.Error
}

func CreateStockOpname(opname *model.StockOpname) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		items := opname.Items
		opname.Items = nil

		if err := tx.Create(opname).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].StockOpnameID = opname.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		opname.Items = items
		return nil
	})
}

func CompleteStockOpname(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var opname model.StockOpname
		if err := tx.Preload("Items").First(&opname, id).Error; err != nil {
			return err
		}

		if opname.Status != model.OpnameStatusDRAFT && opname.Status != model.OpnameStatusINPROGRESS {
			return gorm.ErrInvalidTransaction
		}

		totalDiff := 0
		for i := range opname.Items {
			opname.Items[i].CalculateDifference()
			totalDiff += opname.Items[i].Difference
			if err := tx.Save(&opname.Items[i]).Error; err != nil {
				return err
			}

			var item model.Item
			if err := tx.First(&item, opname.Items[i].ItemID).Error; err != nil {
				return err
			}

			item.Stock += opname.Items[i].Difference
			if err := tx.Model(&item).Update("stock", item.Stock).Error; err != nil {
				return err
			}
		}

		opname.Status = model.OpnameStatusCOMPLETED
		opname.TotalDiff = totalDiff
		return tx.Save(&opname).Error
	})
}

func applyStockMovementFilter(query *gorm.DB, filter StockMovementFilter) *gorm.DB {
	if filter.ItemID != nil {
		query = query.Where("item_id = ?", *filter.ItemID)
	}

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	return query
}
