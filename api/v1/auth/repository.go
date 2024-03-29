package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type authRepository struct {
	tx *sql.Tx
}

func NewAuthRepository(transaction *sql.Tx) *authRepository {
	return &authRepository{
		tx: transaction,
	}
}

func (r *authRepository) CreateUser(newUser User) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO users
		(
			type,
			identifier,
			organization_id,
			credential,
			birthday,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, 2, ?, "SIGNUP", ?, "SIGNUP")
	`, newUser.Type, newUser.Identifier, newUser.OrganizationID, newUser.Credential, newUser.Birthday, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if newUser.Type == 3 {
		_, err := r.tx.Exec(`
			INSERT INTO clients
			(
				user_id,
				organization_id,
				status,
				created,
				created_by,
				updated,
				updated_by
			)
			VALUES (?, ?, 2, ?, "SIGNUP", ?, "SIGNUP")
		`, id, newUser.OrganizationID, time.Now(), time.Now())
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *authRepository) GetUserByID(id int64) (*UserResponse, error) {
	var res UserResponse
	row := r.tx.QueryRow(`	
	SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.avatar as avatar, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
	FROM users u
	LEFT JOIN organizations o
	ON u.organization_id = o.id
	WHERE u.id = ?
	AND u.status > 0
	`, id)
	err := row.Scan(&res.ID, &res.Type, &res.Identifier, &res.OrganizationID, &res.PositionID, &res.RoleID, &res.Name, &res.Email, &res.Gender, &res.Phone, &res.Birthday, &res.Address, &res.Avatar, &res.Status, &res.OrganizationName)
	if err != nil {
		msg := "用户不存在:" + err.Error()
		return nil, errors.New(msg)
	}
	return &res, nil
}

func (r *authRepository) CheckConfict(authType int, identifier string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM users WHERE type = ? AND identifier = ?", authType, identifier)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}
func (r *authRepository) UpdateUser(id int64, info UserResponse, by string) error {
	_, err := r.tx.Exec(`
		Update users SET
		name = ?,
		email = ?,
		role_id = ?,
		position_id = ?, 
		gender = ?,
		phone = ?,
		birthday = ?,
		address = ?,
		avatar = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Email, info.RoleID, info.PositionID, info.Gender, info.Phone, info.Birthday, info.Address, info.Avatar, info.Status, time.Now(), by, id)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *authRepository) DeleteUser(id int64, by string) error {
	_, err := r.tx.Exec(`
		Update users SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, -1, time.Now(), by, id)
	return err
}

func (r *authRepository) CreateRole(info RoleNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO roles
		(
			name,
			priority,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Priority, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepository) UpdateRole(id int64, info RoleNew) (int64, error) {
	result, err := r.tx.Exec(`
		Update roles SET
		name = ?,
		priority = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Priority, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *authRepository) GetRoleByID(id int64) (*Role, error) {
	var res Role
	row := r.tx.QueryRow(`SELECT id, priority, name, status, created, created_by, updated, updated_by FROM roles WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Priority, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		msg := "角色不存在:" + err.Error()
		return nil, errors.New(msg)
	}
	return &res, nil
}

func (r *authRepository) CreateAPI(info APINew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO apis
		(
			name,
			route,
			method,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Route, info.Method, info.Status, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (r *authRepository) UpdateAPI(id int64, info APINew) (int64, error) {
	result, err := r.tx.Exec(`
		Update apis SET
		name = ?,
		route = ?,
		method = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Route, info.Method, info.Status, time.Now(), info.User, id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	return affected, err
}

func (r *authRepository) GetAPIByID(id int64) (*API, error) {
	var res API
	row := r.tx.QueryRow(`SELECT id, name, route, method, status, created, created_by, updated, updated_by FROM apis WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Route, &res.Method, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		msg := "API不存在:" + err.Error()
		return nil, errors.New(msg)
	}
	return &res, nil
}

func (r *authRepository) GetMenuByID(id int64) (*Menu, error) {
	var res Menu
	row := r.tx.QueryRow(`SELECT id, name, action, title, path, component, is_hidden, parent_id, status, created, created_by, updated, updated_by FROM menus WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Action, &res.Title, &res.Path, &res.Component, &res.IsHidden, &res.ParentID, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		msg := "菜单不存在:" + err.Error()
		return nil, errors.New(msg)
	}
	return &res, nil
}

func (r *authRepository) CreateMenu(info MenuNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO menus
		(
			name,
			action,
			title,
			path,
			component,
			is_hidden,
			parent_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *authRepository) UpdateMenu(id int64, info Menu, byUser string) error {
	fmt.Println(info.Component)
	_, err := r.tx.Exec(`
		Update menus SET
		name = ?,
		action = ?,
		title = ?,
		path = ?,
		component = ?,
		is_hidden = ?,
		parent_id = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, info.Status, time.Now(), byUser, id)
	return err
}

func (r *authRepository) DeleteMenu(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update menus SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *authRepository) NewRoleMenu(role_id int64, info RoleMenuNew) error {
	_, err := r.tx.Exec(`
		Update role_menus SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE role_id = ?
	`, time.Now(), info.User, role_id)
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO role_menus
	(
		role_id,
		menu_id,
		status,
		created,
		created_by,
		updated,
		updated_by
	)
	VALUES
	`
	for i := 0; i < len(info.IDS); i++ {
		sql += "(" + fmt.Sprint(role_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
	}
	sql = sql[:len(sql)-1]
	_, err = r.tx.Exec(sql)
	return err
}

func (r *authRepository) NewMenuAPI(menu_id int64, info MenuAPINew) error {
	_, err := r.tx.Exec(`
		Update menu_apis SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE menu_id = ?
	`, time.Now(), info.User, menu_id)
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO menu_apis
	(
		menu_id,
		api_id,
		status,
		created,
		created_by,
		updated,
		updated_by
	)
	VALUES
	`
	for i := 0; i < len(info.IDS); i++ {
		sql += "(" + fmt.Sprint(menu_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
	}
	sql = sql[:len(sql)-1]
	_, err = r.tx.Exec(sql)
	return err
}

func (r *authRepository) DeleteRole(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update roles SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *authRepository) UpdatePassword(id int64, password, by string) error {
	_, err := r.tx.Exec(`
		Update users SET
		credential = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, password, time.Now(), by, id)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (r *authRepository) GetWxmoduleByID(id int64) (*Wxmodule, error) {
	var res Wxmodule
	row := r.tx.QueryRow(`SELECT id, name, code, parent_id, status, created, created_by, updated, updated_by FROM wxmodules WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.Name, &res.Code, &res.ParentID, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		msg := "模块不存在:" + err.Error()
		return nil, errors.New(msg)
	}
	return &res, nil
}

func (r *authRepository) CreateWxmodule(info WxmoduleNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO wxmodules
		(
			name,
			code,
			parent_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.Name, info.Code, info.ParentID, 1, time.Now(), info.User, time.Now(), info.User)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *authRepository) UpdateWxmodule(id int64, info Wxmodule, byUser string) error {
	_, err := r.tx.Exec(`
		Update wxmodules SET
		name = ?,
		code = ?,
		parent_id = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Code, info.ParentID, info.Status, time.Now(), byUser, id)
	return err
}

func (r *authRepository) DeleteWxmodule(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update wxmodules SET 
		status = -1,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *authRepository) NewPositionWxmodule(position_id int64, info PositionWxmoduleNew) error {
	_, err := r.tx.Exec(`
		Update position_wxmodules SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE position_id = ?
	`, time.Now(), info.User, position_id)
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO position_wxmodules
	(
		position_id,
		wxmodule_id,
		status,
		created,
		created_by,
		updated,
		updated_by
	)
	VALUES
	`
	for i := 0; i < len(info.IDS); i++ {
		sql += "(" + fmt.Sprint(position_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
	}
	sql = sql[:len(sql)-1]
	_, err = r.tx.Exec(sql)
	return err
}

func (r authRepository) GetUserMemberCount(userID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(1) FROM project_members WHERE user_id = ? AND status > 0`, userID)
	err := row.Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r authRepository) GetUserClientCount(userID int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`SELECT count(id) FROM projects WHERE client_id = (SELECT id FROM clients WHERE user_id = ?) AND status > 0`, userID)
	err := row.Scan(&res)
	return res, err
}

func (r *authRepository) DeleteClient(userID int64, by string) error {
	_, err := r.tx.Exec(`
		Update clients SET
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE user_id = ?
	`, -1, time.Now(), by, userID)
	return err
}

func (r *authRepository) GetUserLimit(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`	
	SELECT user_limit
	FROM organizations 
	WHERE id = ?
	AND status > 0
	`, id)
	err := row.Scan(&res)
	if err != nil {
		msg := "获取用户数失败:" + err.Error()
		return 0, errors.New(msg)
	}
	return res, nil
}

func (r *authRepository) GetUserCount(id int64) (int, error) {
	var res int
	row := r.tx.QueryRow(`	
	SELECT count(1)
	FROM users 
	WHERE organization_id = ?
	AND type = 2
	AND status = 1
	`, id)
	err := row.Scan(&res)
	if err != nil {
		msg := "获取用户总数失败:" + err.Error()
		return 0, errors.New(msg)
	}
	return res, nil
}
