package model

import (
	"time"
)

type Warehouse struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(20);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Address   string    `json:"address" gorm:"type:varchar(500)"`
	City      string    `json:"city" gorm:"type:varchar(100)"`
	Country   string    `json:"country" gorm:"type:varchar(100)"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w *Warehouse) Validate() []ValidationError {
	var errs []ValidationError

	if w.Code == "" {
		errs = append(errs, ValidationError{Field: "code", Message: "Code is required"})
	} else if len(w.Code) > 20 {
		errs = append(errs, ValidationError{Field: "code", Message: "Code max 20 chars"})
	}

	if w.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Name is required"})
	} else if len(w.Name) > 100 {
		errs = append(errs, ValidationError{Field: "name", Message: "Name max 100 chars"})
	}

	if len(w.Address) > 500 {
		errs = append(errs, ValidationError{Field: "address", Message: "Address max 500 chars"})
	}

	if len(w.City) > 100 {
		errs = append(errs, ValidationError{Field: "city", Message: "City max 100 chars"})
	}

	if len(w.Country) > 100 {
		errs = append(errs, ValidationError{Field: "country", Message: "Country max 100 chars"})
	}

	return errs
}
