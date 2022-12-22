package auth

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/signin", Signin)
	g.POST("/signup", Signup)
}

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/roles", GetRoleList)
	g.GET("/roles/:id", GetRoleByID)
	g.PUT("/roles/:id", UpdateRole)
	g.POST("/roles", NewRole)
	g.DELETE("/roles/:id", DeleteRole)

	g.PUT("/users/:id", UpdateUser)
	g.GET("/users", GetUserList)
	g.GET("/users/:id", GetUserByID)
	g.POST("/password", UpdatePassword)

	g.GET("/apis", GetAPIList)
	g.GET("/apis/:id", GetAPIByID)
	g.PUT("/apis/:id", UpdateAPI)
	g.POST("/apis", NewAPI)

	g.GET("/menus", GetMenuList)
	g.GET("/menus/:id", GetMenuByID)
	g.POST("/menus", NewMenu)
	g.PUT("/menus/:id", UpdateMenu)
	g.DELETE("/menus/:id", DeleteMenu)

	g.GET("/rolemenus/:id", GetRoleMenu)
	g.POST("/rolemenus/:id", NewRoleMenu)
	g.GET("/menuapis/:id", GetMenuApi)
	g.POST("/menuapis/:id", NewMenuApi)

	g.GET("/mymenu", GetMyMenu)

	g.GET("/wxmodules", GetWxmoduleList)
	g.GET("/wxmodules/:id", GetWxmoduleByID)
	g.POST("/wxmodules", NewWxmodule)
	g.PUT("/wxmodules/:id", UpdateWxmodule)
	g.DELETE("/wxmodules/:id", DeleteWxmodule)

	g.GET("/positionwxmodules/:id", GetPositionWxmodule)
	g.POST("/positionwxmodules/:id", NewPositionWxmodule)

}
func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/users", WxGetUserList)
	g.PUT("/wx/users/:id", WxUpdateUser)
	g.GET("/wx/mywxmodule/:id", GetMyWxmodule)
}
