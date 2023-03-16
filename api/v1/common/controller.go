package common

import (
	"bpm/core/response"
	"bpm/service"

	"github.com/gin-gonic/gin"
)

// @Summary 品牌列表
// @Id C001
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "品牌名称"
// @Success 200 object response.ListRes{data=[]BrandResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands [GET]
func GetBrandList(c *gin.Context) {
	var filter BrandFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	count, list, err := commonService.GetBrandList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建品牌
// @Id C002
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param brand_info body BrandNew true "品牌信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands [POST]
func NewBrand(c *gin.Context) {
	var brand BrandNew
	if err := c.ShouldBindJSON(&brand); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	brand.User = claims.Username
	commonService := NewCommonService()
	err := commonService.NewBrand(brand)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取品牌
// @Id C003
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Success 200 object response.SuccessRes{data=BrandResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [GET]
func GetBrandByID(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	common, err := commonService.GetBrandByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, common)

}

// @Summary 根据ID更新品牌
// @Id C004
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Param brand_info body BrandNew true "品牌信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [PUT]
func UpdateBrand(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var brand BrandNew
	if err := c.ShouldBindJSON(&brand); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	brand.User = claims.Username
	commonService := NewCommonService()
	err := commonService.UpdateBrand(uri.ID, brand)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除品牌
// @Id C005
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [DELETE]
func DeleteBrand(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	commonService := NewCommonService()
	err := commonService.DeleteBrand(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 材料列表
// @Id C006
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "材料名称"
// @Success 200 object response.ListRes{data=[]MaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /materials [GET]
func GetMaterialList(c *gin.Context) {
	var filter MaterialFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	count, list, err := commonService.GetMaterialList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建材料
// @Id C007
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param material_info body MaterialNew true "材料信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /materials [POST]
func NewMaterial(c *gin.Context) {
	var material MaterialNew
	if err := c.ShouldBindJSON(&material); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	material.User = claims.Username
	commonService := NewCommonService()
	err := commonService.NewMaterial(material)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取材料
// @Id C008
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料ID"
// @Success 200 object response.SuccessRes{data=MaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /materials/:id [GET]
func GetMaterialByID(c *gin.Context) {
	var uri MaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	common, err := commonService.GetMaterialByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, common)

}

// @Summary 根据ID更新材料
// @Id C009
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料ID"
// @Param material_info body MaterialNew true "材料信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /materials/:id [PUT]
func UpdateMaterial(c *gin.Context) {
	var uri MaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var material MaterialNew
	if err := c.ShouldBindJSON(&material); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	material.User = claims.Username
	commonService := NewCommonService()
	err := commonService.UpdateMaterial(uri.ID, material)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除材料
// @Id C010
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /materials/:id [DELETE]
func DeleteMaterial(c *gin.Context) {
	var uri MaterialID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	commonService := NewCommonService()
	err := commonService.DeleteMaterial(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 材料列表
// @Id C011
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "材料名称"
// @Success 200 object response.ListRes{data=[]MaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/materials [GET]
func PortalGetMaterialList(c *gin.Context) {
	GetMaterialList(c)
}

// @Summary 根据ID获取材料
// @Id C012
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "材料ID"
// @Success 200 object response.SuccessRes{data=MaterialResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/materials/:id [GET]
func PortalGetMaterialByID(c *gin.Context) {
	GetMaterialByID(c)

}

// @Summary banner列表
// @Id C013
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param type query string false "all所有/index首页"
// @Success 200 object response.ListRes{data=[]BannerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /banners [GET]
func GetBannerList(c *gin.Context) {
	var filter BannerFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	count, list, err := commonService.GetBannerList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建banner
// @Id C014
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param banner_info body BannerNew true "banner信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /banners [POST]
func NewBanner(c *gin.Context) {
	var banner BannerNew
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	banner.User = claims.Username
	commonService := NewCommonService()
	err := commonService.NewBanner(banner)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID获取banner
// @Id C015
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "bannerID"
// @Success 200 object response.SuccessRes{data=BannerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /banners/:id [GET]
func GetBannerByID(c *gin.Context) {
	var uri BannerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	commonService := NewCommonService()
	common, err := commonService.GetBannerByID(uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, common)

}

// @Summary 根据ID更新banner
// @Id C016
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "bannerID"
// @Param banner_info body BannerNew true "banner信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /banners/:id [PUT]
func UpdateBanner(c *gin.Context) {
	var uri BannerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var banner BannerNew
	if err := c.ShouldBindJSON(&banner); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	banner.User = claims.Username
	commonService := NewCommonService()
	err := commonService.UpdateBanner(uri.ID, banner)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除banner
// @Id C017
// @Tags 基础信息管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "bannerID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /banners/:id [DELETE]
func DeleteBanner(c *gin.Context) {
	var uri BannerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	commonService := NewCommonService()
	err := commonService.DeleteBanner(uri.ID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary banner列表
// @Id C018
// @Tags 门户接口
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param type query string false "all所有/index首页"
// @Success 200 object response.ListRes{data=[]BannerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /portal/banners [GET]
func PortalGetBannerList(c *gin.Context) {
	GetBannerList(c)
}
