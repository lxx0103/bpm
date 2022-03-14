package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"bpm/core/config"
	"bpm/core/database"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

type AuthService interface {
	CreateAuth(SignupRequest) (int64, error)
	VerifyWechatSignin(string) (*WechatCredential, error)
	VerifyCredential(SigninRequest) (*User, error)
	//User Management
	GetUserInfo(string) (*User, error)
	UpdateUser(int64, UserUpdate, int64) (*User, error)
	GetUserByID(int64, int64) (*User, error)
	GetUserList(UserFilter, int64) (int, *[]User, error)
	//Role Management
	GetRoleByID(int64) (*Role, error)
	NewRole(RoleNew) (*Role, error)
	GetRoleList(RoleFilter) (int, *[]Role, error)
	UpdateRole(int64, RoleNew) (*Role, error)
	// //API Management
	GetAPIByID(int64) (*API, error)
	NewAPI(APINew) (*API, error)
	GetAPIList(APIFilter) (int, *[]API, error)
	UpdateAPI(int64, APINew) (*API, error)
	// //Menu Management
	// GetMenuByID(int64) (UserMenu, error)
	// NewMenu(MenuNew) (UserMenu, error)
	// GetMenuList(MenuFilter) (int, []UserMenu, error)
	// UpdateMenu(int64, MenuNew) (UserMenu, error)
	// //Privilege Management
	// GetRoleMenuByID(int64) ([]int64, error)
	// NewRoleMenu(int64, RoleMenuNew) ([]int64, error)
	// GetMenuAPIByID(int64) ([]int64, error)
	// NewMenuAPI(int64, MenuAPINew) ([]int64, error)
	// GetMyMenu(int64) ([]UserMenu, error)
}

func (s authService) CreateAuth(signupInfo SignupRequest) (int64, error) {
	hashed, err := hashPassword(signupInfo.Credential)
	if err != nil {
		return 0, err
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	var newUser User
	newUser.Credential = hashed
	isConflict, err := repo.CheckConfict(1, signupInfo.Identifier)
	if err != nil {
		return 0, err
	}
	if isConflict {
		errMessage := "用户名已存在"
		return 0, errors.New(errMessage)
	}
	newUser.Identifier = signupInfo.Identifier
	newUser.Type = 1
	newUser.OrganizationID = signupInfo.OrganizationID
	authID, err := repo.CreateUser(newUser)
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return authID, nil
}

func (s *authService) VerifyWechatSignin(code string) (*WechatCredential, error) {
	var credential WechatCredential
	httpClient := &http.Client{}
	signin_uri := config.ReadConfig("Wechat.signin_uri")
	appID := config.ReadConfig("Wechat.app_id")
	appSecret := config.ReadConfig("Wechat.app_secret")
	uri := signin_uri + "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

func (s *authService) GetUserInfo(openID string) (*User, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	user, err := query.GetUserByOpenID(openID)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, err
		}
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
		var newUser User
		newUser.Type = 2
		newUser.Identifier = openID
		repo := NewAuthRepository(tx)
		userID, err := repo.CreateUser(newUser)
		if err != nil {
			return nil, err
		}
		user, err = repo.GetUserByID(userID)
		if err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return user, nil
}

func (s *authService) VerifyCredential(signinInfo SigninRequest) (*User, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	userInfo, err := query.GetUserByUserName(signinInfo.Identifier)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(signinInfo.Credential, userInfo.Credential) {
		errMessage := "密码错误"
		return nil, errors.New(errMessage)
	}
	return userInfo, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) UpdateUser(userID int64, info UserUpdate, byUserID int64) (*User, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		msg := "事务开启错误" + err.Error()
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)

	oldUser, err := repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	byUser, err := repo.GetUserByID(byUserID)
	if err != nil {
		return nil, err
	}
	byRole, err := repo.GetRoleByID(byUser.RoleID)
	if err != nil {
		return nil, err
	}
	if oldUser.RoleID != 0 {
		targetRole, err := repo.GetRoleByID(oldUser.RoleID)
		if err != nil {
			return nil, err
		}
		if byRole.Priority <= targetRole.Priority && userID != byUserID { //只能修改角色比自己优先级低的用户,或者用户自身
			msg := "你无法修改角色为" + targetRole.Name + "的用户"
			return nil, errors.New(msg)
		}
	}
	if info.RoleID != 0 {
		toRole, err := repo.GetRoleByID(info.RoleID)
		if err != nil {
			return nil, err
		}
		if byRole.Priority < toRole.Priority { //只能将目标修改为和自己同级的角色
			msg := "你无法将目标角色改为:" + toRole.Name
			return nil, errors.New(msg)
		}
		oldUser.RoleID = info.RoleID
	}
	if info.PositionID != 0 {
		oldUser.PositionID = info.PositionID
	}
	if info.Name != "" {
		oldUser.Name = info.Name
	}
	if info.Email != "" {
		oldUser.Email = info.Email
	}
	if info.Gender != "" {
		oldUser.Gender = info.Gender
	}
	if info.Birthday != "" {
		oldUser.Birthday = info.Birthday
	}
	if info.Phone != "" {
		oldUser.Phone = info.Phone
	}
	if info.Address != "" {
		oldUser.Address = info.Address
	}
	if info.Status != 0 {
		oldUser.Status = info.Status
	}
	err = repo.UpdateUser(userID, *oldUser, (*byUser).Identifier)
	if err != nil {
		return nil, err
	}
	user, err := repo.GetUserByID(userID)
	tx.Commit()
	return user, err
}

