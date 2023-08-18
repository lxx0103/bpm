package upload

import (
	"bpm/core/config"
	"bpm/core/response"
	"bpm/service"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary 文件上传列表
// @Id O001
// @Tags 文件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param organization_id query int true "组织ID"
// @Param name query string false "创建人"
// @Success 200 object response.ListRes{data=[]Upload} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /uploads [GET]
func GetUploadList(c *gin.Context) {
	var filter UploadFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	uploadService := NewUploadService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := uploadService.GetUploadList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 上传文件
// @Id O002
// @Tags 文件管理
// @version 1.0
// @Accept multipart/form-data
// @Produce application/json
// @Param  file formData file true  "上传文件"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /uploads [POST]
func NewUpload(c *gin.Context) {
	uploaded, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	dest := config.ReadConfig("file.path")
	extension := filepath.Ext(uploaded.Filename)
	newName := uuid.NewString() + extension
	path := dest + newName
	err = c.SaveUploadedFile(uploaded, path)
	if err != nil {
		response.ResponseError(c, "保存文件错误", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	uploadService := NewUploadService()
	err = uploadService.NewUpload(newName, claims.OrganizationID, claims.Username)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, newName)
}

// @Summary WX文件上传列表
// @Id O003
// @Tags 文件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param organization_id query int true "组织ID"
// @Param name query string false "创建人"
// @Success 200 object response.ListRes{data=[]Upload} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/uploads [GET]
func WxGetUploadList(c *gin.Context) {
	GetUploadList(c)
}

// @Summary WX上传文件
// @Id O004
// @Tags 文件管理
// @version 1.0
// @Accept multipart/form-data
// @Produce application/json
// @Param  file formData file true  "上传文件"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/uploads [POST]
func WxNewUpload(c *gin.Context) {
	NewUpload(c)
}

// @Summary 获取临时密钥
// @Id O005
// @Tags 文件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param app_id query string false "APPID"
// @Success 200 object response.SuccessRes{data=KeyRes} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /key [GET]
func GetUploadKey(c *gin.Context) {
	var filter KeyFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	uploadService := NewUploadService()
	key, err := uploadService.GetUploadKey(filter)
	if err != nil {
		response.ResponseError(c, "ServiceError", err)
		return
	}
	response.Response(c, key)
}

// @Summary 获取临时密钥
// @Id O006
// @Tags 文件管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param app_id query string false "APPID"
// @Success 200 object response.SuccessRes{data=KeyRes} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /wx/key [GET]
func WxGetUploadKey(c *gin.Context) {
	GetUploadKey(c)
}
