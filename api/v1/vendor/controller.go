package vendor

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 商家列表
// @Id 151
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "商家名称"
// @Param brand query string false "品牌名称"
// @Param material query string false "材料名称"
// @Success 200 object response.ListRes{data=[]VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [GET]
func GetVendorList(c *gin.Context) {
	var filter VendorFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	vendorService := NewVendorService()
	count, list, err := vendorService.GetVendorList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建商家
// @Id 152
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param vendor_info body VendorNew true "商家信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [POST]
func NewVendor(c *gin.Context) {
	var vendor VendorNew
	if err := c.ShouldBindJSON(&vendor); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendor.User = claims.Username
	vendorService := NewVendorService()
	err := vendorService.NewVendor(vendor)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取商家
// @Id 153
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [GET]
func GetVendorByID(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	vendorService := NewVendorService()
	vendor, err := vendorService.GetVendorByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, vendor)

}

// @Summary 根据ID更新商家
// @Id 154
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Param vendor_info body VendorNew true "商家信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [PUT]
func UpdateVendor(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var vendor VendorNew
	if err := c.ShouldBindJSON(&vendor); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendor.User = claims.Username
	vendorService := NewVendorService()
	err := vendorService.UpdateVendor(uri.ID, vendor)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除商家
// @Id 155
// @Tags 商家管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [DELETE]
func DeleteVendor(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	vendorService := NewVendorService()
	err := vendorService.DeleteVendor(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 商家列表
// @Id 156
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "商家名称"
// @Param brand query string false "品牌名称"
// @Param material query string false "材料名称"
// @Success 200 object response.ListRes{data=[]VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/vendors [GET]
func PortalGetVendorList(c *gin.Context) {
	GetVendorList(c)
}

// @Summary 根据ID获取商家
// @Id 157
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "商家ID"
// @Success 200 object response.SuccessRes{data=VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/vendors/:id [GET]
func PortalGetVendorByID(c *gin.Context) {
	GetVendorByID(c)
}
