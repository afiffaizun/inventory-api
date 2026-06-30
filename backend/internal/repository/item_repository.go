package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type SearchFilter struct {
	Search   string
	Location string
	MinStock int
	MaxStock int
}

func GetAllItems() ([]model.Item, error) {
	var items []model.Item
	result := database.DB.Find(&items)
	return items, result.Error
}

func GetItemByID(id uint) (*model.Item, error) {
	var item model.Item
	result := database.DB.First(&item, id)
	return &item, result.Error
}

func CreateItem(item *model.Item) error {
	result := database.DB.Create(item)
	return result.Error
}

func UpdateItem(id uint, item *model.Item) error {
	result := database.DB.Model(&model.Item{}).Where("id = ?", id).Updates(item)
	return result.Error
}

func DeleteItem(id uint) error {
	result := database.DB.Delete(&model.Item{}, id)
	return result.Error
}

func GetItemsWithFilter(filter SearchFilter, page, limit int) ([]model.Item, int64, error) {
	var items []model.Item
	var total int64

	query := applyFilter(database.DB.Model(&model.Item{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func ExportItems(filter SearchFilter) ([]model.Item, error) {
	var items []model.Item

	query := applyFilter(database.DB.Model(&model.Item{}), filter)

	if err := query.Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func applyFilter(query *gorm.DB, filter SearchFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if filter.Location != "" {
		query = query.Where("location = ?", filter.Location)
	}

	if filter.MinStock > 0 {
		query = query.Where("stock >= ?", filter.MinStock)
	}

	if filter.MaxStock > 0 {
		query = query.Where("stock <= ?", filter.MaxStock)
	}

	return query
}
