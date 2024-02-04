package team

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 班组列表
// @Id T001
// @Tags 班组管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "班组名称"
// @Param status query string false "状态"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]TeamResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /teams [GET]
func GetTeamList(c *gin.Context) {
	var filter TeamFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	teamService := NewTeamService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := teamService.GetTeamList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建班组
// @Id T002
// @Tags 班组管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param team_info body TeamNew true "班组信息"
// @Success 200 object response.SuccessRes{data=Team} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /teams [POST]
func NewTeam(c *gin.Context) {
	var team TeamNew
	if err := c.ShouldBindJSON(&team); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	team.User = claims.Username
	organizationID := claims.OrganizationID
	teamService := NewTeamService()
	new, err := teamService.NewTeam(team, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取班组
// @Id T003
// @Tags 班组管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "班组ID"
// @Success 200 object response.SuccessRes{data=Team} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /teams/:id [GET]
func GetTeamByID(c *gin.Context) {
	var uri TeamID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	teamService := NewTeamService()
	team, err := teamService.GetTeamByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, team)

}

// @Summary 根据ID更新班组
// @Id T004
// @Tags 班组管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "班组ID"
// @Param team_info body TeamNew true "班组信息"
// @Success 200 object response.SuccessRes{data=Team} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /teams/:id [PUT]
func UpdateTeam(c *gin.Context) {
	var uri TeamID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var team TeamNew
	if err := c.ShouldBindJSON(&team); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	team.User = claims.Username
	organizationID := claims.OrganizationID
	teamService := NewTeamService()
	new, err := teamService.UpdateTeam(uri.ID, team, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 班组列表
// @Id T005
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "班组编码"
// @Success 200 object response.ListRes{data=[]Team} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/teams [GET]
func WxGetTeamList(c *gin.Context) {
	GetTeamList(c)
}
