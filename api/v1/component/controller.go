package component

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 组件列表
// @Id 13
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param event_id query int true "事件ID"
// @Param name query string false "组件编码"
// @Success 200 object response.ListRes{data=[]Component} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /components [GET]
func GetComponentList(c *gin.Context) {
	var filter ComponentFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	componentService := NewComponentService()
	count, list, err := componentService.GetComponentList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建组件
// @Id 14
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param component_info body ComponentNew true "组件信息"
// @Success 200 object response.SuccessRes{data=Component} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /components [POST]
func NewComponent(c *gin.Context) {
	var component ComponentNew
	if err := c.ShouldBindJSON(&component); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	component.User = claims.Username
	componentService := NewComponentService()
	new, err := componentService.NewComponent(component)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取组件
// @Id 15
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组件ID"
// @Success 200 object response.SuccessRes{data=Component} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /components/:id [GET]
func GetComponentByID(c *gin.Context) {
	var uri ComponentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	componentService := NewComponentService()
	component, err := componentService.GetComponentByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, component)

}

// @Summary 根据ID更新组件
// @Id 16
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组件ID"
// @Param component_info body ComponentUpdate true "组件信息"
// @Success 200 object response.SuccessRes{data=Component} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /components/:id [PUT]
func UpdateComponent(c *gin.Context) {
	var uri ComponentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var component ComponentUpdate
	if err := c.ShouldBindJSON(&component); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	component.User = claims.Username
	componentService := NewComponentService()
	new, err := componentService.UpdateComponent(uri.ID, component)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新组件
// @Id 49
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组件ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /components/:id [DELETE]
func DelComponent(c *gin.Context) {
	var uri ComponentID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	componentService := NewComponentService()
	err := componentService.DeleteComponent(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
