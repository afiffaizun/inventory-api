package model

import (
	"time"
)

type Category struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string     `json:"code" gorm:"type:varchar(20);uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Description string     `json:"description" gorm:"type:varchar(500)"`
	ParentID    *uint      `json:"parent_id" gorm:"index"`
	Parent      *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (c *Category) Validate() []ValidationError {
	var errs []ValidationError

	if c.Code == "" {
		errs = append(errs, ValidationError{Field: "code", Message: "Code is required"})
	} else if len(c.Code) > 20 {
		errs = append(errs, ValidationError{Field: "code", Message: "Code max 20 chars"})
	}

	if c.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Name is required"})
	} else if len(c.Name) > 100 {
		errs = append(errs, ValidationError{Field: "name", Message: "Name max 100 chars"})
	}

	if len(c.Description) > 500 {
		errs = append(errs, ValidationError{Field: "description", Message: "Description max 500 chars"})
	}

	return errs
}

type CategoryTree struct {
	ID          uint           `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active"`
	Children    []CategoryTree `json:"children"`
}
