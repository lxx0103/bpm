package auth

import (
	"database/sql"
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
	UpdateUser(int64, UserUpdate) error
	// GetAuthCount(filter AuthFilter) (int, error)
	// GetAuthList(filter AuthFilter) ([]Auth, error)

	// //Role Management
	CreateRole(info RoleNew) (int64, error)
	UpdateRole(id int64, info RoleNew) (int64, error)
	GetRoleByID(int64) (*Role, error)
	// //API Management
	// GetAPIByID(id int64) (UserAPI, error)
	// CreateAPI(info APINew) (int64, error)
	// GetAPICount(filter APIFilter) (int, error)
	// GetAPIList(filter APIFilter) ([]UserAPI, error)
	// UpdateAPI(id int64, info APINew) (int64, error)
	// //Menu Management
	// GetMenuByID(id int64) (UserMenu, error)
	// CreateMenu(info MenuNew) (int64, error)
	// GetMenuCount(filter MenuFilter) (int, error)
	// GetMenuList(filter MenuFilter) ([]UserMenu, error)
	// UpdateMenu(id int64, info MenuNew) (int64, error)
	// //Privilege Management
	// GetRoleMenuByID(int64) ([]int64, error)
	// NewRoleMenu(int64, RoleMenuNew) (int64, error)
	// GetMenuAPIByID(int64) ([]int64, error)
	// NewMenuAPI(int64, MenuAPINew) (int64, error)
	// GetMyMenu(int64) ([]UserMenu, error)
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
	row := r.tx.QueryRow(`SELECT id, organization_id, type, identifier, credential, role_id, name, email, gender, phone, birthday, address, status, created, created_by, updated, updated_by FROM users WHERE id = ? LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Type, &res.Identifier, &res.Credential, &res.RoleID, &res.Name, &res.Email, &res.Gender, &res.Phone, &res.Birthday, &res.Address, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
	if err != nil {
		return nil, err
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
func (r *authRepository) UpdateUser(id int64, info UserUpdate) error {
	_, err := r.tx.Exec(`
		Update user_profiles SET 
		name = ?,
		email = ?,
		role_id = ?,
		gender = ?,
		phone = ?,
		birthday = ?,
		address = ?,
		status = ?,
		updated = ?,
		updated_by = ? 
		WHERE id = ?
	`, info.Name, info.Email, info.RoleID, info.Gender, info.Phone, info.Birthday, info.Address, info.Status, time.Now(), info.User, id)
	if err != nil {
		return err
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
		return nil, err
	}
	return &res, nil
}

// func (r *authRepository) GetAPIByID(id int64) (UserAPI, error) {
// 	var api UserAPI
// 	err := r.conn.Get(&api, "SELECT * FROM user_apis WHERE id = ? ", id)
// 	if err != nil {
// 		return UserAPI{}, err
// 	}
// 	return api, nil
// }
// func (r *authRepository) CreateAPI(info APINew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	result, err := tx.Exec(`
// 		INSERT INTO user_apis
// 		(
// 			name,
// 			route,
// 			method,
// 			enabled,
// 			created,
// 			created_by,
// 			updated,
// 			updated_by
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
// 	`, info.Name, info.Route, info.Method, info.Enabled, time.Now(), info.User, time.Now(), info.User)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return id, nil
// }

// func (r *authRepository) GetAPICount(filter APIFilter) (int, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Route; v != "" {
// 		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM user_apis
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authRepository) GetAPIList(filter APIFilter) ([]UserAPI, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Route; v != "" {
// 		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var apis []UserAPI
// 	err := r.conn.Select(&apis, `
// 		SELECT *
// 		FROM user_apis
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return apis, nil
// }

// func (r *authRepository) UpdateAPI(id int64, info APINew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	result, err := tx.Exec(`
// 		Update user_apis SET
// 		name = ?,
// 		route = ?,
// 		method = ?,
// 		enabled = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Route, info.Method, info.Enabled, time.Now(), info.User, id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return affected, nil
// }

// func (r *authRepository) GetMenuByID(id int64) (UserMenu, error) {
// 	var menu UserMenu
// 	err := r.conn.Get(&menu, "SELECT * FROM user_menus WHERE id = ? ", id)
// 	if err != nil {
// 		return UserMenu{}, err
// 	}
// 	return menu, nil
// }
// func (r *authRepository) CreateMenu(info MenuNew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	result, err := tx.Exec(`
// 		INSERT INTO user_menus
// 		(
// 			name,
// 			action,
// 			title,
// 			path,
// 			component,
// 			is_hidden,
// 			parent_id,
// 			enabled,
// 			created,
// 			created_by,
// 			updated,
// 			updated_by
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, info.Enabled, time.Now(), info.User, time.Now(), info.User)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return id, nil
// }

// func (r *authRepository) GetMenuCount(filter MenuFilter) (int, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.OnlyTop; v {
// 		where, args = append(where, "parent_id = ?"), append(args, 0)
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM user_menus
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authRepository) GetMenuList(filter MenuFilter) ([]UserMenu, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.OnlyTop; v {
// 		where, args = append(where, "parent_id = ?"), append(args, 0)
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var menus []UserMenu
// 	err := r.conn.Select(&menus, `
// 		SELECT *
// 		FROM user_menus
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return menus, nil
// }

// func (r *authRepository) UpdateMenu(id int64, info MenuNew) (int64, error) {
// 	tx, err := r.conn.Begin()
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer tx.Rollback()
// 	result, err := tx.Exec(`
// 		Update user_menus SET
// 		name = ?,
// 		action = ?,
// 		title = ?,
// 		path = ?,
// 		component = ?,
// 		is_hidden = ?,
// 		parent_id = ?,
// 		enabled = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, info.Enabled, time.Now(), info.User, id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	tx.Commit()
// 	return affected, nil
// }

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
// func (r *authRepository) GetMyMenu(roleID int64) ([]UserMenu, error) {
// 	var menu []UserMenu
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