func (s *authService) GetRoleByID(id int64) (*Role, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	role, err := query.GetRoleByID(id)
	return role, err
}

func (s *authService) NewRole(info RoleNew) (*Role, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	roleID, err := repo.CreateRole(info)
	if err != nil {
		return nil, err
	}
	role, err := repo.GetRoleByID(roleID)
	tx.Commit()
	return role, err
}

func (s *authService) GetRoleList(filter RoleFilter) (int, *[]Role, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	count, err := query.GetRoleCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetRoleList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *authService) UpdateRole(roleID int64, info RoleNew) (*Role, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	_, err = repo.UpdateRole(roleID, info)
	if err != nil {
		return nil, err
	}
	role, err := repo.GetRoleByID(roleID)
	tx.Commit()
	return role, err
}

func (s *authService) GetUserByID(id int64, organizationID int64) (*User, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	user, err := query.GetUserByID(id, organizationID)
	return user, err
}

func (s *authService) GetUserList(filter UserFilter, organizationID int64) (int, *[]User, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	count, err := query.GetUserCount(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetUserList(filter, organizationID)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *authService) GetAPIByID(id int64) (*API, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	api, err := query.GetAPIByID(id)
	return api, err
}

func (s *authService) GetAPIList(filter APIFilter) (int, *[]API, error) {
	db := database.InitMySQL()
	query := NewAuthQuery(db)
	count, err := query.GetAPICount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetAPIList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *authService) NewAPI(info APINew) (*API, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	apiID, err := repo.CreateAPI(info)
	if err != nil {
		return nil, err
	}
	api, err := repo.GetAPIByID(apiID)
	tx.Commit()
	return api, err
}

func (s *authService) UpdateAPI(apiID int64, info APINew) (*API, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	_, err = repo.UpdateAPI(apiID, info)
	if err != nil {
		return nil, err
	}
	api, err := repo.GetAPIByID(apiID)
	tx.Commit()
	return api, err
}

// func (s *authService) GetMenuByID(id int64) (UserMenu, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	menu, err := repo.GetMenuByID(id)
// 	return menu, err
// }

// func (s *authService) NewMenu(info MenuNew) (UserMenu, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	menuID, err := repo.CreateMenu(info)
// 	if err != nil {
// 		return UserMenu{}, err
// 	}
// 	menu, err := repo.GetMenuByID(menuID)
// 	return menu, err
// }

// func (s *authService) GetMenuList(filter MenuFilter) (int, []UserMenu, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	count, err := repo.GetMenuCount(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	list, err := repo.GetMenuList(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	return count, list, err
// }

// func (s *authService) UpdateMenu(menuID int64, info MenuNew) (UserMenu, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	_, err := repo.UpdateMenu(menuID, info)
// 	if err != nil {
// 		return UserMenu{}, err
// 	}
// 	menu, err := repo.GetMenuByID(menuID)
// 	return menu, err
// }

// func (s *authService) GetRoleMenuByID(id int64) ([]int64, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	menu, err := repo.GetRoleMenuByID(id)
// 	return menu, err
// }

// func (s *authService) NewRoleMenu(id int64, info RoleMenuNew) ([]int64, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	_, err := repo.NewRoleMenu(id, info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	menu, err := repo.GetRoleMenuByID(id)
// 	return menu, err
// }

// func (s *authService) GetMenuAPIByID(id int64) ([]int64, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	menu, err := repo.GetMenuAPIByID(id)
// 	return menu, err
// }

// func (s *authService) NewMenuAPI(id int64, info MenuAPINew) ([]int64, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	_, err := repo.NewMenuAPI(id, info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	menu, err := repo.GetMenuAPIByID(id)
// 	return menu, err
// }

// func (s *authService) GetMyMenu(roleID int64) ([]UserMenu, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	menu, err := repo.GetMyMenu(roleID)
// 	return menu, err
// }
