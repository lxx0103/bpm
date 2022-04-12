package member

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 项目成员列表
// @Id 34
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int true "项目ID"
// @Success 200 object response.SuccessRes{data=[]MemberResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /members [GET]
func GetMemberList(c *gin.Context) {
	var filter MemberFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	memberService := NewMemberService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	members, err := memberService.GetMemberList(filter.ProjectID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, members)
}

// @Summary 新建项目成员
// @Id 35
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param member_info body MemberNew true "成员信息"
// @Success 200 object response.SuccessRes{data=[]MemberResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /members [POST]
func NewMember(c *gin.Context) {
	var member MemberNew
	if err := c.ShouldBindJSON(&member); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	member.User = claims.Username
	organizationID := claims.OrganizationID
	memberService := NewMemberService()
	members, err := memberService.NewMember(member, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, members)
}

// @Summary 项目成员列表
// @Id 82
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_id query int true "项目ID"
// @Success 200 object response.SuccessRes{data=[]MemberResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/members [GET]
func WxGetMemberList(c *gin.Context) {
	GetMemberList(c)
}

// @Summary 新建项目成员
// @Id 83
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param member_info body MemberNew true "成员信息"
// @Success 200 object response.SuccessRes{data=[]MemberResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/members [POST]
func WxNewMember(c *gin.Context) {
	NewMember(c)
}
