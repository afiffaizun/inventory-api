package repository

import (
	"fmt"

	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type PurchaseOrderFilter struct {
	Search     string
	SupplierID *uint
	Status     string
}

func GetPurchaseOrders(filter PurchaseOrderFilter, page, limit int) ([]model.PurchaseOrder, int64, error) {
	var orders []model.PurchaseOrder
	var total int64

	query := applyPurchaseOrderFilter(database.DB.Model(&model.PurchaseOrder{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Supplier").Preload("Warehouse").Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetPurchaseOrderByID(id uint) (*model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	result := database.DB.Preload("Supplier").Preload("Warehouse").Preload("Items").Preload("Items.Item").First(&order, id)
	return &order, result.Error
}

func CreatePurchaseOrder(order *model.PurchaseOrder) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		order.OrderNumber = generateOrderNumber("PO")
		items := order.Items
		order.Items = nil

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].PurchaseOrderID = order.ID
			items[i].CalculateSubtotal()
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		order.Items = items
		return recalculatePurchaseOrderTotal(order.ID)
	})
}

func UpdatePurchaseOrder(id uint, order *model.PurchaseOrder) error {
	result := database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeletePurchaseOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("purchase_order_id = ?", id).Delete(&model.PurchaseOrderItem{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.PurchaseOrder{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func ConfirmPurchaseOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.PurchaseOrder
		if err := tx.Preload("Items").First(&order, id).Error; err != nil {
			return err
		}

		if order.Status != model.PurchaseOrderStatusDRAFT {
			return fmt.Errorf("order cannot be confirmed in status: %s", order.Status)
		}

		for _, item := range order.Items {
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
				ReferenceType: model.ReferenceTypePURCHASE,
				ReferenceID:   order.ID,
				Notes:         fmt.Sprintf("Purchase Order %s", order.OrderNumber),
				CreatedBy:     order.CreatedBy,
			}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		order.Status = model.PurchaseOrderStatusCONFIRMED
		return tx.Save(&order).Error
	})
}

func CancelPurchaseOrder(id uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.PurchaseOrder
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}

		if order.Status == model.PurchaseOrderStatusCOMPLETED {
			return fmt.Errorf("cannot cancel completed order")
		}

		if order.Status == model.PurchaseOrderStatusCONFIRMED {
			if err := reversePurchaseOrderStock(tx, &order); err != nil {
				return err
			}
		}

		order.Status = model.PurchaseOrderStatusCANCELLED
		return tx.Save(&order).Error
	})
}

func reversePurchaseOrderStock(tx *gorm.DB, order *model.PurchaseOrder) error {
	var items []model.PurchaseOrderItem
	if err := tx.Where("purchase_order_id = ?", order.ID).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		var dbItem model.Item
		if err := tx.First(&dbItem, item.ItemID).Error; err != nil {
			return err
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
			ReferenceType: model.ReferenceTypePURCHASE,
			ReferenceID:   order.ID,
			Notes:         fmt.Sprintf("Cancelled Purchase Order %s", order.OrderNumber),
			CreatedBy:     order.CreatedBy,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
	}

	return nil
}

func recalculatePurchaseOrderTotal(orderID uint) error {
	var items []model.PurchaseOrderItem
	if err := database.DB.Where("purchase_order_id = ?", orderID).Find(&items).Error; err != nil {
		return err
	}

	var total float64
	for _, item := range items {
		total += item.Subtotal
	}

	return database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", orderID).Update("total_amount", total).Error
}

func applyPurchaseOrderFilter(query *gorm.DB, filter PurchaseOrderFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Joins("LEFT JOIN suppliers ON suppliers.id = purchase_orders.supplier_id").
			Where("purchase_orders.order_number ILIKE ? OR suppliers.name ILIKE ?", search, search)
	}

	if filter.SupplierID != nil {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}
