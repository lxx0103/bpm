package organization

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type organizationQuery struct {
	conn *sqlx.DB
}

func NewOrganizationQuery(connection *sqlx.DB) *organizationQuery {
	return &organizationQuery{
		conn: connection,
	}
}

func (r *organizationQuery) GetOrganizationByID(id int64) (*OrganizationResponse, error) {
	var organization OrganizationResponse
	err := r.conn.Get(&organization, "SELECT id, name, logo, description, phone, contact, address, city, type, status FROM organizations WHERE id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *organizationQuery) GetOrganizationCount(filter OrganizationFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.City; v != "" {
		where, args = append(where, "city like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM organizations 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *organizationQuery) GetOrganizationList(filter OrganizationFilter) (*[]OrganizationResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.City; v != "" {
		where, args = append(where, "city like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v != 0 {
		where, args = append(where, "type = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var organizations []OrganizationResponse
	err := r.conn.Select(&organizations, `
		SELECT id, name, logo, description, phone, contact, address, city, type, status
		FROM organizations 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}

func (r *organizationQuery) GetQrCodeByPath(path, source string) (string, error) {
	var res string
	err := r.conn.Get(&res, "SELECT img FROM qr_codes WHERE path = ? AND source = ?", path, source)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *organizationQuery) GetAccessToken(code string) (string, error) {
	var res string
	err := r.conn.Get(&res, "SELECT access_token FROM wx_access_token WHERE code = ? AND expires_in > DATE_ADD(now(), INTERVAL 5 MINUTE) order by id desc", code)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *organizationQuery) GetOrganizationTopExamples(organizationID int64) (*[]ExampleResponse, error) {
	var examples []ExampleResponse
	err := r.conn.Select(&examples, `
		SELECT id, name, cover, status
		FROM examples 
		WHERE organization_id = ?
		AND example_type = 1
		AND status > 0 
		ORDER BY id desc
		LIMIT 3
	`, organizationID)
	return &examples, err
}

func (r *organizationQuery) GetOrganizationQrcode(organizationID int64) (*[]OrganizationQrcode, error) {
	var organizationQrcode []OrganizationQrcode
	err := r.conn.Select(&organizationQrcode, `
		SELECT type,name 
		FROM organization_qrcodes 
		WHERE organization_id = ?
		AND status > 0
	`, organizationID)
	return &organizationQrcode, err
}
