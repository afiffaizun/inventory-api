package model

import "time"

type PurchaseOrder struct {
	ID          uint                `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNumber string              `json:"order_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	SupplierID  uint                `json:"supplier_id" gorm:"not null;index"`
	Supplier    Supplier            `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	WarehouseID uint                `json:"warehouse_id" gorm:"not null;index"`
	Warehouse   Warehouse           `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	OrderDate   time.Time           `json:"order_date" gorm:"not null"`
	Status      string              `json:"status" gorm:"type:varchar(20);not null;default:DRAFT"`
	Notes       string              `json:"notes" gorm:"type:varchar(500)"`
	TotalAmount float64             `json:"total_amount" gorm:"type:decimal(15,2);default:0"`
	CreatedBy   string              `json:"created_by" gorm:"type:varchar(100)"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Items       []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
}

type PurchaseOrderItem struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PurchaseOrderID uint      `json:"purchase_order_id" gorm:"not null;index"`
	ItemID          uint      `json:"item_id" gorm:"not null;index"`
	Item            Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Quantity        int       `json:"quantity" gorm:"not null"`
	UnitCost        float64   `json:"unit_cost" gorm:"type:decimal(15,2);not null"`
	Subtotal        float64   `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	CreatedAt       time.Time `json:"created_at"`
}

const (
	PurchaseOrderStatusDRAFT     = "DRAFT"
	PurchaseOrderStatusCONFIRMED = "CONFIRMED"
	PurchaseOrderStatusRECEIVED  = "RECEIVED"
	PurchaseOrderStatusCOMPLETED = "COMPLETED"
	PurchaseOrderStatusCANCELLED = "CANCELLED"
)

func (p *PurchaseOrder) Validate() []ValidationError {
	var errs []ValidationError

	if p.SupplierID == 0 {
		errs = append(errs, ValidationError{Field: "supplier_id", Message: "SupplierID is required"})
	}

	if p.WarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "warehouse_id", Message: "WarehouseID is required"})
	}

	if p.OrderDate.IsZero() {
		errs = append(errs, ValidationError{Field: "order_date", Message: "OrderDate is required"})
	}

	if p.Status == "" {
		p.Status = PurchaseOrderStatusDRAFT
	}

	return errs
}

func (p *PurchaseOrderItem) Validate() []ValidationError {
	var errs []ValidationError

	if p.ItemID == 0 {
		errs = append(errs, ValidationError{Field: "item_id", Message: "ItemID is required"})
	}

	if p.Quantity <= 0 {
		errs = append(errs, ValidationError{Field: "quantity", Message: "Quantity must be positive"})
	}

	if p.UnitCost < 0 {
		errs = append(errs, ValidationError{Field: "unit_cost", Message: "UnitCost must be non-negative"})
	}

	return errs
}

func (p *PurchaseOrderItem) CalculateSubtotal() {
	p.Subtotal = float64(p.Quantity) * p.UnitCost
}

type CreatePurchaseOrderRequest struct {
	SupplierID  uint                      `json:"supplier_id"`
	WarehouseID uint                      `json:"warehouse_id"`
	OrderDate   string                    `json:"order_date"`
	Notes       string                    `json:"notes"`
	CreatedBy   string                    `json:"created_by"`
	Items       []CreatePurchaseOrderItem `json:"items"`
}

type CreatePurchaseOrderItem struct {
	ItemID   uint    `json:"item_id"`
	Quantity int     `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
}
