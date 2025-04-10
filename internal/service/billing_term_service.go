package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BillingTermService struct {
	db *gorm.DB
}

func NewBillingTermService(db *gorm.DB) *BillingTermService {
	return &BillingTermService{db: db}
}

func (s *BillingTermService) GetByID(billingTermID uint) (*model.MstBillingTerm, error) {
	var billingTerm model.MstBillingTerm
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&billingTerm, billingTermID).Error; err != nil {
		return nil, err
	}

	return &billingTerm, nil
}

func (s *BillingTermService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstBillingTerm{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BillingTermService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstBillingTerm, error) {
	var billingTerms []model.MstBillingTerm

	query := s.db.Model(&model.MstBillingTerm{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&billingTerms).Error; err != nil {
		return nil, err
	}

	return billingTerms, nil
}

func (s *BillingTermService) Create(billingTerm *model.MstBillingTerm) error {
	return s.db.Create(billingTerm).Error
}

func (s *BillingTermService) Update(billingTerm *model.MstBillingTerm) error {
	return s.db.Save(billingTerm).Error
}

func (s *BillingTermService) Delete(billingTerm *model.MstBillingTerm) error {
	return s.db.Delete(billingTerm).Error
}
