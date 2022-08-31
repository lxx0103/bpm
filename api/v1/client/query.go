package client

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type clientQuery struct {
	conn *sqlx.DB
}

func NewClientQuery(connection *sqlx.DB) ClientQuery {
	return &clientQuery{
		conn: connection,
	}
}

type ClientQuery interface {
	//Client Management
	GetClientByID(int64, int64) (*Client, error)
	GetClientCount(ClientFilter, int64) (int, error)
	GetClientList(ClientFilter, int64) (*[]Client, error)
	GetClientByUserID(int64, int64) (*Client, error)
}

func (r *clientQuery) GetClientByID(id int64, organizationID int64) (*Client, error) {
	var client Client
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&client, "SELECT * FROM clients WHERE id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&client, "SELECT * FROM clients WHERE id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientQuery) GetClientCount(filter ClientFilter, organizationID int64) (int, error) {
	if organizationID == 0 && filter.OrganizationID != 0 {
		organizationID = filter.OrganizationID
	}
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count 
		FROM clients 
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *clientQuery) GetClientList(filter ClientFilter, organizationID int64) (*[]Client, error) {
	if organizationID == 0 && filter.OrganizationID != 0 {
		organizationID = filter.OrganizationID
	}
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var clients []Client
	err := r.conn.Select(&clients, `
		SELECT * 
		FROM clients 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &clients, nil
}

func (r *clientQuery) GetClientByUserID(id int64, organizationID int64) (*Client, error) {
	var client Client
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&client, "SELECT * FROM clients WHERE user_id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&client, "SELECT * FROM clients WHERE user_id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}
