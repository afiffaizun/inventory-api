package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type SupplierFilter struct {
	Search string
}

func GetSuppliers(filter SupplierFilter, page, limit int) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := applySupplierFilter(database.DB.Model(&model.Supplier{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func GetAllSuppliers() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	result := database.DB.Where("is_active = ?", true).Order("id ASC").Find(&suppliers)
	return suppliers, result.Error
}

func GetSupplierByID(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	result := database.DB.First(&supplier, id)
	return &supplier, result.Error
}

func CreateSupplier(supplier *model.Supplier) error {
	return database.DB.Create(supplier).Error
}

func UpdateSupplier(id uint, supplier *model.Supplier) error {
	result := database.DB.Model(&model.Supplier{}).Where("id = ?", id).Updates(supplier)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteSupplier(id uint) error {
	result := database.DB.Delete(&model.Supplier{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func applySupplierFilter(query *gorm.DB, filter SupplierFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR email ILIKE ?", search, search, search)
	}
	return query
}
