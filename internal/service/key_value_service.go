package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type KeyValueService struct {
	db *gorm.DB
}

func NewKeyValueService(db *gorm.DB) *KeyValueService {
	return &KeyValueService{db: db}
}

func (s *KeyValueService) GetByID(keyValueID uint) (*model.MstKeyValue, error) {
	var keyValue model.MstKeyValue
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&keyValue, keyValueID).Error; err != nil {
		return nil, err
	}

	return &keyValue, nil
}

func (s *KeyValueService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstKeyValue{})

	if search != "" {
		query = query.Where("key ILIKE ? OR value ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *KeyValueService) GetAll(offset, limit int, search string) ([]model.MstKeyValue, error) {
	var keyValues []model.MstKeyValue

	query := s.db.Model(&model.MstKeyValue{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("key ILIKE ? OR value ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&keyValues).Error; err != nil {
		return nil, err
	}

	return keyValues, nil
}

func (s *KeyValueService) Create(keyValue *model.MstKeyValue) error {
	return s.db.Create(keyValue).Error
}

func (s *KeyValueService) Update(keyValue *model.MstKeyValue) error {
	return s.db.Save(keyValue).Error
}

func (s *KeyValueService) Delete(keyValue *model.MstKeyValue) error {
	return s.db.Delete(keyValue).Error
}
