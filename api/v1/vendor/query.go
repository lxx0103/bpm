package vendor

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type vendorQuery struct {
	conn *sqlx.DB
}

func NewVendorQuery(connection *sqlx.DB) *vendorQuery {
	return &vendorQuery{
		conn: connection,
	}
}

func (r *vendorQuery) GetVendorByID(id int64) (*VendorResponse, error) {
	var vendor VendorResponse
	err := r.conn.Get(&vendor, "SELECT id, name, phone, address, longitude, latitude, cover, description, status FROM vendors WHERE id = ? AND status > 0 ", id)
	return &vendor, err
}

func (r *vendorQuery) GetVendorCount(filter VendorFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Brand; v != "" {
		where, args = append(where, "id in (SELECT vendor_id from vendor_brands WHERE brand_id in (SELECT id FROM brands WHERE name like ?))"), append(args, "%"+v+"%")
	}
	if v := filter.Material; v != "" {
		where, args = append(where, "id in (SELECT vendor_id from vendor_materials WHERE material_id in (SELECT id FROM materials WHERE name like ?))"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM vendors 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *vendorQuery) GetVendorList(filter VendorFilter) (*[]VendorResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Brand; v != "" {
		where, args = append(where, "id in (SELECT vendor_id from vendor_brands WHERE brand_id in (SELECT id FROM brands WHERE name like ?))"), append(args, "%"+v+"%")
	}
	if v := filter.Material; v != "" {
		where, args = append(where, "id in (SELECT vendor_id from vendor_materials WHERE material_id in (SELECT id FROM materials WHERE name like ?))"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var vendors []VendorResponse
	err := r.conn.Select(&vendors, `
		SELECT id, name, phone, address, longitude, latitude, cover, description, status
		FROM vendors 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &vendors, err
}

func (r *vendorQuery) GetVendorMaterial(vendorID int64) (*[]VendorMaterial, error) {
	var vendorMaterial []VendorMaterial
	err := r.conn.Select(&vendorMaterial, `
		SELECT vm.material_id, m.name as material_name 
		FROM vendor_materials vm LEFT JOIN materials m ON vm.material_id = m.id 
		WHERE vm.vendor_id = ?
		AND vm.status > 0
	`, vendorID)
	return &vendorMaterial, err
}

func (r *vendorQuery) GetVendorBrand(vendorID int64) (*[]VendorBrand, error) {
	var vendorBrand []VendorBrand
	err := r.conn.Select(&vendorBrand, `
		SELECT vm.brand_id, m.name as brand_name 
		FROM vendor_brands vm LEFT JOIN brands m ON vm.brand_id = m.id 
		WHERE vm.vendor_id = ?
		AND vm.status > 0
	`, vendorID)
	return &vendorBrand, err
}

func (r *vendorQuery) GetVendorPicture(vendorID int64) (*[]string, error) {
	var vendorPicture []string
	err := r.conn.Select(&vendorPicture, `
		SELECT name 
		FROM vendor_pictures 
		WHERE vendor_id = ?
		AND status > 0
	`, vendorID)
	return &vendorPicture, err
}
