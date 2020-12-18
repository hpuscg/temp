/*
#Time      :  2020/12/17 5:22 下午
#Author    :  chuangangshen@deepglint.com
#File      :  casbin.go
#Software  :  GoLand
*/
package middleware

import (
	"github.com/gin-gonic/gin"
	"temp/GoModelTest/casbinTest/service"
)

func CasbinHandler(c *gin.Context) {
	// 获取请求的URI
	obj := c.Request.URL.RequestURI()
	// 获取请求方法
	act := c.Request.Method
	// 获取用户的角色
	sub := c.MustGet("role_id")
	e := service.Casbin()
	// 判断策略中是否存在
	success, _ := e.Enforce(sub, obj, act)
	if global.GVA_CONFIG.System.Env == "develop" || success {
		c.Next()
	} else {
		response.FailWithDetailed(gin.H{}, "权限不足", c)
		c.Abort()
		return
	}
}
