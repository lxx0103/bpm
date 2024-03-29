package component

import (
	"bpm/core/response"

	"github.com/gin-gonic/gin"
)

// @Summary 组件列表
// @Id D001
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

// @Summary 根据ID获取组件
// @Id D002
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

// @Summary 组件列表
// @Id D003
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param event_id query int true "事件ID"
// @Param name query string false "组件编码"
// @Success 200 object response.ListRes{data=[]Component} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/components [GET]
func WxGetComponentList(c *gin.Context) {
	GetComponentList(c)
}
