package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bpm/core/config"
	"bpm/core/database"
)

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

type AuthService interface {
	VerifyWechatSignin(string) (*WechatCredential, error)
	GetUserInfo(string) (*User, error)
	//Role Management
	GetRoleByID(int64) (*Role, error)
	NewRole(RoleNew) (*Role, error)
	GetRoleList(RoleFilter) (int, *[]Role, error)
	UpdateRole(int64, RoleNew) (*Role, error)
	// //API Management
	// GetAPIByID(int64) (UserAPI, error)
	// NewAPI(APINew) (UserAPI, error)
	// GetAPIList(APIFilter) (int, []UserAPI, error)
	// UpdateAPI(int64, APINew) (UserAPI, error)
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
		repo := NewAuthRepository(tx)
		userID, err := repo.CreateUser(openID)
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

// func (s *authService) GetAPIByID(id int64) (UserAPI, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	api, err := repo.GetAPIByID(id)
// 	return api, err
// }

// func (s *authService) NewAPI(info APINew) (UserAPI, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	apiID, err := repo.CreateAPI(info)
// 	if err != nil {
// 		return UserAPI{}, err
// 	}
// 	api, err := repo.GetAPIByID(apiID)
// 	return api, err
// }

// func (s *authService) GetAPIList(filter APIFilter) (int, []UserAPI, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	count, err := repo.GetAPICount(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	list, err := repo.GetAPIList(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	return count, list, err
// }

// func (s *authService) UpdateAPI(apiID int64, info APINew) (UserAPI, error) {
// 	db := database.InitMySQL()
// 	repo := NewAuthRepository(db)
// 	_, err := repo.UpdateAPI(apiID, info)
// 	if err != nil {
// 		return UserAPI{}, err
// 	}
// 	api, err := repo.GetAPIByID(apiID)
// 	return api, err
// }

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
