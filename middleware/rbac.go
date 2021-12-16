package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"bpm/core/response"
	"bpm/service"
)

func RbacCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*service.CustomClaims)
		role_id := claims.RoleID
		path := c.FullPath()
		method := c.Request.Method
		fmt.Println(role_id)
		fmt.Println(path)
		fmt.Println(method)
		checked := service.NewRbacService().CheckPrivilege(role_id, path, method)
		if !checked {
			response.ResponseUnauthorized(c, "AuthError", errors.New("NO PRIVILEGE"))
			return
		}
		c.Next()
	}
}
