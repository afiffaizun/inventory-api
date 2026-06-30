package model

type Item struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code     string `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Stock    int    `json:"stock" gorm:"default:0"`
	Location string `json:"location" gorm:"type:varchar(255)"`
}
