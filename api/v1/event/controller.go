package event

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 事件列表
// @Id 9
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

// // @Summary 新建事件
// // @Id 10
// // @Tags 事件管理
// // @version 1.0
// // @Accept application/json
// // @Produce application/json
// // @Param event_info body EventNew true "事件信息"
// // @Success 200 object response.SuccessRes{data=Event} 成功
// // @Failure 400 object response.ErrorRes 内部错误
// // @Router /events [POST]
// func NewEvent(c *gin.Context) {
// 	var event EventNew
// 	if err := c.ShouldBindJSON(&event); err != nil {
// 		response.ResponseError(c, "BindingError", err)
// 		return
// 	}
// 	claims := c.MustGet("claims").(*service.CustomClaims)
// 	event.User = claims.Username
// 	organizationID := claims.OrganizationID
// 	eventService := NewEventService()
// 	new, err := eventService.NewEvent(event, organizationID)
// 	if err != nil {
// 		response.ResponseError(c, "DatabaseError", err)
// 		return
// 	}
// 	response.Response(c, new)
// }

// @Summary 根据ID获取事件
// @Id 11
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
	eventService := NewEventService()
	event, err := eventService.GetEventByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, event)

}

// @Summary 根据ID更新事件
// @Id 12
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
// @Id 53
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Success 200 object response.SuccessRes{data=[]MyEvent} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events [GET]
func WxGetEvents(c *gin.Context) {
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

// @Summary 获取我创建的任务
// @Id 70
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param status query string true "显示所有all/激活active"
// @Param project_id query int64 true "项目id"
// @Success 200 object response.SuccessRes{data=[]MyEvent} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/myevents [GET]
func WxGetMyEvents(c *gin.Context) {
	var filter MyEventFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	eventService := NewEventService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	list, err := eventService.GetMyEvent(filter, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 保存事件
// @Id 71
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "事件ID"
// @Param info body SaveEventInfo true "组件内容"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/events/:id [PUT]
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
