package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemSubCategoryService struct {
	db *gorm.DB
}

func NewItemSubCategoryService(db *gorm.DB) *ItemSubCategoryService {
	return &ItemSubCategoryService{db: db}
}

func (s *ItemSubCategoryService) GetByID(itemSubCategoryID uint) (*model.MstItemSubCategory, error) {
	var itemSubCategory model.MstItemSubCategory
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemSubCategory, itemSubCategoryID).Error; err != nil {
		return nil, err
	}

	return &itemSubCategory, nil
}

func (s *ItemSubCategoryService) GetTotal(search string, itemCategoryCode string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemSubCategory{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_sub_categories.id_item_category")

	if itemCategoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", itemCategoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_sub_categories.code ILIKE ? OR mst_item_sub_categories.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemSubCategoryService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, itemCategoryCode string) ([]model.MstItemSubCategory, error) {
	var itemSubCategorys []model.MstItemSubCategory

	query := s.db.Model(&model.MstItemSubCategory{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_sub_categories.id_item_category").
		Preload("ItemCategory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).
		Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("mst_item_sub_categories.code ASC")
	}

	if itemCategoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", itemCategoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_sub_categories.code ILIKE ? OR mst_item_sub_categories.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemSubCategorys).Error; err != nil {
		return nil, err
	}

	return itemSubCategorys, nil
}

func (s *ItemSubCategoryService) Create(itemSubCategory *model.MstItemSubCategory) error {
	return s.db.Create(itemSubCategory).Error
}

func (s *ItemSubCategoryService) Update(itemSubCategory *model.MstItemSubCategory) error {
	return s.db.Save(itemSubCategory).Error
}

func (s *ItemSubCategoryService) Delete(itemSubCategory *model.MstItemSubCategory) error {
	return s.db.Delete(itemSubCategory).Error
}
