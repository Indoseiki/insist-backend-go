package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemGroupTypeService struct {
	db *gorm.DB
}

func NewItemGroupTypeService(db *gorm.DB) *ItemGroupTypeService {
	return &ItemGroupTypeService{db: db}
}

func (s *ItemGroupTypeService) GetByID(itemGroupTypeID uint) (*model.MstItemGroupType, error) {
	var itemGroupType model.MstItemGroupType
	if err := s.db.Preload("ItemGroup", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemGroupType, itemGroupTypeID).Error; err != nil {
		return nil, err
	}

	return &itemGroupType, nil
}

func (s *ItemGroupTypeService) GetTotal(search string, idItemGroup uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemGroupType{})

	if idItemGroup != 0 {
		query = query.Where("id_item_group = ?", idItemGroup)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemGroupTypeService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, idItemGroup uint) ([]model.MstItemGroupType, error) {
	var itemGroupTypes []model.MstItemGroupType

	query := s.db.Model(&model.MstItemGroupType{}).
		Preload("ItemGroup", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Offset(offset).Limit(limit)

	if idItemGroup != 0 {
		query = query.Where("id_item_group = ?", idItemGroup)
	}

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("code ASC")
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemGroupTypes).Error; err != nil {
		return nil, err
	}

	return itemGroupTypes, nil
}

func (s *ItemGroupTypeService) Create(itemGroupType *model.MstItemGroupType) error {
	return s.db.Create(itemGroupType).Error
}

func (s *ItemGroupTypeService) Update(itemGroupType *model.MstItemGroupType) error {
	return s.db.Save(itemGroupType).Error
}

func (s *ItemGroupTypeService) Delete(itemGroupType *model.MstItemGroupType) error {
	return s.db.Delete(itemGroupType).Error
}
