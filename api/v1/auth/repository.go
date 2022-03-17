package auth

import (
	"database/sql"
	"errors"
	"time"
)

type authRepository struct {
	tx *sql.Tx
}

func NewAuthRepository(transaction *sql.Tx) AuthRepository {
	return &authRepository{
		tx: transaction,
	}
}

type AuthRepository interface {
	// GetCredential(SigninRequest) (UserAuth, error)
	CreateUser(User) (int64, error)
	GetUserByID(int64) (*User, error)
	CheckConfict(int, string) (bool, error)
	UpdateUser(int64, User, string) error
	// GetAuthCount(filter AuthFilter) (int, error)
	// GetAuthList(filter AuthFilter) ([]Auth, error)

	// //Role Management
	CreateRole(info RoleNew) (int64, error)
	UpdateRole(int64, RoleNew) (int64, error)
	GetRoleByID(int64) (*Role, error)
	// //API Management
	CreateAPI(APINew) (int64, error)
	UpdateAPI(int64, APINew) (int64, error)
	GetAPIByID(int64) (*API, error)
	//Menu Management
	GetMenuByID(id int64) (*Menu, error)
	CreateMenu(info MenuNew) (int64, error)
	UpdateMenu(int64, Menu, string) error
	DeleteMenu(int64, string) error
	// //Privilege Management
	// GetRoleMenuByID(int64) ([]int64, error)
	// NewRoleMenu(int64, RoleMenuNew) (int64, error)
	// GetMenuAPIByID(int64) ([]int64, error)
	// NewMenuAPI(int64, MenuAPINew) (int64, error)
	// GetMyMenu(int64) ([]Menu, error)
}

func (r *authRepository) CreateUser(newUser User) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO users
		(
			type,
			identifier,
			organization_id,
			credential,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, 1, ?, "SIGNUP", ?, "SIGNUP")
	`, newUser.Type, newUser.Identifier, newUser.OrganizationID, newUser.Credential, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepository) GetUserByID(id int64) (*User, error) {
	var res User
	row := r.tx.QueryRow(`SELECT id, organization_id, type, identifier, credential, role_id, position_id, name, email, gender, phone, birthday, address, status, created, created_by, updated, updated_by FROM users WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Type, &res.Identifier, &res.Credential, &res.RoleID, &res.PositionID, &res.Name, &res.Email, &res.Gender, &res.Phone, &res.Birthday, &res.Address, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
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
func (r *authRepository) UpdateUser(id int64, info User, by string) error {
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
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Email, info.RoleID, info.PositionID, info.Gender, info.Phone, info.Birthday, info.Address, info.Status, time.Now(), by, id)
	if err != nil {
		msg := "更新失败:" + err.Error()
		return errors.New(msg)
	}
	return nil
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
	row := r.tx.QueryRow(`SELECT id, priority, name, status, created, created_by, updated, updated_by FROM roles WHERE id = ? LIMIT 1`, id)
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

// func (r *authRepository) GetRoleMenuByID(id int64) ([]int64, error) {
// 	var menu []int64
// 	err := r.conn.Select(&menu, "SELECT menu_id FROM user_role_menus WHERE role_id = ? and enabled = 1", id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return menu, nil
// }
// func (r *authRepository) NewRoleMenu(role_id int64, info RoleMenuNew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	_, err = tx.Exec(`
// 		Update user_role_menus SET
// 		enabled = 2,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE role_id = ?
// 		AND enabled = 1
// 	`, time.Now(), info.User, role_id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sql := `
// 	INSERT INTO user_role_menus
// 	(
// 		role_id,
// 		menu_id,
// 		enabled,
// 		created,
// 		created_by,
// 		updated,
// 		updated_by
// 	)
// 	VALUES
// 	`
// 	for i := 0; i < len(info.IDS); i++ {
// 		sql += "(" + fmt.Sprint(role_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
// 	}
// 	sql = sql[:len(sql)-1]
// 	result, err := tx.Exec(sql)
// 	if err != nil {
// 		return 0, err
// 	}
// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return rows, nil
// }

// func (r *authRepository) GetMenuAPIByID(id int64) ([]int64, error) {
// 	var apis []int64
// 	err := r.conn.Select(&apis, "SELECT api_id FROM user_menu_apis WHERE menu_id = ? and enabled = 1", id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return apis, nil
// }
// func (r *authRepository) NewMenuAPI(menu_id int64, info MenuAPINew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	_, err = tx.Exec(`
// 		Update user_menu_apis SET
// 		enabled = 2,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE menu_id = ?
// 		AND enabled = 1
// 	`, time.Now(), info.User, menu_id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	sql := `
// 	INSERT INTO user_menu_apis
// 	(
// 		menu_id,
// 		api_id,
// 		enabled,
// 		created,
// 		created_by,
// 		updated,
// 		updated_by
// 	)
// 	VALUES
// 	`
// 	for i := 0; i < len(info.IDS); i++ {
// 		sql += "(" + fmt.Sprint(menu_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
// 	}
// 	sql = sql[:len(sql)-1]
// 	result, err := tx.Exec(sql)
// 	if err != nil {
// 		return 0, err
// 	}
// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return rows, nil
// }
// func (r *authRepository) GetMyMenu(roleID int64) ([]Menu, error) {
// 	var menu []Menu
// 	err := r.conn.Select(&menu, `
// 		SELECT um.* FROM user_role_menus urm
// 		LEFT JOIN user_menus um
// 		ON urm.menu_id = um.id
// 		WHERE urm.role_id = ?
// 		AND um.enabled = 1
// 		AND urm.enabled = 1
// 		ORDER BY parent_id ASC, ID ASC
// 	`, roleID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return menu, nil
// }
