package meeting

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 会议列表
// @Id H001
// @Tags 会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "会议名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]MeetingResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /meetings [GET]
func GetMeetingList(c *gin.Context) {
	var filter MeetingFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	meetingService := NewMeetingService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := meetingService.GetMeetingList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建会议
// @Id H002
// @Tags 会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param meeting_info body MeetingNew true "会议信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /meetings [POST]
func NewMeeting(c *gin.Context) {
	var meeting MeetingNew
	if err := c.ShouldBindJSON(&meeting); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	meeting.User = claims.Username
	meeting.UserID = claims.UserID
	organizationID := claims.OrganizationID
	meetingService := NewMeetingService()
	err := meetingService.NewMeeting(meeting, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取会议
// @Id H003
// @Tags 会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Success 200 object response.SuccessRes{data=MeetingResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /meetings/:id [GET]
func GetMeetingByID(c *gin.Context) {
	var uri MeetingID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	meetingService := NewMeetingService()
	meeting, err := meetingService.GetMeetingByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, meeting)

}

// @Summary 根据ID更新会议
// @Id H004
// @Tags 会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Param meeting_info body MeetingNew true "会议信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /meetings/:id [PUT]
func UpdateMeeting(c *gin.Context) {
	var uri MeetingID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var meeting MeetingNew
	if err := c.ShouldBindJSON(&meeting); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	meeting.User = claims.Username
	meeting.UserID = claims.UserID
	organizationID := claims.OrganizationID
	meetingService := NewMeetingService()
	err := meetingService.UpdateMeeting(uri.ID, meeting, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除会议
// @Id H005
// @Tags 会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /meetings/:id [DELETE]
func DeleteMeeting(c *gin.Context) {
	var uri MeetingID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	meetingService := NewMeetingService()
	err := meetingService.DeleteMeeting(uri.ID, claims.OrganizationID, claims.Username, claims.UserID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 会议列表
// @Id H006
// @Tags 小程序接口-会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "会议名称"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]MeetingResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/meetings [GET]
func WxGetMeetingList(c *gin.Context) {
	GetMeetingList(c)
}

// @Summary 新建会议
// @Id H007
// @Tags 小程序接口-会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param meeting_info body MeetingNew true "会议信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/meetings [POST]
func WxNewMeeting(c *gin.Context) {
	NewMeeting(c)
}

// @Summary 根据ID获取会议
// @Id H008
// @Tags 小程序接口-会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Success 200 object response.SuccessRes{data=MeetingResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/meetings/:id [GET]
func WxGetMeetingByID(c *gin.Context) {
	GetMeetingByID(c)

}

// @Summary 根据ID更新会议
// @Id H009
// @Tags 小程序接口-会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Param meeting_info body MeetingNew true "会议信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/meetings/:id [PUT]
func WxUpdateMeeting(c *gin.Context) {
	UpdateMeeting(c)
}

// @Summary 根据ID删除会议
// @Id H010
// @Tags 小程序接口-会议管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "会议ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/meetings/:id [DELETE]
func WxDeleteMeeting(c *gin.Context) {
	DeleteMeeting(c)
}
