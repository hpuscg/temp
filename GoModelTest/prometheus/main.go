package main

import (
	"net/http"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := gin.Default()
	e.Use(ginprom.PromMiddleware(nil))

	e.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "home"})
	})
	e.Run("0.0.0.0:8080")
}
