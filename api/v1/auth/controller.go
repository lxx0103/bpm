package auth

import (
	"bpm/core/response"
	"bpm/service"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// @Summary 登录
// @Id 17
// @Tags 用户权限
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param signin_info body SigninRequest true "登录类型"
// @Success 200 object response.SuccessRes{data=SigninResponse} 登录成功
// @Failure 400 object response.ErrorRes 内部错误
// @Failure 401 object response.ErrorRes 登录失败
// @Router /signin [POST]
func Signin(c *gin.Context) {
	var signinInfo SigninRequest
	var userInfo *User
	err := c.ShouldBindJSON(&signinInfo)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	authService := NewAuthService()
	if signinInfo.AuthType == 2 {
		wechatCredential, err := authService.VerifyWechatSignin(signinInfo.Identifier)
		if err != nil {
			response.ResponseUnauthorized(c, "AuthError", err)
			return
		}
		if wechatCredential.ErrCode != 0 {
			response.ResponseUnauthorized(c, "AuthError", errors.New(wechatCredential.ErrMsg))
			return
		}
		userInfo, err = authService.GetUserInfo(wechatCredential.OpenID)
		if err != nil {
			response.ResponseUnauthorized(c, "AuthError", err)
			return
		}
	} else if signinInfo.AuthType == 1 {
		userInfo, err = authService.VerifyCredential(signinInfo)
		if err != nil {
			response.ResponseUnauthorized(c, "AuthError", err)
			return
		}
	} else {
		errMessage := "登陆类型错误"
		response.ResponseUnauthorized(c, "AuthError", errors.New(errMessage))
		return
	}
	claims := service.CustomClaims{
		UserID:         userInfo.ID,
		Username:       userInfo.Identifier,
		RoleID:         userInfo.RoleID,
		OrganizationID: userInfo.OrganizationID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + 72000,
			Issuer:    "bpm",
		},
	}
	jwtServices := service.JWTAuthService()
	generatedToken := jwtServices.GenerateToken(claims)
	var res SigninResponse
	res.Token = generatedToken
	res.User = *userInfo
	response.Response(c, res)
}

// @Id 22
// @Tags 用户权限
// @Summary 用户注册
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param signup_info body SignupRequest true "登录类型"
// @Success 200 object response.SuccessRes{data=int} 注册成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /signup [POST]
func Signup(c *gin.Context) {
	var signupInfo SignupRequest
	err := c.ShouldBindJSON(&signupInfo)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	authService := NewAuthService()
	authID, err := authService.CreateAuth(signupInfo)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, authID)
}

