package organization

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type organizationQuery struct {
	conn *sqlx.DB
}

func NewOrganizationQuery(connection *sqlx.DB) OrganizationQuery {
	return &organizationQuery{
		conn: connection,
	}
}

type OrganizationQuery interface {
	//Organization Management
	GetOrganizationByID(id int64) (*Organization, error)
	GetOrganizationCount(filter OrganizationFilter) (int, error)
	GetOrganizationList(filter OrganizationFilter) (*[]Organization, error)
}

func (r *organizationQuery) GetOrganizationByID(id int64) (*Organization, error) {
	var organization Organization
	err := r.conn.Get(&organization, "SELECT * FROM organizations WHERE id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *organizationQuery) GetOrganizationCount(filter OrganizationFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
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

func (r *organizationQuery) GetOrganizationList(filter OrganizationFilter) (*[]Organization, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var organizations []Organization
	err := r.conn.Select(&organizations, `
		SELECT * 
		FROM organizations 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}
