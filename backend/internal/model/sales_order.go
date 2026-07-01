package model

import "time"

type SalesOrder struct {
	ID          uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNumber string           `json:"order_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	CustomerID  uint             `json:"customer_id" gorm:"not null;index"`
	Customer    Customer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	WarehouseID uint             `json:"warehouse_id" gorm:"not null;index"`
	Warehouse   Warehouse        `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	OrderDate   time.Time        `json:"order_date" gorm:"not null"`
	Status      string           `json:"status" gorm:"type:varchar(20);not null;default:DRAFT"`
	Notes       string           `json:"notes" gorm:"type:varchar(500)"`
	TotalAmount float64          `json:"total_amount" gorm:"type:decimal(15,2);default:0"`
	CreatedBy   string           `json:"created_by" gorm:"type:varchar(100)"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Items       []SalesOrderItem `json:"items" gorm:"foreignKey:SalesOrderID"`
}

type SalesOrderItem struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SalesOrderID uint      `json:"sales_order_id" gorm:"not null;index"`
	ItemID       uint      `json:"item_id" gorm:"not null;index"`
	Item         Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Quantity     int       `json:"quantity" gorm:"not null"`
	UnitPrice    float64   `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Subtotal     float64   `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	SalesOrderStatusDRAFT     = "DRAFT"
	SalesOrderStatusCONFIRMED = "CONFIRMED"
	SalesOrderStatusSHIPPED   = "SHIPPED"
	SalesOrderStatusCOMPLETED = "COMPLETED"
	SalesOrderStatusCANCELLED = "CANCELLED"
)

func (s *SalesOrder) Validate() []ValidationError {
	var errs []ValidationError

	if s.CustomerID == 0 {
		errs = append(errs, ValidationError{Field: "customer_id", Message: "CustomerID is required"})
	}

	if s.WarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "warehouse_id", Message: "WarehouseID is required"})
	}

	if s.OrderDate.IsZero() {
		errs = append(errs, ValidationError{Field: "order_date", Message: "OrderDate is required"})
	}

	if s.Status == "" {
		s.Status = SalesOrderStatusDRAFT
	}

	return errs
}

func (s *SalesOrderItem) Validate() []ValidationError {
	var errs []ValidationError

	if s.ItemID == 0 {
		errs = append(errs, ValidationError{Field: "item_id", Message: "ItemID is required"})
	}

	if s.Quantity <= 0 {
		errs = append(errs, ValidationError{Field: "quantity", Message: "Quantity must be positive"})
	}

	if s.UnitPrice < 0 {
		errs = append(errs, ValidationError{Field: "unit_price", Message: "UnitPrice must be non-negative"})
	}

	return errs
}

func (s *SalesOrderItem) CalculateSubtotal() {
	s.Subtotal = float64(s.Quantity) * s.UnitPrice
}

type CreateSalesOrderRequest struct {
	CustomerID  uint                   `json:"customer_id"`
	WarehouseID uint                   `json:"warehouse_id"`
	OrderDate   string                 `json:"order_date"`
	Notes       string                 `json:"notes"`
	CreatedBy   string                 `json:"created_by"`
	Items       []CreateSalesOrderItem `json:"items"`
}

type CreateSalesOrderItem struct {
	ItemID    uint    `json:"item_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
