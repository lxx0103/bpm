package vendors

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type vendorsQuery struct {
	conn *sqlx.DB
}

func NewVendorsQuery(connection *sqlx.DB) *vendorsQuery {
	return &vendorsQuery{
		conn: connection,
	}
}

func (r *vendorsQuery) GetVendorsByID(id int64) (*VendorsResponse, error) {
	var vendors VendorsResponse
	err := r.conn.Get(&vendors, "SELECT id, contact, name, phone, address, longitude, latitude, cover, description, status FROM vendors WHERE id = ? AND status > 0 ", id)
	return &vendors, err
}

func (r *vendorsQuery) GetVendorsCount(filter VendorsFilter) (int, error) {
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

func (r *vendorsQuery) GetVendorsList(filter VendorsFilter) (*[]VendorsResponse, error) {
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
	var vendorss []VendorsResponse
	err := r.conn.Select(&vendorss, `
		SELECT id, contact, name, phone, address, longitude, latitude, cover, description, status
		FROM vendors
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &vendorss, err
}

func (r *vendorsQuery) GetVendorsMaterial(vendorsID int64) (*[]VendorsMaterial, error) {
	var vendorsMaterial []VendorsMaterial
	err := r.conn.Select(&vendorsMaterial, `
		SELECT vm.material_id, m.name as material_name 
		FROM vendor_materials vm LEFT JOIN materials m ON vm.material_id = m.id 
		WHERE vm.vendor_id = ?
		AND vm.status > 0
	`, vendorsID)
	return &vendorsMaterial, err
}

func (r *vendorsQuery) GetVendorsBrand(vendorsID int64) (*[]VendorsBrand, error) {
	var vendorsBrand []VendorsBrand
	err := r.conn.Select(&vendorsBrand, `
		SELECT vm.brand_id, m.name as brand_name 
		FROM vendor_brands vm LEFT JOIN brands m ON vm.brand_id = m.id 
		WHERE vm.vendor_id = ?
		AND vm.status > 0
	`, vendorsID)
	return &vendorsBrand, err
}

func (r *vendorsQuery) GetVendorsPicture(vendorsID int64) (*[]string, error) {
	var vendorsPicture []string
	err := r.conn.Select(&vendorsPicture, `
		SELECT name 
		FROM vendor_pictures 
		WHERE vendor_id = ?
		AND status > 0
	`, vendorsID)
	return &vendorsPicture, err
}

func (r *vendorsQuery) GetVendorsQrcode(vendorsID int64) (*[]VendorsQrcode, error) {
	var vendorsQrcode []VendorsQrcode
	err := r.conn.Select(&vendorsQrcode, `
		SELECT type,name 
		FROM vendor_qrcodes 
		WHERE vendor_id = ?
		AND status > 0
	`, vendorsID)
	return &vendorsQrcode, err
}
