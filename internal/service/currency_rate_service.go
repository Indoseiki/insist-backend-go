package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CurrencyRateService struct {
	db *gorm.DB
}

func NewCurrencyRateService(db *gorm.DB) *CurrencyRateService {
	return &CurrencyRateService{db: db}
}

func (s *CurrencyRateService) GetByID(currencyRateID uint) (*model.MstCurrencyRate, error) {
	var currencyRate model.MstCurrencyRate
	if err := s.db.Preload("FromCurrency").Preload("ToCurrency").Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&currencyRate, currencyRateID).Error; err != nil {
		return nil, err
	}

	return &currencyRate, nil
}

func (s *CurrencyRateService) GetTotal(idCurrency uint, search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstCurrencyRate{})

	if idCurrency != 0 {
		query = query.Where("id_from_currency = ?", idCurrency)
	}

	if search != "" {
		query = query.Joins("LEFT JOIN mst_currencies AS fc ON fc.id = mst_currency_rates.id_from_currency").
			Joins("LEFT JOIN mst_currencies AS tc ON tc.id = mst_currency_rates.id_to_currency").
			Where("fc.currency ILIKE ? OR tc.currency ILIKE ?", "%"+search+"%", "%"+search+"%")

	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *CurrencyRateService) GetAll(idCurrency uint, offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstCurrencyRate, error) {
	var currencyRates []model.MstCurrencyRate

	query := s.db.Model(&model.MstCurrencyRate{}).Preload("FromCurrency").Preload("ToCurrency").Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if idCurrency != 0 {
		query = query.Where("id_from_currency = ?", idCurrency)
	}

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Joins("LEFT JOIN mst_currencies AS fc ON fc.id = mst_currency_rates.id_from_currency").
			Joins("LEFT JOIN mst_currencies AS tc ON tc.id = mst_currency_rates.id_to_currency").
			Where("fc.currency ILIKE ? OR tc.currency ILIKE ?", "%"+search+"%", "%"+search+"%")

	}

	if err := query.Find(&currencyRates).Error; err != nil {
		return nil, err
	}

	return currencyRates, nil
}

func (s *CurrencyRateService) Create(currencyRate *model.MstCurrencyRate) error {
	return s.db.Create(currencyRate).Error
}

func (s *CurrencyRateService) Update(currencyRate *model.MstCurrencyRate) error {
	return s.db.Save(currencyRate).Error
}

func (s *CurrencyRateService) Delete(currencyRate *model.MstCurrencyRate) error {
	return s.db.Delete(currencyRate).Error
}
