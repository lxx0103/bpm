package organization

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 组织列表
// @Id 1
// @Tags 组织管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "组织编码"
// @Success 200 object response.ListRes{data=[]Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /organizations [GET]
func GetOrganizationList(c *gin.Context) {
	var filter OrganizationFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	organizationService := NewOrganizationService()
	count, list, err := organizationService.GetOrganizationList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建组织
// @Id 2
// @Tags 组织管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param organization_info body OrganizationNew true "组织信息"
// @Success 200 object response.SuccessRes{data=Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /organizations [POST]
func NewOrganization(c *gin.Context) {
	var organization OrganizationNew
	if err := c.ShouldBindJSON(&organization); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organization.User = claims.Username
	organizationService := NewOrganizationService()
	new, err := organizationService.NewOrganization(organization)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取组织
// @Id 3
// @Tags 组织管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组织ID"
// @Success 200 object response.SuccessRes{data=Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /organizations/:id [GET]
func GetOrganizationByID(c *gin.Context) {
	var uri OrganizationID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	organizationService := NewOrganizationService()
	organization, err := organizationService.GetOrganizationByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, organization)

}

// @Summary 根据ID更新组织
// @Id 4
// @Tags 组织管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组织ID"
// @Param organization_info body OrganizationNew true "组织信息"
// @Success 200 object response.SuccessRes{data=Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /organizations/:id [PUT]
func UpdateOrganization(c *gin.Context) {
	var uri OrganizationID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var organization OrganizationNew
	if err := c.ShouldBindJSON(&organization); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organization.User = claims.Username
	organizationService := NewOrganizationService()
	new, err := organizationService.UpdateOrganization(uri.ID, organization)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 获取小程序码
// @Id 99
// @Tags 组织管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param info body QrcodeFilter true "页面路径参数"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /qrcode [POST]
func GetQrCode(c *gin.Context) {
	var filter QrcodeFilter
	err := c.ShouldBindJSON(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	organizationService := NewOrganizationService()
	res, err := organizationService.GetQrCodeByPath(filter.Path)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}

// @Summary 根据ID获取组织
// @Id 109
// @Tags 小程序接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组织ID"
// @Success 200 object response.SuccessRes{data=Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/organizations/:id [GET]
func WxGetOrganizationByID(c *gin.Context) {
	GetOrganizationByID(c)
}

// @Summary 组织列表
// @Id 140
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "组织名称"
// @Param city query string false "区域"
// @Param type query string false "类型1/2"
// @Success 200 object response.ListRes{data=[]OrganizationExampleResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/organizations [GET]
func PortalGetOrganizationList(c *gin.Context) {
	var filter OrganizationFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	organizationService := NewOrganizationService()
	count, list, err := organizationService.GetPortalOrganizationList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 根据ID获取组织
// @Id 143
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "组织ID"
// @Success 200 object response.SuccessRes{data=Organization} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/organizations/:id [GET]
func PortalGetOrganizationByID(c *gin.Context) {
	GetOrganizationByID(c)
}
