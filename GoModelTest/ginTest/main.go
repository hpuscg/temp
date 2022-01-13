/*
#Time      :  2020/12/4 7:26 下午
#Author    :  chuangangshen@deepglint.com
#File      :  ginTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// zeroTest()
	// printTest()
	GinTest()
}

func printTest() {
	a := "123"
	fmt.Printf("%+v\n", a)
	fmt.Printf("%v\n", a)
}

func zeroTest() {
	var (
		a int         = 0
		b int64       = 0
		c interface{} = int(0)
		d interface{} = int64(0)
	)

	println(c == 0)
	println(c == a)
	println(c == b)
	println(d == b)
	println(d == 0)
}

func GinTest() {
	gin.SetMode(gin.DebugMode)
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	e.LoadHTMLGlob("/Users/hpu_scg/gocode/src/temp/GoModelTest/ginTest/dist/index.html")
	e.Static("static", "/Users/hpu_scg/gocode/src/temp/GoModelTest/ginTest/dist/static")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	e.POST("/file", UploadFile)
	e.NoRoute()
	server := &http.Server{
		Addr: net.JoinHostPort("127.0.0.1",
			"9000"),
		Handler: e,
	}
	fmt.Printf("app run on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("app run error: %s", err)
	}
}

func UploadFile(c *gin.Context) {
	fmt.Println(c.PostForm("eventId"))
	fmt.Println(c.PostForm("ip"))
	file, header, err := c.Request.FormFile("uploadFile")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", header)
	out, err := os.Create("/Users/hpu_scg/gocode/src/temp/GoModelTest/ginTest/" + header.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
	}
}
