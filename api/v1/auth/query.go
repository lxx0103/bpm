package auth

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type authQuery struct {
	conn *sqlx.DB
}

func NewAuthQuery(connection *sqlx.DB) AuthQuery {
	return &authQuery{
		conn: connection,
	}
}

type AuthQuery interface {
	//User Management
	GetUserByID(int64, int64) (*User, error)
	GetUserByOpenID(openID string) (*User, error)
	GetUserByUserName(userName string) (*User, error)
	GetUserCount(UserFilter, int64) (int, error)
	GetUserList(UserFilter, int64) (*[]User, error)
	//Role Management
	GetRoleByID(id int64) (*Role, error)
	GetRoleCount(filter RoleFilter) (int, error)
	GetRoleList(filter RoleFilter) (*[]Role, error)
}

func (r *authQuery) GetUserByID(id int64, organizationID int64) (*User, error) {
	var user User
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? AND organization_id = ?", id, organizationID)
	} else {
		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? ", id)
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authQuery) GetUserByOpenID(openID string) (*User, error) {
	var user User
	err := r.conn.Get(&user, "SELECT * FROM users WHERE open_id = ? ", openID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authQuery) GetUserByUserName(userName string) (*User, error) {
	var user User
	err := r.conn.Get(&user, "SELECT * FROM users WHERE identifier = ? ", userName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authQuery) GetUserCount(filter UserFilter, organizationID int64) (int, error) {
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
		FROM users
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetUserList(filter UserFilter, organizationID int64) (*[]User, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := organizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var users []User
	err := r.conn.Select(&users, `
		SELECT *
		FROM users
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *authQuery) GetRoleByID(id int64) (*Role, error) {
	var role Role
	err := r.conn.Get(&role, "SELECT * FROM roles WHERE id = ? ", id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *authQuery) GetRoleCount(filter RoleFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM roles
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetRoleList(filter RoleFilter) (*[]Role, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var roles []Role
	err := r.conn.Select(&roles, `
		SELECT *
		FROM roles
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &roles, nil
}
