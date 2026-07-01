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

// Customer functions
func GetCustomers(filter repository.CustomerFilter, page, limit int) ([]model.Customer, int64, error) {
	return repository.GetCustomers(filter, page, limit)
}

func GetAllCustomers() ([]model.Customer, error) {
	return repository.GetAllCustomers()
}

func GetCustomerByID(id uint) (*model.Customer, error) {
	return repository.GetCustomerByID(id)
}

func CreateCustomer(customer *model.Customer) error {
	return repository.CreateCustomer(customer)
}

func UpdateCustomer(id uint, customer *model.Customer) error {
	return repository.UpdateCustomer(id, customer)
}

func DeleteCustomer(id uint) error {
	return repository.DeleteCustomer(id)
}

// Supplier functions
func GetSuppliers(filter repository.SupplierFilter, page, limit int) ([]model.Supplier, int64, error) {
	return repository.GetSuppliers(filter, page, limit)
}

func GetAllSuppliers() ([]model.Supplier, error) {
	return repository.GetAllSuppliers()
}

func GetSupplierByID(id uint) (*model.Supplier, error) {
	return repository.GetSupplierByID(id)
}

func CreateSupplier(supplier *model.Supplier) error {
	return repository.CreateSupplier(supplier)
}

func UpdateSupplier(id uint, supplier *model.Supplier) error {
	return repository.UpdateSupplier(id, supplier)
}

func DeleteSupplier(id uint) error {
	return repository.DeleteSupplier(id)
}

// Sales Order functions
func GetSalesOrders(filter repository.SalesOrderFilter, page, limit int) ([]model.SalesOrder, int64, error) {
	return repository.GetSalesOrders(filter, page, limit)
}

func GetSalesOrderByID(id uint) (*model.SalesOrder, error) {
	return repository.GetSalesOrderByID(id)
}

func CreateSalesOrder(order *model.SalesOrder) error {
	return repository.CreateSalesOrder(order)
}

func UpdateSalesOrder(id uint, order *model.SalesOrder) error {
	return repository.UpdateSalesOrder(id, order)
}

func DeleteSalesOrder(id uint) error {
	return repository.DeleteSalesOrder(id)
}

func ConfirmSalesOrder(id uint) error {
	return repository.ConfirmSalesOrder(id)
}

func CancelSalesOrder(id uint) error {
	return repository.CancelSalesOrder(id)
}

// Purchase Order functions
func GetPurchaseOrders(filter repository.PurchaseOrderFilter, page, limit int) ([]model.PurchaseOrder, int64, error) {
	return repository.GetPurchaseOrders(filter, page, limit)
}

func GetPurchaseOrderByID(id uint) (*model.PurchaseOrder, error) {
	return repository.GetPurchaseOrderByID(id)
}

func CreatePurchaseOrder(order *model.PurchaseOrder) error {
	return repository.CreatePurchaseOrder(order)
}

func UpdatePurchaseOrder(id uint, order *model.PurchaseOrder) error {
	return repository.UpdatePurchaseOrder(id, order)
}

func DeletePurchaseOrder(id uint) error {
	return repository.DeletePurchaseOrder(id)
}

func ConfirmPurchaseOrder(id uint) error {
	return repository.ConfirmPurchaseOrder(id)
}

func CancelPurchaseOrder(id uint) error {
	return repository.CancelPurchaseOrder(id)
}
