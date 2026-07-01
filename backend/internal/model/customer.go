package model

import "time"

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255)"`
	Phone     string    `json:"phone" gorm:"type:varchar(50)"`
	Address   string    `json:"address" gorm:"type:varchar(500)"`
	City      string    `json:"city" gorm:"type:varchar(100)"`
	Country   string    `json:"country" gorm:"type:varchar(100)"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Customer) Validate() []ValidationError {
	var errs []ValidationError

	if c.Code == "" {
		errs = append(errs, ValidationError{Field: "code", Message: "Code is required"})
	}

	if c.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "Name is required"})
	}

	return errs
}
