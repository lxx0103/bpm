package vendors

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 商家列表
// @Id P001
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "商家名称"
// @Param brand query string false "品牌名称"
// @Param material query string false "材料名称"
// @Success 200 object response.ListRes{data=[]VendorsResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [GET]
func GetVendorsList(c *gin.Context) {
	var filter VendorsFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	vendorsService := NewVendorsService()
	count, list, err := vendorsService.GetVendorsList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建商家
// @Id P002
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param vendor_info body VendorsNew true "商家信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [POST]
func NewVendors(c *gin.Context) {
	var vendors VendorsNew
	if err := c.ShouldBindJSON(&vendors); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendors.User = claims.Username
	vendorsService := NewVendorsService()
	err := vendorsService.NewVendors(vendors)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取商家
// @Id P003
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=VendorsResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [GET]
func GetVendorsByID(c *gin.Context) {
	var uri VendorsID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	vendorsService := NewVendorsService()
	vendors, err := vendorsService.GetVendorsByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, vendors)

}

// @Summary 根据ID更新商家
// @Id P004
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Param vendor_info body VendorsNew true "商家信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [PUT]
func UpdateVendors(c *gin.Context) {
	var uri VendorsID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var vendors VendorsNew
	if err := c.ShouldBindJSON(&vendors); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendors.User = claims.Username
	vendorsService := NewVendorsService()
	err := vendorsService.UpdateVendors(uri.ID, vendors)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除商家
// @Id P005
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [DELETE]
func DeleteVendors(c *gin.Context) {
	var uri VendorsID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendorsService := NewVendorsService()
	err := vendorsService.DeleteVendors(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 商家列表
// @Id P006
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "商家名称"
// @Param brand query string false "品牌名称"
// @Param material query string false "材料名称"
// @Success 200 object response.ListRes{data=[]VendorsResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/vendors [GET]
func PortalGetVendorsList(c *gin.Context) {
	GetVendorsList(c)
}

// @Summary 根据ID获取商家
// @Id P007
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=VendorsResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/vendors/:id [GET]
func PortalGetVendorsByID(c *gin.Context) {
	GetVendorsByID(c)
}
