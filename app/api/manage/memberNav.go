package manage

import (
	"fmt"
	"github.com/123456/c_code"
	"github.com/gin-gonic/gin"
	"v2ex/app/api"

	"v2ex/model"
)

func MemberNav(c *gin.Context) {
	user_info := api.GetNowUserInfo(c)
	u2 := model.UrlViewMemberConfig
	list := []gin.H{
		{
			"t": "我的主页",
			"u": fmt.Sprintf("/member/%d", user_info.MID),
		},
		{
			"t": "个人资料",
			"u": u2 + "/user_info",
		},
		{
			"t": "发文章",
			"u": u2 + "/send_article",
		},

		{
			"t": "发提问",
			"u": u2 + "/send_question",
		},
	}
	root_list := []gin.H{}
	if user_info.MemberType == model.MemberTypeRoot {
		u := model.UrlViewMemberManage
		root_list = append(root_list, []gin.H{
			{
				"t": "SEO设置",
				"u": u + "/seo",
			},
			{
				"t": "网站权限",
				"u": u + "/api_auth",
			},
			{
				"t": "索引查看",
				"u": u + "/db_index",
			},
			{
				"t": "数据管理",
				"u": u + "/data_manage",
			},
			{
				"t": "数据审核",
				"u": u2 + "/data_check",
			},
		}...)

	} else {
		list = append(list, gin.H{
			"t": "审核中的内容",
			"u": u2 + "/data_check",
		})
	}

	h := gin.H{
		"nav": []gin.H{
			{
				"t": "我的文章",
				"u": model.UrlViewMemberConfig + "/article",
			},
			{
				"t": "我的提问",
				"u": model.UrlViewMemberConfig + "/question",
			},
		},
		"xxx":  list,
		"root": root_list,
	}

	//登录成功
	result_json := c_code.V1GinSuccess(h)
	c.JSON(200, result_json)
	return
}
