package repository

import (
	"github.com/afiffazun/inventory-api/internal/database"
	"github.com/afiffazun/inventory-api/internal/model"
	"gorm.io/gorm"
)

type CategoryFilter struct {
	Search   string
	ParentID *uint
}

func GetCategories(filter CategoryFilter, page, limit int) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := applyCategoryFilter(database.DB.Model(&model.Category{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Parent").Offset(offset).Limit(limit).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	result := database.DB.Where("is_active = ?", true).Order("id ASC").Find(&categories)
	return categories, result.Error
}

func GetCategoryTree() ([]model.CategoryTree, error) {
	var categories []model.Category
	if err := database.DB.Where("is_active = ? AND parent_id IS NULL", true).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	var trees []model.CategoryTree
	for _, cat := range categories {
		tree := model.CategoryTree{
			ID:          cat.ID,
			Code:        cat.Code,
			Name:        cat.Name,
			Description: cat.Description,
			IsActive:    cat.IsActive,
		}
		children, err := getChildCategories(cat.ID)
		if err != nil {
			return nil, err
		}
		tree.Children = children
		trees = append(trees, tree)
	}

	return trees, nil
}

func getChildCategories(parentID uint) ([]model.CategoryTree, error) {
	var categories []model.Category
	if err := database.DB.Where("is_active = ? AND parent_id = ?", true, parentID).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	var trees []model.CategoryTree
	for _, cat := range categories {
		tree := model.CategoryTree{
			ID:          cat.ID,
			Code:        cat.Code,
			Name:        cat.Name,
			Description: cat.Description,
			IsActive:    cat.IsActive,
		}
		children, err := getChildCategories(cat.ID)
		if err != nil {
			return nil, err
		}
		tree.Children = children
		trees = append(trees, tree)
	}

	return trees, nil
}

func GetCategoryByID(id uint) (*model.Category, error) {
	var category model.Category
	result := database.DB.Preload("Parent").First(&category, id)
	return &category, result.Error
}

func CreateCategory(category *model.Category) error {
	result := database.DB.Create(category)
	return result.Error
}

func UpdateCategory(id uint, category *model.Category) error {
	result := database.DB.Model(&model.Category{}).Where("id = ?", id).Updates(category)
	return result.Error
}

func DeleteCategory(id uint) error {
	result := database.DB.Delete(&model.Category{}, id)
	return result.Error
}

func applyCategoryFilter(query *gorm.DB, filter CategoryFilter) *gorm.DB {
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}

	return query
}
