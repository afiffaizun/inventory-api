package service

import (
	"github.com/afiffazun/inventory-api/internal/model"
	"github.com/afiffazun/inventory-api/internal/repository"
)

func GetHome() model.Response {
	return model.Response{
		Application: "Inventory API",
		Author:      "Afif",
		Status:      "Running",
	}
}

func GetVersion() model.Response {
	return model.Response{
		Version: "v1.0.0",
	}
}

func GetAllItems() ([]model.Item, error) {
	return repository.GetAllItems()
}

func GetItemByID(id uint) (*model.Item, error) {
	return repository.GetItemByID(id)
}

func CreateItem(item *model.Item) error {
	return repository.CreateItem(item)
}

func UpdateItem(id uint, item *model.Item) error {
	return repository.UpdateItem(id, item)
}

func DeleteItem(id uint) error {
	return repository.DeleteItem(id)
}

func GetItemsWithFilter(filter repository.SearchFilter, page, limit int) ([]model.Item, int64, error) {
	return repository.GetItemsWithFilter(filter, page, limit)
}

func ExportItems(filter repository.SearchFilter) ([]model.Item, error) {
	return repository.ExportItems(filter)
}
