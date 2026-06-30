package model

import "regexp"

type Item struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code     string `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Stock    int    `json:"stock" gorm:"default:0"`
	Location string `json:"location" gorm:"type:varchar(255)"`
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

	if i.Stock < 0 {
		errs = append(errs, ValidationError{Field: "stock", Message: "Stock must be non-negative"})
	}

	if len(i.Location) > 255 {
		errs = append(errs, ValidationError{Field: "location", Message: "Location max 255 chars"})
	}

	return errs
}
