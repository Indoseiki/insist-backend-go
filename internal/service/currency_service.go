package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CurrencyService struct {
	db *gorm.DB
}

func NewCurrencyService(db *gorm.DB) *CurrencyService {
	return &CurrencyService{db: db}
}

func (s *CurrencyService) GetByID(currencyID uint) (*model.MstCurrency, error) {
	var currency model.MstCurrency
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&currency, currencyID).Error; err != nil {
		return nil, err
	}

	return &currency, nil
}

func (s *CurrencyService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstCurrency{})

	if search != "" {
		query = query.Where("currency ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *CurrencyService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstCurrency, error) {
	var currencys []model.MstCurrency

	query := s.db.Model(&model.MstCurrency{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("currency ASC")
	}

	if search != "" {
		query = query.Where("currency ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&currencys).Error; err != nil {
		return nil, err
	}

	return currencys, nil
}

func (s *CurrencyService) Create(currency *model.MstCurrency) error {
	return s.db.Create(currency).Error
}

func (s *CurrencyService) Update(currency *model.MstCurrency) error {
	return s.db.Save(currency).Error
}

func (s *CurrencyService) Delete(currency *model.MstCurrency) error {
	return s.db.Delete(currency).Error
}

func (s *CurrencyService) GetByCurrencyCode(code string) (*model.MstCurrency, error) {
	var currency model.MstCurrency
	if err := s.db.Where("currency = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}
