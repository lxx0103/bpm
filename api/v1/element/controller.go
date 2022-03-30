package element

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
// @Param node_id query int true "事件ID"
// @Param name query string false "组件编码"
// @Success 200 object response.ListRes{data=[]Element} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /elements [GET]
func GetElementList(c *gin.Context) {
	var filter ElementFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	elementService := NewElementService()
	count, list, err := elementService.GetElementList(filter)
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
// @Param element_info body ElementNew true "组件信息"
// @Success 200 object response.SuccessRes{data=Element} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /elements [POST]
func NewElement(c *gin.Context) {
	var element ElementNew
	if err := c.ShouldBindJSON(&element); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	element.User = claims.Username
	elementService := NewElementService()
	new, err := elementService.NewElement(element, claims.OrganizationID)
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
// @Success 200 object response.SuccessRes{data=Element} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /elements/:id [GET]
func GetElementByID(c *gin.Context) {
	var uri ElementID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	elementService := NewElementService()
	element, err := elementService.GetElementByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, element)

}

// @Summary 根据ID更新组件
// @Id 16
// @Tags 组件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组件ID"
// @Param element_info body ElementUpdate true "组件信息"
// @Success 200 object response.SuccessRes{data=Element} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /elements/:id [PUT]
func UpdateElement(c *gin.Context) {
	var uri ElementID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var element ElementUpdate
	if err := c.ShouldBindJSON(&element); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	element.User = claims.Username
	elementService := NewElementService()
	new, err := elementService.UpdateElement(uri.ID, element, claims.OrganizationID)
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
// @Router /elements/:id [DELETE]
func DelElement(c *gin.Context) {
	var uri ElementID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	elementService := NewElementService()
	err := elementService.DeleteElement(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
