package client

import (
	"bpm/core/response"
	"bpm/service"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// @Summary 客户列表
// @Id B001
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "客户名称"
// @Param organization_id query int64 false "客户名称"
// @Success 200 object response.ListRes{data=[]Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /clients [GET]
func GetClientList(c *gin.Context) {
	var filter ClientFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	clientService := NewClientService()
	count, list, err := clientService.GetClientList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建客户
// @Id B002
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param client_info body ClientNew true "客户信息"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /clients [POST]
func NewClient(c *gin.Context) {
	var client ClientNew
	if err := c.ShouldBindJSON(&client); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	fmt.Println(claims)
	client.User = claims.Username
	organizationID := claims.OrganizationID
	if organizationID == 0 {
		msg := "此用户没有组织"
		response.ResponseError(c, "DatabaseError", errors.New(msg))
		return
	}
	clientService := NewClientService()
	new, err := clientService.NewClient(client, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取客户
// @Id B003
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /clients/:id [GET]
func GetClientByID(c *gin.Context) {
	var uri ClientID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	clientService := NewClientService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	client, err := clientService.GetClientByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, client)

}

// @Summary 根据ID更新客户
// @Id B004
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Param client_info body ClientNew true "客户信息"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /clients/:id [PUT]
func UpdateClient(c *gin.Context) {
	var uri ClientID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var client ClientNew
	if err := c.ShouldBindJSON(&client); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	client.User = claims.Username
	organizationID := claims.OrganizationID
	clientService := NewClientService()
	new, err := clientService.UpdateClient(uri.ID, client, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 客户列表
// @Id B005
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "客户名称"
// @Success 200 object response.ListRes{data=[]Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/clients [GET]
func WxGetClientList(c *gin.Context) {
	GetClientList(c)
}

// @Summary 新建客户
// @Id B006
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param client_info body ClientNew true "客户信息"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/clients [POST]
func WxNewClient(c *gin.Context) {
	NewClient(c)
}

// @Summary 根据ID获取客户
// @Id B007
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/clients/:id [GET]
func WxGetClientByID(c *gin.Context) {
	GetClientByID(c)

}

// @Summary 根据ID更新客户
// @Id B008
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Param client_info body ClientNew true "客户信息"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/clients/:id [PUT]
func WxUpdateClient(c *gin.Context) {
	UpdateClient(c)
}

// @Summary 根据UserID获取客户
// @Id B009
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/clients/user/:id [GET]
func WxGetClientByUserID(c *gin.Context) {
	GetClientByUserID(c)

}

// @Summary 根据ID获取客户
// @Id B010
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Success 200 object response.SuccessRes{data=Client} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /clients/user/:id [GET]
func GetClientByUserID(c *gin.Context) {
	var uri ClientID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	clientService := NewClientService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	client, err := clientService.GetClientByUserID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, client)

}
