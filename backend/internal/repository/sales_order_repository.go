package repository

import (
	"fmt"
	"time"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type SalesOrderFilter struct {
	Search     string
	CustomerID *uint
	Status     string
}

func GetSalesOrders(filter SalesOrderFilter, page, limit int) ([]model.SalesOrder, int64, error) {
	var orders []model.SalesOrder
	var total int64

	query := applySalesOrderFilter(database.DB.Model(&model.SalesOrder{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Customer").Preload("Warehouse").Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetSalesOrderByID(id uint) (*model.SalesOrder, error) {
	var order model.SalesOrder
	result := database.DB.Preload("Customer").Preload("Warehouse").Preload("Items").Preload("Items.Item").First(&order, id)
	return &order, result.Error
}

func CreateSalesOrder(order *model.SalesOrder) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		order.OrderNumber = generateOrderNumber("SO")
		items := order.Items
		order.Items = nil

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].SalesOrderID = order.ID
			items[i].CalculateSubtotal()
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		order.Items = items
		return recalculateSalesOrderTotal(order.ID)
	})
}

func UpdateSalesOrder(id uint, order *model.SalesOrder) error {
	result := database.DB.Model(&model.SalesOrder{}).Where("id = ?", id).Updates(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteSalesOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("sales_order_id = ?", id).Delete(&model.SalesOrderItem{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.SalesOrder{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func ConfirmSalesOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.SalesOrder
		if err := tx.Preload("Items").First(&order, id).Error; err != nil {
			return err
		}

		if order.Status != model.SalesOrderStatusDRAFT {
			return fmt.Errorf("order cannot be confirmed in status: %s", order.Status)
		}

		for _, item := range order.Items {
			var dbItem model.Item
			if err := tx.First(&dbItem, item.ItemID).Error; err != nil {
				return err
			}

			if dbItem.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for item %s: have %d, need %d", dbItem.Name, dbItem.Stock, item.Quantity)
			}

			dbItem.Stock -= item.Quantity
			if err := tx.Model(&dbItem).Update("stock", dbItem.Stock).Error; err != nil {
				return err
			}

			movement := model.StockMovement{
				ItemID:        item.ItemID,
				WarehouseID:   order.WarehouseID,
				Type:          model.MovementTypeOUT,
				Quantity:      item.Quantity,
				ReferenceType: model.ReferenceTypeSALE,
				ReferenceID:   order.ID,
				Notes:         fmt.Sprintf("Sales Order %s", order.OrderNumber),
				CreatedBy:     order.CreatedBy,
			}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		order.Status = model.SalesOrderStatusCONFIRMED
		return tx.Save(&order).Error
	})
}

func CancelSalesOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.SalesOrder
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}

		if order.Status == model.SalesOrderStatusCOMPLETED {
			return fmt.Errorf("cannot cancel completed order")
		}

		if order.Status == model.SalesOrderStatusCONFIRMED {
			if err := reverseSalesOrderStock(tx, &order); err != nil {
				return err
			}
		}

		order.Status = model.SalesOrderStatusCANCELLED
		return tx.Save(&order).Error
	})
}

func reverseSalesOrderStock(tx *gorm.DB, order *model.SalesOrder) error {
	var items []model.SalesOrderItem
	if err := tx.Where("sales_order_id = ?", order.ID).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		var dbItem model.Item
		if err := tx.First(&dbItem, item.ItemID).Error; err != nil {
			return err
		}

		dbItem.Stock += item.Quantity
		if err := tx.Model(&dbItem).Update("stock", dbItem.Stock).Error; err != nil {
			return err
		}

		movement := model.StockMovement{
			ItemID:        item.ItemID,
			WarehouseID:   order.WarehouseID,
			Type:          model.MovementTypeIN,
			Quantity:      item.Quantity,
			ReferenceType: model.ReferenceTypeSALE,
			ReferenceID:   order.ID,
			Notes:         fmt.Sprintf("Cancelled Sales Order %s", order.OrderNumber),
			CreatedBy:     order.CreatedBy,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
	}

	return nil
}

func recalculateSalesOrderTotal(orderID uint) error {
	var items []model.SalesOrderItem
	if err := database.DB.Where("sales_order_id = ?", orderID).Find(&items).Error; err != nil {
		return err
	}

	var total float64
	for _, item := range items {
		total += item.Subtotal
	}

	return database.DB.Model(&model.SalesOrder{}).Where("id = ?", orderID).Update("total_amount", total).Error
}

func generateOrderNumber(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s-%d-%04d", prefix, now.Year(), now.UnixNano()%10000)
}

func applySalesOrderFilter(query *gorm.DB, filter SalesOrderFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Joins("LEFT JOIN customers ON customers.id = sales_orders.customer_id").
			Where("sales_orders.order_number ILIKE ? OR customers.name ILIKE ?", search, search)
	}

	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}
