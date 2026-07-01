package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type CustomerFilter struct {
	Search string
}

func GetCustomers(filter CustomerFilter, page, limit int) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := applyCustomerFilter(database.DB.Model(&model.Customer{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func GetAllCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	result := database.DB.Where("is_active = ?", true).Order("id ASC").Find(&customers)
	return customers, result.Error
}

func GetCustomerByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	result := database.DB.First(&customer, id)
	return &customer, result.Error
}

func CreateCustomer(customer *model.Customer) error {
	return database.DB.Create(customer).Error
}

func UpdateCustomer(id uint, customer *model.Customer) error {
	result := database.DB.Model(&model.Customer{}).Where("id = ?", id).Updates(customer)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteCustomer(id uint) error {
	result := database.DB.Delete(&model.Customer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func applyCustomerFilter(query *gorm.DB, filter CustomerFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR email ILIKE ?", search, search, search)
	}
	return query
}
