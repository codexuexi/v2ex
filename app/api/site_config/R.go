package site_config

import (
	"github.com/123456/c_code"
	"github.com/gin-gonic/gin"
	"v2ex/app/api"
)

func R(r *gin.RouterGroup) {
	r1 := r.Group("/config")
	r1.Use(isok)
	r1.GET("/seo", seo)
	r1.POST("/seo", seopost)
}

func isok(c *gin.Context) {
	user := api.GetNowUserInfo(c)
	if user.MemberType != 1 {
		result := c_code.V1GinError(500, "没权限啊")
		c.JSON(200, result)
		c.Abort()
		return
	}
}