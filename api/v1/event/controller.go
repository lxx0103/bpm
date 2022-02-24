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

// @Summary 新建事件
// @Id 10
// @Tags 事件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param event_info body EventNew true "事件信息"
// @Success 200 object response.SuccessRes{data=Event} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /events [POST]
func NewEvent(c *gin.Context) {
	var event EventNew
	if err := c.ShouldBindJSON(&event); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	event.User = claims.Username
	organizationID := claims.OrganizationID
	eventService := NewEventService()
	new, err := eventService.NewEvent(event, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

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
// @Param event_info body EventNew true "事件信息"
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
