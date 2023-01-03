package example

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 案例列表
// @Id 103
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "案例编码"
// @Param style query string false "装修风格"
// @Param type query string false "类型"
// @Param room query string false "居室"
// @Param status query int false "状态"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]ExampleResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples [GET]
func GetExampleList(c *gin.Context) {
	var filter ExampleFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	exampleService := NewExampleService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := exampleService.GetExampleList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建案例
// @Id 104
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param example_info body ExampleNew true "案例信息"
// @Success 200 object response.SuccessRes{data=Example} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples [POST]
func NewExample(c *gin.Context) {
	var example ExampleNew
	if err := c.ShouldBindJSON(&example); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	example.User = claims.Username
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	new, err := exampleService.NewExample(example, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取案例
// @Id 105
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Success 200 object response.SuccessRes{data=Example} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id [GET]
func GetExampleByID(c *gin.Context) {
	var uri ExampleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	example, err := exampleService.GetExampleByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, example)

}

// @Summary 根据ID更新案例
// @Id 106
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Param example_info body ExampleNew true "案例信息"
// @Success 200 object response.SuccessRes{data=Example} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id [PUT]
func UpdateExample(c *gin.Context) {
	var uri ExampleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var example ExampleNew
	if err := c.ShouldBindJSON(&example); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	example.User = claims.Username
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	new, err := exampleService.UpdateExample(uri.ID, example, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 案例列表
// @Id 107
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "案例编码"
// @Param style query string false "装修风格"
// @Param type query string false "类型"
// @Param room query string false "居室"
// @Param status query int false "状态"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]ExampleListResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/examples [GET]
func WxGetExampleList(c *gin.Context) {
	GetExampleList(c)
}

// @Summary 根据ID获取案例
// @Id 108
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Success 200 object response.SuccessRes{data=ExampleResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/examples/:id [GET]
func WxGetExampleByID(c *gin.Context) {
	GetExampleByID(c)
}

// @Summary 案例列表
// @Id 141
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "案例编码"
// @Param style query string false "装修风格"
// @Param type query string false "类型"
// @Param room query string false "居室"
// @Param status query int false "状态"
// @Param organization_id query int64 false "组织ID"
// @Success 200 object response.ListRes{data=[]ExampleListResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/examples [GET]
func PortalGetExampleList(c *gin.Context) {
	var filter ExampleFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	exampleService := NewExampleService()
	count, list, err := exampleService.GetExampleList(filter, 0)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取案例
// @Id 142
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Success 200 object response.SuccessRes{data=ExampleResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/examples/:id [GET]
func PortalGetExampleByID(c *gin.Context) {
	var uri ExampleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	exampleService := NewExampleService()
	example, err := exampleService.GetExampleByID(uri.ID, 0)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, example)
}
