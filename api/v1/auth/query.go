package auth

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type authQuery struct {
	conn *sqlx.DB
}

func NewAuthQuery(connection *sqlx.DB) *authQuery {
	return &authQuery{
		conn: connection,
	}
}

func (r *authQuery) GetUserByID(id int64, organizationID int64) (*User, error) {
	var user User
	var err error
	if organizationID != 0 {
		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? AND organization_id = ? AND status > 0 ", id, organizationID)
	} else {
		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? AND status > 0", id)
	}
	if err != nil {
		return nil, err
	}
	user.Credential = ""
	return &user, nil
}

func (r *authQuery) GetUserByOpenID(openID string) (*UserResponse, error) {
	var user UserResponse
	err := r.conn.Get(&user, `	
		SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.avatar as avatar, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
		FROM users u
		LEFT JOIN organizations o
		ON u.organization_id = o.id
		WHERE identifier = ?
		AND u.status > 0 
	`, openID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authQuery) GetUserCredential(id int64) (string, error) {
	var credential string
	err := r.conn.Get(&credential, "SELECT credential FROM users WHERE id = ? ", id)
	if err != nil {
		return "", err
	}
	return credential, nil
}

func (r *authQuery) GetUserCount(filter UserFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v == "wx" {
		where, args = append(where, "type = ?"), append(args, 2)
	}
	if v := filter.Type; v == "admin" {
		where, args = append(where, "type = ?"), append(args, 1)
	}
	if v := filter.OrganizationID; v != 0 {
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

func (r *authQuery) GetUserList(filter UserFilter) (*[]UserResponse, error) {
	where, args := []string{"u.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "u.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.Type; v == "wx" {
		where, args = append(where, "u.type = ?"), append(args, 2)
	}
	if v := filter.Type; v == "admin" {
		where, args = append(where, "u.type = ?"), append(args, 1)
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "u.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var users []UserResponse
	err := r.conn.Select(&users, `
		SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.avatar as avatar, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
		FROM users u
		LEFT JOIN organizations o
		ON u.organization_id = o.id
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
	err := r.conn.Get(&role, "SELECT * FROM roles WHERE id = ? AND status > 0", id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
func (r *authQuery) GetRoleCount(filter RoleFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
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
	where, args := []string{"status > 0"}, []interface{}{}
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
	where, args := []string{"status > 0"}, []interface{}{}
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
	where, args := []string{"status > 0"}, []interface{}{}
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

func (r *authQuery) GetRoleMenuByID(roleID int64) ([]int64, error) {
	var menu []int64
	err := r.conn.Select(&menu, "SELECT menu_id FROM role_menus WHERE role_id = ? and status > 0", roleID)
	return menu, err
}

func (r *authQuery) GetMyMenu(roleID int64) ([]Menu, error) {
	var menu []Menu
	err := r.conn.Select(&menu, `
		SELECT m.* FROM role_menus rm
		LEFT JOIN menus m
		ON rm.menu_id = m.id
		WHERE rm.role_id = ?
		AND m.status > 0
		AND rm.status > 0
		ORDER BY parent_id ASC, ID ASC
	`, roleID)
	return menu, err
}

func (r *authQuery) GetWxmoduleByID(id int64) (*Wxmodule, error) {
	var wxmodule Wxmodule
	err := r.conn.Get(&wxmodule, "SELECT * FROM wxmodules WHERE id = ? AND status > 0 ", id)
	return &wxmodule, err
}

func (r *authQuery) GetWxmoduleCount(filter WxmoduleFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM wxmodules
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *authQuery) GetWxmoduleList(filter WxmoduleFilter) (*[]Wxmodule, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var wxmodules []Wxmodule
	err := r.conn.Select(&wxmodules, `
		SELECT *
		FROM wxmodules
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &wxmodules, err
}

func (r *authQuery) GetPositionWxmoduleByID(positionID int64) ([]int64, error) {
	var wxmodule []int64
	err := r.conn.Select(&wxmodule, "SELECT wxmodule_id FROM position_wxmodules WHERE position_id = ? and status > 0", positionID)
	return wxmodule, err
}

func (r *authQuery) GetMyWxmodule(positionID, parentID int64) ([]Wxmodule, error) {
	var wxmodule []Wxmodule
	err := r.conn.Select(&wxmodule, `
		SELECT m.* FROM position_wxmodules rm
		LEFT JOIN wxmodules m
		ON rm.wxmodule_id = m.id
		WHERE rm.position_id = ?
		AND m.status > 0
		AND rm.status > 0
		AND m.parent_id = ?
		ORDER BY parent_id ASC, ID ASC
	`, positionID, parentID)
	return wxmodule, err
}
