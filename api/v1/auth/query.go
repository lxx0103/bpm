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
	// //API Management
	GetAPIByID(id int64) (*API, error)
	GetAPICount(filter APIFilter) (int, error)
	GetAPIList(filter APIFilter) (*[]API, error)
	//Menu Management
	GetMenuByID(id int64) (*Menu, error)
	GetMenuCount(filter MenuFilter) (int, error)
	GetMenuList(filter MenuFilter) (*[]Menu, error)
	GetMenuAPIByID(int64) ([]int64, error)
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
	err := r.conn.Get(&user, "SELECT * FROM users WHERE identifier = ? ", openID)
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

func (r *authQuery) GetAPIByID(id int64) (*API, error) {
	var api API
	err := r.conn.Get(&api, "SELECT * FROM apis WHERE id = ? ", id)
	return &api, err
}

func (r *authQuery) GetAPICount(filter APIFilter) (int, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Route; v != "" {
		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM apis
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetAPIList(filter APIFilter) (*[]API, error) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Route; v != "" {
		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var apis []API
	err := r.conn.Select(&apis, `
		SELECT *
		FROM apis
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &apis, nil
}

func (r *authQuery) GetMenuByID(id int64) (*Menu, error) {
	var menu Menu
	err := r.conn.Get(&menu, "SELECT * FROM menus WHERE id = ? AND status > 0 ", id)
	return &menu, err
}

func (r *authQuery) GetMenuCount(filter MenuFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OnlyTop; v {
		where, args = append(where, "parent_id = ?"), append(args, 0)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM menus
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetMenuList(filter MenuFilter) (*[]Menu, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OnlyTop; v {
		where, args = append(where, "parent_id = ?"), append(args, 0)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var menus []Menu
	err := r.conn.Select(&menus, `
		SELECT *
		FROM menus
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &menus, err
}

func (r *authQuery) GetMenuAPIByID(menuID int64) ([]int64, error) {
	var apis []int64
	err := r.conn.Select(&apis, "SELECT api_id FROM menu_apis WHERE menu_id = ? and status > 0", menuID)
	return apis, err
}
