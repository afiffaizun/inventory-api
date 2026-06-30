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

// Warehouse functions
func GetWarehouses(filter repository.WarehouseFilter, page, limit int) ([]model.Warehouse, int64, error) {
	return repository.GetWarehouses(filter, page, limit)
}

func GetAllWarehouses() ([]model.Warehouse, error) {
	return repository.GetAllWarehouses()
}

func GetWarehouseByID(id uint) (*model.Warehouse, error) {
	return repository.GetWarehouseByID(id)
}

func CreateWarehouse(warehouse *model.Warehouse) error {
	return repository.CreateWarehouse(warehouse)
}

func UpdateWarehouse(id uint, warehouse *model.Warehouse) error {
	return repository.UpdateWarehouse(id, warehouse)
}

func DeleteWarehouse(id uint) error {
	return repository.DeleteWarehouse(id)
}

func SetDefaultWarehouse(id uint) error {
	return repository.SetDefaultWarehouse(id)
}

// Category functions
func GetCategories(filter repository.CategoryFilter, page, limit int) ([]model.Category, int64, error) {
	return repository.GetCategories(filter, page, limit)
}

func GetAllCategories() ([]model.Category, error) {
	return repository.GetAllCategories()
}

func GetCategoryTree() ([]model.CategoryTree, error) {
	return repository.GetCategoryTree()
}

func GetCategoryByID(id uint) (*model.Category, error) {
	return repository.GetCategoryByID(id)
}

func CreateCategory(category *model.Category) error {
	return repository.CreateCategory(category)
}

func UpdateCategory(id uint, category *model.Category) error {
	return repository.UpdateCategory(id, category)
}

func DeleteCategory(id uint) error {
	return repository.DeleteCategory(id)
}

// Stock Movement functions
func GetStockMovements(filter repository.StockMovementFilter, page, limit int) ([]model.StockMovement, int64, error) {
	return repository.GetStockMovements(filter, page, limit)
}

func GetStockMovementByID(id uint) (*model.StockMovement, error) {
	return repository.GetStockMovementByID(id)
}

func CreateStockMovement(movement *model.StockMovement) error {
	return repository.CreateStockMovement(movement)
}

func GetStockHistory(itemID uint, page, limit int) ([]model.StockMovement, int64, error) {
	return repository.GetStockHistory(itemID, page, limit)
}

func GetStockByWarehouse(warehouseID uint) ([]model.Item, error) {
	return repository.GetStockByWarehouse(warehouseID)
}

func GetStockSummary() ([]repository.StockSummary, error) {
	return repository.GetStockSummary()
}

// Stock Opname functions
func GetStockOpnames(page, limit int) ([]model.StockOpname, int64, error) {
	return repository.GetStockOpnames(page, limit)
}

func GetStockOpnameByID(id uint) (*model.StockOpname, error) {
	return repository.GetStockOpnameByID(id)
}

func CreateStockOpname(opname *model.StockOpname) error {
	return repository.CreateStockOpname(opname)
}

func CompleteStockOpname(id uint) error {
	return repository.CompleteStockOpname(id)
}
