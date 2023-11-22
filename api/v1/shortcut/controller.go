package shortcut

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 快捷模版列表
// @Id R001
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param keyword query string false "快捷模版名称"
// @Param organization_id query int64 false "组织ID"
// @Param shortcut_type query int64 false "快捷模版类型"
// @Success 200 object response.ListRes{data=[]ShortcutResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcuts [GET]
func GetShortcutList(c *gin.Context) {
	var filter ShortcutFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	shortcutService := NewShortcutService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := shortcutService.GetShortcutList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建快捷模版
// @Id R002
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param shortcut_info body ShortcutNew true "快捷模版信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcuts [POST]
func NewShortcut(c *gin.Context) {
	var shortcut ShortcutNew
	if err := c.ShouldBindJSON(&shortcut); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcut.User = claims.Username
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	err := shortcutService.NewShortcut(shortcut, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取快捷模版
// @Id R003
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Success 200 object response.SuccessRes{data=ShortcutResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcuts/:id [GET]
func GetShortcutByID(c *gin.Context) {
	var uri ShortcutID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	shortcut, err := shortcutService.GetShortcutByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, shortcut)

}

// @Summary 根据ID更新快捷模版
// @Id R004
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Param info body ShortcutUpdate true "快捷模版信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcuts/:id [PUT]
func UpdateShortcut(c *gin.Context) {
	var uri ShortcutID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var shortcut ShortcutUpdate
	if err := c.ShouldBindJSON(&shortcut); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcut.User = claims.Username
	shortcut.UserID = claims.UserID
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	err := shortcutService.UpdateShortcut(uri.ID, shortcut, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除快捷模版
// @Id R005
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcuts/:id [DELETE]
func DeleteShortcut(c *gin.Context) {
	var uri ShortcutID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcutService := NewShortcutService()
	err := shortcutService.DeleteShortcut(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 快捷模版列表
// @Id R006
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param keyword query string false "快捷模版名称"
// @Param organization_id query int64 false "组织ID"
// @Param shortcut_type query int64 false "快捷模版类型"
// @Success 200 object response.ListRes{data=[]ShortcutResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcuts [GET]
func WxGetShortcutList(c *gin.Context) {
	GetShortcutList(c)
}

// @Summary 新建快捷模版
// @Id R007
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param shortcut_info body ShortcutNew true "快捷模版信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcuts [POST]
func WxNewShortcut(c *gin.Context) {
	NewShortcut(c)
}

// @Summary 根据ID获取快捷模版
// @Id R008
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Success 200 object response.SuccessRes{data=ShortcutResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcuts/:id [GET]
func WxGetShortcutByID(c *gin.Context) {
	GetShortcutByID(c)
}

// @Summary 根据ID更新快捷模版
// @Id R009
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Param info body ShortcutUpdate true "快捷模版信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcuts/:id [PUT]
func WxUpdateShortcut(c *gin.Context) {
	UpdateShortcut(c)
}

// @Summary 根据ID删除快捷模版
// @Id R010
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcuts/:id [DELETE]
func WxDeleteShortcut(c *gin.Context) {
	DeleteShortcut(c)
}

// @Summary 快捷模版类别列表
// @Id R011
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param parent_id query int64 false "父级ID"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.SuccessRes{data=[]ShortcutTypeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcut_types [GET]
func GetShortcutTypeList(c *gin.Context) {
	var filter ShortcutTypeFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	shortcutService := NewShortcutService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	list, err := shortcutService.GetShortcutTypeList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 新建快捷模版类别
// @Id R012
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param shortcut_info body ShortcutTypeNew true "快捷模版类别信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcut_types [POST]
func NewShortcutType(c *gin.Context) {
	var shortcutType ShortcutTypeNew
	if err := c.ShouldBindJSON(&shortcutType); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcutType.User = claims.Username
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	err := shortcutService.NewShortcutType(shortcutType, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取快捷模版类别
// @Id R013
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Success 200 object response.SuccessRes{data=ShortcutTypeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcut_types/:id [GET]
func GetShortcutTypeByID(c *gin.Context) {
	var uri ShortcutID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	shortcut, err := shortcutService.GetShortcutTypeByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, shortcut)

}

// @Summary 根据ID更新快捷模版类别
// @Id R014
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Param info body ShortcutTypeUpdate true "快捷模版类别信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcut_types/:id [PUT]
func UpdateShortcutType(c *gin.Context) {
	var uri ShortcutTypeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var shortcut ShortcutTypeUpdate
	if err := c.ShouldBindJSON(&shortcut); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcut.User = claims.Username
	organizationID := claims.OrganizationID
	shortcutService := NewShortcutService()
	err := shortcutService.UpdateShortcutType(uri.ID, shortcut, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除快捷模版类别
// @Id R015
// @Tags 快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shortcut_types/:id [DELETE]
func DeleteShortcutType(c *gin.Context) {
	var uri ShortcutTypeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	shortcutService := NewShortcutService()
	err := shortcutService.DeleteShortcutType(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 快捷模版类别列表
// @Id R016
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param parent_id query int64 false "父级ID"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]ShortcutTypeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcut_types [GET]
func WxGetShortcutTypeList(c *gin.Context) {
	GetShortcutTypeList(c)
}

// @Summary 新建快捷模版类别
// @Id R017
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param shortcut_info body ShortcutTypeNew true "快捷模版类别信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcut_types [POST]
func WxNewShortcutType(c *gin.Context) {
	NewShortcutType(c)
}

// @Summary 根据ID获取快捷模版类别
// @Id R018
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Success 200 object response.SuccessRes{data=ShortcutTypeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcut_types/:id [GET]
func WxGetShortcutTypeByID(c *gin.Context) {
	GetShortcutTypeByID(c)
}

// @Summary 根据ID更新快捷模版类别
// @Id R019
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Param info body ShortcutTypeUpdate true "快捷模版类别信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcut_types/:id [PUT]
func WxUpdateShortcutType(c *gin.Context) {
	UpdateShortcutType(c)
}

// @Summary 根据ID删除快捷模版类别
// @Id R020
// @Tags 小程序接口-快捷模版管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "快捷模版类别ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/shortcut_types/:id [DELETE]
func WxDeleteShortcutType(c *gin.Context) {
	DeleteShortcutType(c)
}
