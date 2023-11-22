package shortcut

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/shortcuts", GetShortcutList)
	g.GET("/shortcuts/:id", GetShortcutByID)
	g.PUT("/shortcuts/:id", UpdateShortcut)
	g.POST("/shortcuts", NewShortcut)
	g.DELETE("/shortcuts/:id", DeleteShortcut)
	g.GET("/shortcut_types", GetShortcutTypeList)
	g.GET("/shortcut_types/:id", GetShortcutTypeByID)
	g.PUT("/shortcut_types/:id", UpdateShortcutType)
	g.POST("/shortcut_types", NewShortcutType)
	g.DELETE("/shortcut_types/:id", DeleteShortcutType)
}

func WxRouters(g *gin.RouterGroup) {
	g.GET("/wx/shortcuts", WxGetShortcutList)
	g.GET("/wx/shortcuts/:id", WxGetShortcutByID)
	g.PUT("/wx/shortcuts/:id", WxUpdateShortcut)
	g.POST("/wx/shortcuts", WxNewShortcut)
	g.DELETE("/wx/shortcuts/:id", WxDeleteShortcut)
	g.GET("/wx/shortcut_types", WxGetShortcutTypeList)
	g.GET("/wx/shortcut_types/:id", WxGetShortcutTypeByID)
	g.PUT("/wx/shortcut_types/:id", WxUpdateShortcutType)
	g.POST("/wx/shortcut_types", WxNewShortcutType)
	g.DELETE("/wx/shortcut_types/:id", WxDeleteShortcutType)
}
