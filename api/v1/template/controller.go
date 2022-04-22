package template

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 模板列表
// @Id 54
// @Tags 模板管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "模板编码"
// @Param type query int false "模板类型1内部2外部"
// @Success 200 object response.ListRes{data=[]Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /templates [GET]
func GetTemplateList(c *gin.Context) {
	var filter TemplateFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	templateService := NewTemplateService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := templateService.GetTemplateList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建模板
// @Id 55
// @Tags 模板管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param template_info body TemplateNew true "模板信息"
// @Success 200 object response.SuccessRes{data=Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /templates [POST]
func NewTemplate(c *gin.Context) {
	var template TemplateNew
	if err := c.ShouldBindJSON(&template); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	template.User = claims.Username
	organizationID := claims.OrganizationID
	templateService := NewTemplateService()
	new, err := templateService.NewTemplate(template, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取模板
// @Id 56
// @Tags 模板管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "模板ID"
// @Success 200 object response.SuccessRes{data=Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /templates/:id [GET]
func GetTemplateByID(c *gin.Context) {
	var uri TemplateID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	templateService := NewTemplateService()
	template, err := templateService.GetTemplateByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, template)

}

// @Summary 根据ID更新模板
// @Id 57
// @Tags 模板管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "模板ID"
// @Param template_info body TemplateUpdate true "模板信息"
// @Success 200 object response.SuccessRes{data=Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /templates/:id [PUT]
func UpdateTemplate(c *gin.Context) {
	var uri TemplateID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var template TemplateUpdate
	if err := c.ShouldBindJSON(&template); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	template.User = claims.Username
	organizationID := claims.OrganizationID
	templateService := NewTemplateService()
	new, err := templateService.UpdateTemplate(uri.ID, template, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除模板
// @Id 58
// @Tags 模板管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "模板ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /templates/:id [DELETE]
func DeleteTemplate(c *gin.Context) {
	var uri TemplateID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	templateService := NewTemplateService()
	err := templateService.DeleteTemplate(uri.ID, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 模板列表
// @Id 79
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "模板编码"
// @Param name query int64 false "模板编码"
// @Success 200 object response.ListRes{data=[]Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/templates [GET]
func WxGetTemplateList(c *gin.Context) {
	GetTemplateList(c)
}

// @Summary 根据ID获取模板
// @Id 80
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "模板ID"
// @Success 200 object response.SuccessRes{data=Template} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/templates/:id [GET]
func WxGetTemplateByID(c *gin.Context) {
	GetTemplateByID(c)

}
