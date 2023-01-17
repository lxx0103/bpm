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
// @Param mixed query string false "搜索名称和楼盘"
// @Param priority query string false "all所有/index推荐"
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
// @Param mixed query string false "搜索名称和楼盘"
// @Param priority query string false "all所有/index推荐"
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
// @Param mixed query string false "搜索名称和楼盘"
// @Param priority query string false "all所有/index推荐"
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

// @Summary 案例材料列表
// @Id 158
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Success 200 object response.ListRes{data=[]ExampleMaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id/materials [GET]
func GetExampleMaterialList(c *gin.Context) {
	var uri ExampleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	exampleService := NewExampleService()
	list, err := exampleService.GetExampleMaterialList(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 根据ID获取案例材料
// @Id 159
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Param material_id path int true "案例材料ID"
// @Success 200 object response.SuccessRes{data=ExampleMaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id/materials/:material_id [GET]
func GetExampleMaterialByID(c *gin.Context) {
	var uri ExampleMaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	exampleService := NewExampleService()
	example, err := exampleService.GetExampleMaterialByID(uri.ID, uri.MaterialID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, example)

}

// @Summary 新建案例材料
// @Id 160
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Param material_id path int true "案例材料ID"
// @Param example_info body ExampleMaterialNew true "案例材料信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id/materials [POST]
func NewExampleMaterial(c *gin.Context) {
	var uri ExampleID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var material ExampleMaterialNew
	if err := c.ShouldBindJSON(&material); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	material.User = claims.Username
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	err := exampleService.NewExampleMaterial(material, uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 更新案例材料
// @Id 161
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Param material_id path int true "案例材料ID"
// @Param example_info body ExampleMaterialNew true "案例材料信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id/materials/material_id [PUT]
func UpdateExampleMaterial(c *gin.Context) {
	var uri ExampleMaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var material ExampleMaterialNew
	if err := c.ShouldBindJSON(&material); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	material.User = claims.Username
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	err := exampleService.UpdateExampleMaterial(material, uri.ID, uri.MaterialID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 删除案例材料
// @Id 162
// @Tags 案例管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Param material_id path int true "案例材料ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /examples/:id/materials/material_id [DELETE]
func DeleteExampleMaterial(c *gin.Context) {
	var uri ExampleMaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	exampleService := NewExampleService()
	err := exampleService.DeleteExampleMaterial(uri.ID, uri.MaterialID, organizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 案例材料列表
// @Id 163
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "案例ID"
// @Success 200 object response.ListRes{data=[]ExampleMaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/examples/:id/materials [GET]
func PortalGetExampleMaterialList(c *gin.Context) {
	GetExampleMaterialList(c)
}
