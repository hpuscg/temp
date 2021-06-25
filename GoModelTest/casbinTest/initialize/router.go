/*
#Time      :  2020/12/18 7:20 下午
#Author    :  chuangangshen@deepglint.com
#File      :  router.go
#Software  :  GoLand
*/
package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"temp/GoModelTest/casbinTest/global"
	"temp/GoModelTest/casbinTest/middleware"
)

func Routers() *gin.Engine {
	var Router = gin.Default()
	Router.StaticFS(global.Server.Local.Path, http.Dir(global.Server.Local.Path)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.Logger.Info("use middleware logger")
	// 跨域
	Router.Use(middleware.Cors)
	global.Logger.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.Logger.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用
	/*ApiGroup := Router.Group("")
	router.InitUserRouter(ApiGroup)                  // 注册用户路由
	router.InitBaseRouter(ApiGroup)                  // 注册基础功能路由 不做鉴权
	router.InitMenuRouter(ApiGroup)                  // 注册menu路由
	router.InitAuthorityRouter(ApiGroup)             // 注册角色路由
	router.InitApiRouter(ApiGroup)                   // 注册功能api路由
	router.InitFileUploadAndDownloadRouter(ApiGroup) // 文件上传下载功能路由
	router.InitSimpleUploaderRouter(ApiGroup)        // 断点续传（插件版）
	router.InitWorkflowRouter(ApiGroup)              // 工作流相关路由
	router.InitCasbinRouter(ApiGroup)                // 权限相关路由
	router.InitJwtRouter(ApiGroup)                   // jwt相关路由
	router.InitSystemRouter(ApiGroup)                // system相关路由
	router.InitCustomerRouter(ApiGroup)              // 客户路由
	router.InitAutoCodeRouter(ApiGroup)              // 创建自动化代码
	router.InitSysDictionaryDetailRouter(ApiGroup)   // 字典详情管理
	router.InitSysDictionaryRouter(ApiGroup)         // 字典管理
	router.InitSysOperationRecordRouter(ApiGroup)    // 操作记录
	router.InitEmailRouter(ApiGroup)                 // 邮件相关路由*/

	global.Logger.Info("router register success")
	return Router
}
