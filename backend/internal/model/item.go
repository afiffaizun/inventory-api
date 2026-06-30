package model

import (
	"regexp"
	"time"
)

type Item struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string     `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:varchar(1000)"`
	CategoryID  *uint      `json:"category_id" gorm:"index"`
	Category    *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	WarehouseID *uint      `json:"warehouse_id" gorm:"index"`
	Warehouse   *Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	SKU         string     `json:"sku" gorm:"type:varchar(100);uniqueIndex"`
	Unit        string     `json:"unit" gorm:"type:varchar(20)"`
	MinStock    int        `json:"min_stock" gorm:"default:0"`
	MaxStock    int        `json:"max_stock" gorm:"default:0"`
	Stock       int        `json:"stock" gorm:"default:0"`
	CostPrice   float64    `json:"cost_price" gorm:"type:decimal(15,2);default:0"`
	SellPrice   float64    `json:"sell_price" gorm:"type:decimal(15,2);default:0"`
	Location    string     `json:"location" gorm:"type:varchar(255)"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

var codeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func (i *Item) Validate() []ValidationError {
	var errs []ValidationError

	if i.Code == "" {
		errs = append(errs, ValidationError{Field: "code", Message: "Code is required"})
	} else if len(i.Code) > 50 {
		errs = append(errs, ValidationError{Field: "code", Message: "Code max 50 chars"})
	} else if !codeRegex.MatchString(i.Code) {
		errs = append(errs, ValidationError{Field: "code", Message: "Code format invalid (alphanumeric, dash, underscore only)"})
	}

	if i.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Name is required"})
	} else if len(i.Name) > 255 {
		errs = append(errs, ValidationError{Field: "name", Message: "Name max 255 chars"})
	}

	if len(i.Description) > 1000 {
		errs = append(errs, ValidationError{Field: "description", Message: "Description max 1000 chars"})
	}

	if len(i.SKU) > 100 {
		errs = append(errs, ValidationError{Field: "sku", Message: "SKU max 100 chars"})
	}

	if len(i.Unit) > 20 {
		errs = append(errs, ValidationError{Field: "unit", Message: "Unit max 20 chars"})
	}

	if i.Stock < 0 {
		errs = append(errs, ValidationError{Field: "stock", Message: "Stock must be non-negative"})
	}

	if i.MinStock < 0 {
		errs = append(errs, ValidationError{Field: "min_stock", Message: "MinStock must be non-negative"})
	}

	if i.MaxStock < 0 {
		errs = append(errs, ValidationError{Field: "max_stock", Message: "MaxStock must be non-negative"})
	}

	if i.MaxStock > 0 && i.MinStock > i.MaxStock {
		errs = append(errs, ValidationError{Field: "min_stock", Message: "MinStock cannot be greater than MaxStock"})
	}

	if i.CostPrice < 0 {
		errs = append(errs, ValidationError{Field: "cost_price", Message: "CostPrice must be non-negative"})
	}

	if i.SellPrice < 0 {
		errs = append(errs, ValidationError{Field: "sell_price", Message: "SellPrice must be non-negative"})
	}

	if len(i.Location) > 255 {
		errs = append(errs, ValidationError{Field: "location", Message: "Location max 255 chars"})
	}

	return errs
}
