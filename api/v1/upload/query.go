package upload

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type uploadQuery struct {
	conn *sqlx.DB
}

func NewUploadQuery(connection *sqlx.DB) UploadQuery {
	return &uploadQuery{
		conn: connection,
	}
}

type UploadQuery interface {
	//Upload Management
	GetUploadByID(id int64) (*Upload, error)
	GetUploadCount(UploadFilter) (int, error)
	GetUploadList(UploadFilter) (*[]Upload, error)
}

func (r *uploadQuery) GetUploadByID(id int64) (*Upload, error) {
	var upload Upload
	err := r.conn.Get(&upload, "SELECT * FROM file_uploads WHERE id = ? AND status > 0 ", id)
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *uploadQuery) GetUploadCount(filter UploadFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "created_by = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM file_uploads
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *uploadQuery) GetUploadList(filter UploadFilter) (*[]Upload, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "created_by = ?"), append(args, v)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var uploads []Upload
	err := r.conn.Select(&uploads, `
		SELECT *
		FROM file_uploads 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &uploads, nil
}
