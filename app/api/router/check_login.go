package router

import (
	"github.com/123456/c_code"
	"github.com/123456/c_code/mc"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"v2ex/model"
)

func checkLogin(c *gin.Context) {
	_, exists := c.Get("mid")
	if !exists {
		token := c.GetHeader("____token")
		//根据 Token 查询本人是否登录
		member_token := model.MemberToken{}
		err := mc.Table(member_token.Table()).Where(bson.M{"token": token}).FindOne(&member_token)
		if err != nil {
			result_json := c_code.V1GinError(500, err.Error())
			c.JSON(200, result_json)
			c.Abort()
			return
		}
		if member_token.MID == 0 {
			result_json := c_code.V1GinError(500, "未登录")
			result_json["token"] = token
			c.JSON(200, result_json)
			c.Abort()
			return
		}
		c.Set("mid", member_token.MID)
	}

}
