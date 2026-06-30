package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type WarehouseFilter struct {
	Search string
}

func GetWarehouses(filter WarehouseFilter, page, limit int) ([]model.Warehouse, int64, error) {
	var warehouses []model.Warehouse
	var total int64

	query := applyWarehouseFilter(database.DB.Model(&model.Warehouse{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func GetAllWarehouses() ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	result := database.DB.Where("is_active = ?", true).Order("id ASC").Find(&warehouses)
	return warehouses, result.Error
}

func GetWarehouseByID(id uint) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	result := database.DB.First(&warehouse, id)
	return &warehouse, result.Error
}

func CreateWarehouse(warehouse *model.Warehouse) error {
	result := database.DB.Create(warehouse)
	return result.Error
}

func UpdateWarehouse(id uint, warehouse *model.Warehouse) error {
	result := database.DB.Model(&model.Warehouse{}).Where("id = ?", id).Updates(warehouse)
	return result.Error
}

func DeleteWarehouse(id uint) error {
	result := database.DB.Delete(&model.Warehouse{}, id)
	return result.Error
}

func SetDefaultWarehouse(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Warehouse{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&model.Warehouse{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

func applyWarehouseFilter(query *gorm.DB, filter WarehouseFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}
	return query
}
