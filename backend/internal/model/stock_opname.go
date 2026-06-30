package model

import (
	"time"
)

type StockOpname struct {
	ID          uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	WarehouseID uint              `json:"warehouse_id" gorm:"not null;index"`
	Warehouse   Warehouse         `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Status      string            `json:"status" gorm:"type:varchar(20);not null;default:DRAFT"`
	Notes       string            `json:"notes" gorm:"type:varchar(500)"`
	TotalDiff   int               `json:"total_diff" gorm:"default:0"`
	CreatedBy   string            `json:"created_by" gorm:"type:varchar(100)"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Items       []StockOpnameItem `json:"items" gorm:"foreignKey:StockOpnameID"`
}

type StockOpnameItem struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StockOpnameID  uint      `json:"stock_opname_id" gorm:"not null;index"`
	ItemID         uint      `json:"item_id" gorm:"not null;index"`
	Item           Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	SystemQuantity int       `json:"system_quantity" gorm:"not null"`
	ActualQuantity int       `json:"actual_quantity"`
	Difference     int       `json:"difference" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
}

const (
	OpnameStatusDRAFT      = "DRAFT"
	OpnameStatusINPROGRESS = "IN_PROGRESS"
	OpnameStatusCOMPLETED  = "COMPLETED"
	OpnameStatusCANCELLED  = "CANCELLED"
)

func (s *StockOpname) Validate() []ValidationError {
	var errs []ValidationError

	if s.WarehouseID == 0 {
		errs = append(errs, ValidationError{Field: "warehouse_id", Message: "WarehouseID is required"})
	}

	if s.Status == "" {
		s.Status = OpnameStatusDRAFT
	} else if s.Status != OpnameStatusDRAFT && s.Status != OpnameStatusINPROGRESS && s.Status != OpnameStatusCOMPLETED && s.Status != OpnameStatusCANCELLED {
		errs = append(errs, ValidationError{Field: "status", Message: "Status must be DRAFT, IN_PROGRESS, COMPLETED, or CANCELLED"})
	}

	return errs
}

func (s *StockOpnameItem) Validate() []ValidationError {
	var errs []ValidationError

	if s.ItemID == 0 {
		errs = append(errs, ValidationError{Field: "item_id", Message: "ItemID is required"})
	}

	if s.SystemQuantity < 0 {
		errs = append(errs, ValidationError{Field: "system_quantity", Message: "SystemQuantity must be non-negative"})
	}

	if s.ActualQuantity < 0 {
		errs = append(errs, ValidationError{Field: "actual_quantity", Message: "ActualQuantity must be non-negative"})
	}

	return errs
}

func (s *StockOpnameItem) CalculateDifference() {
	s.Difference = s.ActualQuantity - s.SystemQuantity
}

type CreateOpnameRequest struct {
	WarehouseID uint               `json:"warehouse_id"`
	Notes       string             `json:"notes"`
	CreatedBy   string             `json:"created_by"`
	Items       []CreateOpnameItem `json:"items"`
}

type CreateOpnameItem struct {
	ItemID         uint `json:"item_id"`
	ActualQuantity int  `json:"actual_quantity"`
}
