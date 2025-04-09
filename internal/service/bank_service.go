package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BankService struct {
	db *gorm.DB
}

func NewBankService(db *gorm.DB) *BankService {
	return &BankService{db: db}
}

func (s *BankService) GetByID(bankID uint) (*model.MstBank, error) {
	var bank model.MstBank
	if err := s.db.Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, account, description")
	}).Preload("Currency", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, currency, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&bank, bankID).Error; err != nil {
		return nil, err
	}

	return &bank, nil
}

func (s *BankService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstBank{})

	if search != "" {
		query = query.Where("code ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BankService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstBank, error) {
	var banks []model.MstBank

	query := s.db.Model(&model.MstBank{}).Preload("Account", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, account, description")
	}).Preload("Currency", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, currency, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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
		query = query.Where("code ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&banks).Error; err != nil {
		return nil, err
	}

	return banks, nil
}

func (s *BankService) Create(bank *model.MstBank) error {
	return s.db.Create(bank).Error
}

func (s *BankService) Update(bank *model.MstBank) error {
	return s.db.Save(bank).Error
}

func (s *BankService) Delete(bank *model.MstBank) error {
	return s.db.Delete(bank).Error
}
