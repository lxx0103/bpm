package common

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type commonQuery struct {
	conn *sqlx.DB
}

func NewCommonQuery(connection *sqlx.DB) *commonQuery {
	return &commonQuery{
		conn: connection,
	}
}

func (r *commonQuery) GetBrandByID(id int64) (*BrandResponse, error) {
	var brand BrandResponse
	err := r.conn.Get(&brand, "SELECT id, name, status FROM brands WHERE id = ? AND status > 0 ", id)
	return &brand, err
}

func (r *commonQuery) GetBrandCount(filter BrandFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM brands 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *commonQuery) GetBrandList(filter BrandFilter) (*[]BrandResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var brands []BrandResponse
	err := r.conn.Select(&brands, `
		SELECT id, name, status
		FROM brands 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &brands, err
}

func (r *commonQuery) GetMaterialByID(id int64) (*MaterialResponse, error) {
	var material MaterialResponse
	err := r.conn.Get(&material, "SELECT id, name, status FROM materials WHERE id = ? AND status > 0 ", id)
	return &material, err
}

func (r *commonQuery) GetMaterialCount(filter MaterialFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM materials 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *commonQuery) GetMaterialList(filter MaterialFilter) (*[]MaterialResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var materials []MaterialResponse
	err := r.conn.Select(&materials, `
		SELECT id, name, status
		FROM materials 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &materials, err
}
