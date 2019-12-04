package router

import (
	"github.com/go-martini/martini"
	"temp/GoModelTest/httpServerTest/Controllers"
	"net/http"
)

var libraM *martini.ClassicMartini

func StartTest(addr string)  {
	libraM = martini.Classic()
	libraM.Post("/event", Controllers.HandlerRequest)
	http.ListenAndServe(addr, libraM)
}
