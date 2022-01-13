package main

import (
	"embed"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"temp/GoModelTest/embedTest/simple"

	"github.com/gin-gonic/gin"
)

// 参考代码
// https://github.com/gotomicro/embedctl

func main() {
	webuiSimpleObj := &webui{
		webuiEmbed: simple.Webui,
		path:       "public",
	}
	handler := gin.Default()
	{
		// 设置ant design的路径，在config.ts里配置
		// handler.StaticFS("/ant/", http.FS(webuiAntObj))
		// 访问首页跳转到ant design的welcome页面
		/* handler.GET("/", func(ctx *gin.Context) {
			ctx.Redirect(302, "/welcome")
			return
		}) */
		// Ant Design前端访问，try file到index.html
		/* handler.GET("/welcome", func(context *gin.Context) {
			context.FileFromFS("/welcome", http.FS(webuiAntIndexObj))
		}) */

		// 设置hello world
		handler.GET("/api/hello", func(ctx *gin.Context) {
			ctx.JSON(200, "Hello EGO")
		})

		webuiIndexObj := &webuiIndex{
			webui: webuiSimpleObj,
		}

		// 设置简单的演示静态资源
		handler.StaticFS("/public", http.FS(webuiSimpleObj))

		// try file到index.html
		handler.GET("/", func(context *gin.Context) {
			context.FileFromFS("/welcome", http.FS(webuiIndexObj))
		})

	}

	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		panic(listener)
	}
	httpServer := &http.Server{
		Handler: handler,
	}
	err = httpServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

type webui struct {
	webuiEmbed embed.FS
	path       string
}

func (w *webui) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullNmae := filepath.Join(w.path, filepath.FromSlash(path.Clean("/"+name)))
	file, err := w.webuiEmbed.Open(fullNmae)
	return file, err
}

// 访问index.html
type webuiIndex struct {
	webui *webui
}

func (w *webuiIndex) Open(name string) (fs.File, error) {
	return w.webui.Open("index.html")
}
