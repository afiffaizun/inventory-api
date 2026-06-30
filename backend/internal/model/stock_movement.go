package model

import (
	"time"
)

type StockMovement struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ItemID        uint      `json:"item_id" gorm:"not null;index"`
	Item          Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	WarehouseID   uint      `json:"warehouse_id" gorm:"not null;index"`
	Warehouse     Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Type          string    `json:"type" gorm:"type:varchar(20);not null"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	ReferenceType string    `json:"reference_type" gorm:"type:varchar(50)"`
	ReferenceID   uint      `json:"reference_id"`
	Notes         string    `json:"notes" gorm:"type:varchar(500)"`
	CreatedBy     string    `json:"created_by" gorm:"type:varchar(100)"`
	CreatedAt     time.Time `json:"created_at"`
}

const (
	MovementTypeIN         = "IN"
	MovementTypeOUT        = "OUT"
	MovementTypeTRANSFER   = "TRANSFER"
	MovementTypeADJUSTMENT = "ADJUSTMENT"
)

const (
	ReferenceTypePURCHASE = "PURCHASE"
	ReferenceTypeSALE     = "SALE"
	ReferenceTypeTRANSFER = "TRANSFER"
	ReferenceTypeOPNAME   = "OPNAME"
	ReferenceTypeADJUST   = "ADJUST"
)

func (s *StockMovement) Validate() []ValidationError {
	var errs []ValidationError

	if s.ItemID == 0 {
		errs = append(errs, ValidationError{Field: "item_id", Message: "ItemID is required"})
	}

	if s.WarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "warehouse_id", Message: "WarehouseID is required"})
	}

	if s.Type == "" {
		errs = append(errs, ValidationError{Field: "type", Message: "Type is required"})
	} else if s.Type != MovementTypeIN && s.Type != MovementTypeOUT && s.Type != MovementTypeTRANSFER && s.Type != MovementTypeADJUSTMENT {
		errs = append(errs, ValidationError{Field: "type", Message: "Type must be IN, OUT, TRANSFER, or ADJUSTMENT"})
	}

	if s.Quantity == 0 {
		errs = append(errs, ValidationError{Field: "quantity", Message: "Quantity must be non-zero"})
	}

	return errs
}

type TransferRequest struct {
	ItemID          uint   `json:"item_id"`
	FromWarehouseID uint   `json:"from_warehouse_id"`
	ToWarehouseID   uint   `json:"to_warehouse_id"`
	Quantity        int    `json:"quantity"`
	Notes           string `json:"notes"`
	CreatedBy       string `json:"created_by"`
}

func (t *TransferRequest) Validate() []ValidationError {
	var errs []ValidationError

	if t.ItemID == 0 {
		errs = append(errs, ValidationError{Field: "item_id", Message: "ItemID is required"})
	}

	if t.FromWarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "from_warehouse_id", Message: "FromWarehouseID is required"})
	}

	if t.ToWarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "to_warehouse_id", Message: "ToWarehouseID is required"})
	}

	if t.FromWarehouseID == t.ToWarehouseID {
		errs = append(errs, ValidationError{Field: "to_warehouse_id", Message: "Source and destination warehouses must be different"})
	}

	if t.Quantity <= 0 {
		errs = append(errs, ValidationError{Field: "quantity", Message: "Quantity must be positive"})
	}

	return errs
}
