package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaxCodeService struct {
	db *gorm.DB
}

func NewTaxCodeService(db *gorm.DB) *TaxCodeService {
	return &TaxCodeService{db: db}
}

func (s *TaxCodeService) GetByID(taxCodeID uint) (*model.MstTaxCode, error) {
	var taxCode model.MstTaxCode
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&taxCode, taxCodeID).Error; err != nil {
		return nil, err
	}

	return &taxCode, nil
}

func (s *TaxCodeService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstTaxCode{})

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *TaxCodeService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstTaxCode, error) {
	var taxCodes []model.MstTaxCode

	query := s.db.Model(&model.MstTaxCode{}).
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("AccountAR").
		Preload("AccountARProc").
		Preload("AccountAP").
		Offset(offset).
		Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&taxCodes).Error; err != nil {
		return nil, err
	}

	return taxCodes, nil
}

func (s *TaxCodeService) Create(taxCode *model.MstTaxCode) error {
	return s.db.Create(taxCode).Error
}

func (s *TaxCodeService) Update(taxCode *model.MstTaxCode) error {
	return s.db.Save(taxCode).Error
}

func (s *TaxCodeService) Delete(taxCode *model.MstTaxCode) error {
	return s.db.Delete(taxCode).Error
}
