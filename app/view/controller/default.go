package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"v2ex/config"
	"v2ex/model"
)

func DefaultData(c *gin.Context) (_ht gin.H) {
	_con := config.GetConfig()
	_ht = gin.H{}
	_ht["_______API"] = "/api/manage"
	_ht["___upload_server"] = _con.Run.UploadServer
	//初始化tdk
	seo := model.SiteConfig{}.GetSeo()
	_ht["t"] = seo.T
	_ht["d"] = seo.D
	_ht["k"] = seo.K
	_ht["t_"] = seo.T_
	_ht["title_fgf"] = seo.TitleDelimiter

	_ht["avatar"] = func(s string) string {
		return _con.Run.UploadServer + s
	}
	_ht["mu"] = func(member model.Member) string {
		return fmt.Sprintf("/member/%d", member.MID)
	}
	_ht["mumid"] = func(mid model.MIDTYPE) string {
		return fmt.Sprintf("/member/%d", mid)
	}

	return
}

func defaultData(c *gin.Context) gin.H {
	return DefaultData(c)
}
