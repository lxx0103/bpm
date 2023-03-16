package node

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 节点列表
// @Id J001
// @Tags 节点管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "节点名称"
// @Param template_id query int64 true "模板ID"
// @Success 200 object response.ListRes{data=[]Node} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nodes [GET]
func GetNodeList(c *gin.Context) {
	var filter NodeFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	nodeService := NewNodeService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := nodeService.GetNodeList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建节点
// @Id J002
// @Tags 节点管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param node_info body NodeNew true "节点信息"
// @Success 200 object response.SuccessRes{data=Node} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nodes [POST]
func NewNode(c *gin.Context) {
	var node NodeNew
	if err := c.ShouldBindJSON(&node); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	node.User = claims.Username
	organizationID := claims.OrganizationID
	nodeService := NewNodeService()
	new, err := nodeService.NewNode(node, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取节点
// @Id J003
// @Tags 节点管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "节点ID"
// @Success 200 object response.SuccessRes{data=Node} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nodes/:id [GET]
func GetNodeByID(c *gin.Context) {
	var uri NodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	nodeService := NewNodeService()
	node, err := nodeService.GetNodeByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, node)

}

// @Summary 根据ID更新节点
// @Id J004
// @Tags 节点管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "节点ID"
// @Param node_info body NodeUpdate true "节点信息"
// @Success 200 object response.SuccessRes{data=Node} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nodes/:id [PUT]
func UpdateNode(c *gin.Context) {
	var uri NodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var node NodeUpdate
	if err := c.ShouldBindJSON(&node); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	node.User = claims.Username
	organizationID := claims.OrganizationID
	nodeService := NewNodeService()
	new, err := nodeService.UpdateNode(uri.ID, node, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除节点
// @Id J005
// @Tags 节点管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "节点ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nodes/:id [DELETE]
func DeleteNode(c *gin.Context) {
	var uri NodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	nodeService := NewNodeService()
	err := nodeService.DeleteNode(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
