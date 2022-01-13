package maim

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func MelodyTest() {
	m := melody.New()
	e := gin.Default()
	e.GET("/ws", func(c *gin.Context) { m.HandleRequest(c.Writer, c.Request) })
	m.HandleMessage(func(s *melody.Session, b []byte) {})
}
