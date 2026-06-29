package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
)

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
