package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MachineService struct {
	db *gorm.DB
}

func NewMachineService(db *gorm.DB) *MachineService {
	return &MachineService{db: db}
}

func (s *MachineService) BeginTx() *gorm.DB {
	return s.db.Begin()
}

func (s *MachineService) GetByID(machineID uint) (*model.ViewMstMachine, error) {
	var machine model.ViewMstMachine
	if err := s.db.First(&machine, machineID).Error; err != nil {
		return nil, err
	}
	return &machine, nil
}

func (s *MachineService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.ViewMstMachine{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *MachineService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.ViewMstMachine, error) {
	var machines []model.ViewMstMachine

	query := s.db.Model(&model.ViewMstMachine{}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("machine_updated_at ASC")
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&machines).Error; err != nil {
		return nil, err
	}

	return machines, nil
}

func (s *MachineService) Create(machine *model.MstMachine) error {
	return s.db.Create(machine).Error
}

func (s *MachineService) Update(machine *model.MstMachine) error {
	return s.db.Save(machine).Error
}

func (s *MachineService) Delete(machineID uint) error {
	return s.db.Where("id = ?", machineID).Delete(&model.MstMachine{}).Error
}

func (s *MachineService) GetTotalDetail(machineID uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.ViewMstMachineDetail{}).Where("id_machine = ?", machineID)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *MachineService) GetAllDetail(offset, limit int, sortBy string, sortDirection bool, machineID uint) ([]model.ViewMstMachineDetail, error) {
	var machineDetails []model.ViewMstMachineDetail

	query := s.db.Model(&model.ViewMstMachineDetail{}).
		Where("id_machine = ?", machineID).
		Offset(offset).
		Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("detail_updated_at ASC")
	}

	if err := query.Find(&machineDetails).Error; err != nil {
		return nil, err
	}

	return machineDetails, nil
}

func (s *MachineService) GetDetailByID(machineDetailID uint) (*model.MstMachineDetail, error) {
	var machineDetail model.MstMachineDetail
	if err := s.db.First(&machineDetail, machineDetailID).Error; err != nil {
		return nil, err
	}
	return &machineDetail, nil
}

func (s *MachineService) GetDetailsByMachineID(machineID uint) ([]model.MstMachineDetail, error) {
	var machineDetails []model.MstMachineDetail
	if err := s.db.Where("id_machine = ?", machineID).Find(&machineDetails).Error; err != nil {
		return nil, err
	}
	return machineDetails, nil
}

func (s *MachineService) CreateDetail(machineDetail *model.MstMachineDetail) error {
	return s.db.Create(machineDetail).Error
}

func (s *MachineService) UpdateDetail(machineDetail *model.MstMachineDetail) error {
	return s.db.Save(machineDetail).Error
}

func (s *MachineService) DeleteDetail(machineDetailID uint) error {
	return s.db.Where("id = ?", machineDetailID).Delete(&model.MstMachineDetail{}).Error
}

func (s *MachineService) GetTotalStatus(machineID uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstMachineStatus{}).Where("id_machine = ?", machineID)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *MachineService) GetAllStatus(offset, limit int, sortBy string, sortDirection bool, machineID uint) ([]model.MstMachineStatus, error) {
	var machineStatus []model.MstMachineStatus

	query := s.db.Model(&model.MstMachineStatus{}).Preload("Reason", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, key, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Where("id_machine = ?", machineID).
		Offset(offset).
		Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at DESC")
	}

	if err := query.Find(&machineStatus).Error; err != nil {
		return nil, err
	}

	return machineStatus, nil
}

func (s *MachineService) GetStatusByMachineID(machineID uint) ([]model.MstMachineStatus, error) {
	var machineDetails []model.MstMachineStatus
	if err := s.db.Where("id_machine = ?", machineID).Preload("Reason", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, key, code, description")
	}).Find(&machineDetails).Error; err != nil {
		return nil, err
	}
	return machineDetails, nil
}

func (s *MachineService) CreateStatus(machineStatus *model.MstMachineStatus) error {
	return s.db.Create(machineStatus).Error
}

func (s *MachineService) DeleteStatus(machineID uint) error {
	return s.db.Where("id_machine = ?", machineID).Delete(&model.MstMachineStatus{}).Error
}
