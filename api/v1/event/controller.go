package event

import (
	"bpm/core/response"
	"bpm/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 事件列表
// @Id F001
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param project_id query int64 false "项目ID"
// @Param name query string false "事件编码"
// @Success 200 object response.ListRes{data=[]Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events [GET]
func GetEventList(c *gin.Context) {
	var filter EventFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := eventService.GetEventList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取事件
// @Id F002
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events/:id [GET]
func GetEventByID(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	eventService := NewEventService()
	event, err := eventService.GetEventByID(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, event)

}

// @Summary 根据ID更新事件
// @Id F003
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param event_info body EventUpdate true "事件信息"
// @Success 200 object response.SuccessRes{data=Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events/:id [PUT]
func UpdateEvent(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var event EventUpdate
	if err := c.ShouldBindJSON(&event); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	event.User = claims.Username
	organizationID := claims.OrganizationID
	eventService := NewEventService()
	new, err := eventService.UpdateEvent(uri.ID, event, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 获取我的当前任务
// @Id F004
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Success 200 object response.SuccessRes{data=[]MyEvent} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/myevents [GET]
func WxGetMyEvents(c *gin.Context) {
	var filter AssignedEventFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	list, err := eventService.GetAssignedEvent(filter, claims.UserID, claims.PositionID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 获取项目中的任务
// @Id F005
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Param project_id query int64 true "项目id"
// @Success 200 object response.SuccessRes{data=[]MyEvent} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events [GET]
func WxGetEvents(c *gin.Context) {
	var filter MyEventFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	list, err := eventService.GetProjectEvent(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 保存事件
// @Id F006
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param info body SaveEventInfo true "组件内容"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/saveevents/:id [PUT]
func WxSaveEvent(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info SaveEventInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.PositionID = claims.PositionID
	err := eventService.SaveEvent(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取事件
// @Id F007
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id [GET]
func WxGetEventByID(c *gin.Context) {
	GetEventByID(c)
}

// @Summary 根据ID更新事件
// @Id F008
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param event_info body EventUpdate true "事件信息"
// @Success 200 object response.SuccessRes{data=Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id [PUT]
func WxUpdateEvent(c *gin.Context) {
	UpdateEvent(c)
}

// @Summary 审核事件
// @Id F009
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param info body AuditEventInfo true "组件内容"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/auditevents/:id [PUT]
func WxAuditEvent(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info AuditEventInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.PositionID = claims.PositionID
	err := eventService.AuditEvent(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 获取我的审核任务
// @Id F010
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Success 200 object response.SuccessRes{data=[]MyEvent} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/myaudits [GET]
func WxGetMyAudits(c *gin.Context) {
	var filter AssignedAuditFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	list, err := eventService.GetAssignedAudit(filter, claims.UserID, claims.PositionID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 事件签到
// @Id F011
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param event_info body NewCheckin true "签到信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id/checkin [POST]
func WxNewEventCheckin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var checkin NewCheckin
	if err = c.ShouldBindJSON(&checkin); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	checkin.User = claims.Username
	checkin.PositionID = claims.PositionID
	checkin.UserID = claims.UserID
	checkin.OrganizationID = claims.OrganizationID
	eventService := NewEventService()
	err = eventService.NewCheckin(id, checkin)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	// new := checkin
	response.Response(c, "ok")
}

// @Summary 事件签到列表
// @Id F012
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param project_id query int64 false "项目ID"
// @Param name query string false "用户名称"
// @Param organization_id query int64 false "组织ID"
// @Param event_id query int64 false "事件ID"
// @Param user_id query int64 false "用户ID"
// @Param from query string false "开始日期（2016-01-01）"
// @Param to query string false "结束日期（2016-01-01）"
// @Success 200 object response.ListRes{data=[]CheckinResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /checkins [GET]
func GetCheckinList(c *gin.Context) {
	var filter CheckinFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := eventService.GetCheckinList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 事件签到列表
// @Id F013
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param project_id query int64 false "项目ID"
// @Param name query string false "用户名称"
// @Param organization_id query int64 false "组织ID"
// @Param event_id query int64 false "事件ID"
// @Param user_id query int64 false "用户ID"
// @Param from query string false "开始日期（2016-01-01）"
// @Param to query string false "结束日期（2016-01-01）"
// @Success 200 object response.ListRes{data=[]CheckinResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/checkins [GET]
func WxGetCheckinList(c *gin.Context) {
	var filter CheckinFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	filter.UserID = claims.UserID
	count, list, err := eventService.GetCheckinList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 获取事件审核历史
// @Id F014
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=[]EventAuditHistoryResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events/:id/audits [GET]
func GetAuditHistory(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	historys, err := eventService.GetEventAuditHistory(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, historys)

}

// @Summary 获取事件审核历史
// @Id F015
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=[]EventAuditHistoryResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id/audits [GET]
func WxGetAuditHistory(c *gin.Context) {
	GetAuditHistory(c)
}

// @Summary 顾客反馈
// @Id F016
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param info body EventReviewNew true "组件内容"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id/reviews [POST]
func WxReviewEvent(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info EventReviewNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.PositionID = claims.PositionID
	err := eventService.ReviewEvent(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 获取顾客反馈历史
// @Id F017
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=[]EventReviewResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events/:id/reviews [GET]
func GetReview(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	reviews, err := eventService.GetEventReview(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, reviews)

}

// @Summary 获取顾客反馈历史
// @Id F018
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Success 200 object response.SuccessRes{data=[]EventReviewResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id/reviews [GET]
func WxGetReview(c *gin.Context) {
	GetReview(c)
}

// @Summary 根据ID更新事件截止日期
// @Id F019
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param event_info body EventDeadlineNew true "事件信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events/:id/deadline [PUT]
func UpdateEventDeadline(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var event EventDeadlineNew
	if err := c.ShouldBindJSON(&event); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	event.User = claims.Username
	organizationID := claims.OrganizationID
	eventService := NewEventService()
	err := eventService.UpdateEventDeadline(uri.ID, event, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID更新事件截止日期
// @Id F020
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param event_info body EventDeadlineNew true "事件信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id/deadline [PUT]
func WxUpdateEventDeadline(c *gin.Context) {
	UpdateEventDeadline(c)
}

// @Summary 处理顾客反馈
// @Id F021
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "反馈ID"
// @Param info body HandleReviewInfo true "处理内容"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/reviews/:id/handle [PUT]
func WxHandleReview(c *gin.Context) {
	var uri EventID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info HandleReviewInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Username
	info.UserID = claims.UserID
	info.PositionID = claims.PositionID
	err := eventService.HandleReview(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}
