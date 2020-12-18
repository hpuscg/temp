/*
#Time      :  2020/12/4 7:26 下午
#Author    :  chuangangshen@deepglint.com
#File      :  ginTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func main() {
	// zeroTest()
	printTest()
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
