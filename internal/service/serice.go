package service

import "github.com/afiffazun/inventory-api/internal/model"

func GetHome() model.Response {
	return model.Response{
		Application: "Inventory API",
		Author:      "Afif",
		Status:      "Running",
	}
}

func GetVersion() model.Response {
	return model.Response{
		Version: "v1.0.0",
	}
}