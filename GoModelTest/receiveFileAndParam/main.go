package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	e := gin.New()
	e.POST("/file", UploadFile)
	server := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", "8080"),
		Handler: e,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func UploadFile(c *gin.Context) {
	info := c.PostForm("test")
	fmt.Println(info)
	_, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		fmt.Println(err.Error())
		c.Error(err)
		return
	}
	// fmt.Println(file.Close().Error())
	fmt.Println(header.Filename)
	c.JSON(http.StatusOK, nil)
}