// @Summary 角色列表
// @Id 18
// @Tags 角色管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "角色名称"
// @Success 200 object response.ListRes{data=[]Role} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /roles [GET]
func GetRoleList(c *gin.Context) {
	var filter RoleFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	authService := NewAuthService()
	count, list, err := authService.GetRoleList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建角色
// @Id 19
// @Tags 角色管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param role_info body RoleNew true "角色信息"
// @Success 200 object response.SuccessRes{data=Role} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /roles [POST]
func NewRole(c *gin.Context) {
	var role RoleNew
	if err := c.ShouldBindJSON(&role); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	role.User = claims.Username
	authService := NewAuthService()
	new, err := authService.NewRole(role)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取角色
// @Id 20
// @Tags 角色管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "角色ID"
// @Success 200 object response.SuccessRes{data=Role} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /roles/:id [GET]
func GetRoleByID(c *gin.Context) {
	var uri RoleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	authService := NewAuthService()
	role, err := authService.GetRoleByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, role)

}

// @Summary 根据ID更新角色
// @Id 21
// @Tags 角色管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "角色ID"
// @Param role_info body RoleNew true "角色信息"
// @Success 200 object response.SuccessRes{data=Role} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /roles/:id [PUT]
func UpdateRole(c *gin.Context) {
	var uri RoleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var role RoleNew
	if err := c.ShouldBindJSON(&role); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	role.User = claims.Username
	authService := NewAuthService()
	new, err := authService.UpdateRole(uri.ID, role)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新用户
// @Id 23
// @Tags 用户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Param menu_info body UserUpdate true "用户信息"
// @Success 200 object response.SuccessRes{data=User} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /users/:id [PUT]
func UpdateUser(c *gin.Context) {
	var uri UserID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var user UserUpdate
	if err := c.ShouldBindJSON(&user); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	user.User = claims.Username
	authService := NewAuthService()
	new, err := authService.UpdateUser(uri.ID, user, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// // @Summary API列表
// // @Id 34
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数"
// // @Param name query string false "API名称"
// // @Param route query string false "API路由"
// // @Success 200 object response.ListRes{data=[]UserAPI} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis [GET]
// func GetAPIList(c *gin.Context) {
// 	var filter APIFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	count, list, err := authService.GetAPIList(filter)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 新建API
// // @Id 35
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param api_info body APINew true "API信息"
// // @Success 200 object response.SuccessRes{data=UserAPI} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis [POST]
// func NewAPI(c *gin.Context) {
// 	var api APINew
// 	if err := c.ShouldBindJSON(&api); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	api.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewAPI(api)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID获取API
// // @Id 36
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "API ID"
// // @Success 200 object response.SuccessRes{data=UserAPI} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis/:id [GET]
// func GetAPIByID(c *gin.Context) {
// 	var uri APIID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	api, err := authService.GetAPIByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, api)

// }

// // @Summary 根据ID更新API
// // @Id 37
// // @Tags API管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "API ID"
// // @Param api_info body APINew true "API信息"
// // @Success 200 object response.SuccessRes{data=UserAPI} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /apis/:id [PUT]
// func UpdateAPI(c *gin.Context) {
// 	var uri APIID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var api APINew
// 	if err := c.ShouldBindJSON(&api); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	api.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateAPI(uri.ID, api)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 菜单列表
// // @Id 38
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param page_id query int true "页码"
// // @Param page_size query int true "每页行数（5/10/15/20）"
// // @Param name query string false "菜单名称"
// // @Param only_top query bool false "只显示顶级菜单"
// // @Success 200 object response.ListRes{data=[]UserMenu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus [GET]
// func GetMenuList(c *gin.Context) {
// 	var filter MenuFilter
// 	err := c.ShouldBindQuery(&filter)
// 	if err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	count, list, err := authService.GetMenuList(filter)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
// }

// // @Summary 新建菜单
// // @Id 39
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=UserMenu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus [POST]
// func NewMenu(c *gin.Context) {
// 	var menu MenuNew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewMenu(menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据ID获取菜单
// // @Id 40
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Success 200 object response.SuccessRes{data=UserMenu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus/:id [GET]
// func GetMenuByID(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetMenuByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据ID更新菜单
// // @Id 41
// // @Tags 菜单管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=UserMenu} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menus/:id [PUT]
// func UpdateMenu(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu MenuNew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.UpdateMenu(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据角色ID获取菜单权限
// // @Id 42
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /rolemenus/:id [GET]
// func GetRoleMenu(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetRoleMenuByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据角色ID更新菜单权限
// // @Id 43
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "角色ID"
// // @Param menu_info body RoleMenu true "菜单信息"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /rolemenus/:id [POST]
// func NewRoleMenu(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu RoleMenuNew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewRoleMenu(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 根据菜单ID获取API权限
// // @Id 44
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menuapis/:id [GET]
// func GetMenuApi(c *gin.Context) {
// 	var uri MenuID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	authService := NewAuthService()
// 	menu, err := authService.GetMenuAPIByID(uri.ID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, menu)

// }

// // @Summary 根据菜单ID更新API权限
// // @Id 45
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param id path int true "菜单ID"
// // @Param menu_info body MenuNew true "菜单信息"
// // @Success 200 object response.SuccessRes{data=[]int64} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /menuapis/:id [POST]
// func NewMenuApi(c *gin.Context) {
// 	var uri RoleID
// 	if err := c.ShouldBindUri(&uri); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	var menu MenuAPINew
// 	if err := c.ShouldBindJSON(&menu); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	menu.User = claims.Username
// 	authService := NewAuthService()
// 	new, err := authService.NewMenuAPI(uri.ID, menu)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// // @Summary 获取当前用户的前端路由
// // @Id 46
// // @Tags 权限管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Success 200 object response.SuccessRes{data=interface{}} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /mymenu [GET]
// func GetMyMenu(c *gin.Context) {
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	role_id := claims.RoleID
// 	authService := NewAuthService()
// 	new, err := authService.GetMyMenu(role_id)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	res := make(map[int64]*MyMenuDetail)
// 	for i := 0; i < len(new); i++ {
// 		if new[i].ParentID == 0 {
// 			var m MyMenuDetail
// 			m.Action = new[i].Action
// 			m.Component = new[i].Component
// 			m.Name = new[i].Name
// 			m.Title = new[i].Title
// 			m.Path = new[i].Path
// 			m.IsHidden = new[i].IsHidden
// 			m.Enabled = new[i].Enabled
// 			res[new[i].ID] = &m
// 		} else {
// 			var m MyMenuDetail
// 			m.Action = new[i].Action
// 			m.Component = new[i].Component
// 			m.Name = new[i].Name
// 			m.Title = new[i].Title
// 			m.Path = new[i].Path
// 			m.IsHidden = new[i].IsHidden
// 			m.Enabled = new[i].Enabled
// 			res[new[i].ParentID].Items = append(res[new[i].ParentID].Items, m)
// 		}
// 	}
// 	response.Response(c, res)
// }
