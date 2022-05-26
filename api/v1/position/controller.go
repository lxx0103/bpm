package position

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 职位列表
// @Id 28
// @Tags 职位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "职位编码"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]Position} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /positions [GET]
func GetPositionList(c *gin.Context) {
	var filter PositionFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	positionService := NewPositionService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := positionService.GetPositionList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建职位
// @Id 29
// @Tags 职位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param position_info body PositionNew true "职位信息"
// @Success 200 object response.SuccessRes{data=Position} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /positions [POST]
func NewPosition(c *gin.Context) {
	var position PositionNew
	if err := c.ShouldBindJSON(&position); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	position.User = claims.Username
	organizationID := claims.OrganizationID
	positionService := NewPositionService()
	new, err := positionService.NewPosition(position, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取职位
// @Id 30
// @Tags 职位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "职位ID"
// @Success 200 object response.SuccessRes{data=Position} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /positions/:id [GET]
func GetPositionByID(c *gin.Context) {
	var uri PositionID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	positionService := NewPositionService()
	position, err := positionService.GetPositionByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, position)

}

// @Summary 根据ID更新职位
// @Id 31
// @Tags 职位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "职位ID"
// @Param position_info body PositionNew true "职位信息"
// @Success 200 object response.SuccessRes{data=Position} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /positions/:id [PUT]
func UpdatePosition(c *gin.Context) {
	var uri PositionID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var position PositionNew
	if err := c.ShouldBindJSON(&position); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	position.User = claims.Username
	organizationID := claims.OrganizationID
	positionService := NewPositionService()
	new, err := positionService.UpdatePosition(uri.ID, position, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 职位列表
// @Id 86
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "职位编码"
// @Success 200 object response.ListRes{data=[]Position} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/positions [GET]
func WxGetPositionList(c *gin.Context) {
	GetPositionList(c)
}
