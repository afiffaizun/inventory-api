package database

import (
	"fmt"
	"log"

	"github.com/afiffazun/inventory-api/internal/config"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error

	DB, err = gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&model.Warehouse{},
		&model.Category{},
		&model.Item{},
		&model.StockMovement{},
		&model.StockOpname{},
		&model.StockOpnameItem{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully")
}
