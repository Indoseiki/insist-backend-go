package service

import (
	"fmt"
	"insist-backend-golang/internal/dto"

	"gorm.io/gorm"
)

type ItemInforService struct {
	db *gorm.DB
}

func NewItemInforService(db *gorm.DB) *ItemInforService {
	return &ItemInforService{db: db}
}

func (s *ItemInforService) GetTotal(search string) (int64, error) {
	var count int64

	var raw string
	var args []interface{}

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
		raw = `
			SELECT COUNT(*) FROM (
				SELECT item AS code, Uf_description AS description, u_m AS uom
				FROM item_mst
				WHERE LOWER(item) LIKE LOWER(?) OR LOWER(Uf_description) LIKE LOWER(?)

				UNION

				SELECT item AS code, Uf_description2 AS description, u_m AS uom
				FROM non_inventory_item_mst
				WHERE LOWER(item) LIKE LOWER(?) OR LOWER(Uf_description2) LIKE LOWER(?)
			) AS item_infor
		`
		args = []interface{}{search, search, search, search}
	} else {
		raw = `
			SELECT COUNT(*) FROM (
				SELECT item AS code, Uf_description AS description, u_m AS uom
				FROM item_mst

				UNION

				SELECT item AS code, Uf_description2 AS description, u_m AS uom
				FROM non_inventory_item_mst
			) AS item_infor
		`
	}

	if err := s.db.Raw(raw, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ItemInforService) GetAll(offset, limit int, search, sortBy string, sortAsc bool) ([]dto.ItemInforDTO, error) {
	var items []dto.ItemInforDTO

	sortColumn := "code"
	if sortBy != "" {
		sortColumn = sortBy
	}
	sortDirection := "ASC"
	if !sortAsc {
		sortDirection = "DESC"
	}

	var raw string
	var args []interface{}

	if search != "" {
		search = fmt.Sprintf("%%%s%%", search)
		raw = fmt.Sprintf(`
			SELECT * FROM (
				SELECT item AS code, Uf_description AS description, u_m AS uom
				FROM item_mst
				WHERE LOWER(item) LIKE LOWER(?) OR LOWER(Uf_description) LIKE LOWER(?)

				UNION

				SELECT item AS code, Uf_description2 AS description, u_m AS uom
				FROM non_inventory_item_mst
				WHERE LOWER(item) LIKE LOWER(?) OR LOWER(Uf_description2) LIKE LOWER(?)
			) AS item_infor
			ORDER BY %s %s
			OFFSET ? ROWS FETCH NEXT ? ROWS ONLY
		`, sortColumn, sortDirection)
		args = []interface{}{search, search, search, search, offset, limit}
	} else {
		raw = fmt.Sprintf(`
			SELECT * FROM (
				SELECT item AS code, Uf_description AS description, u_m AS uom
				FROM item_mst

				UNION

				SELECT item AS code, Uf_description2 AS description, u_m AS uom
				FROM non_inventory_item_mst
			) AS item_infor
			ORDER BY %s %s
			OFFSET ? ROWS FETCH NEXT ? ROWS ONLY
		`, sortColumn, sortDirection)
		args = []interface{}{offset, limit}
	}

	if err := s.db.Raw(raw, args...).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
