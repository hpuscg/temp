package Controllers

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

func HandlerRequest(req *http.Request)  {
	req.ParseForm()
	result, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		fmt.Println(string(err.Error()))
		return
	}
	ret := string(result)
	fmt.Println(ret)
}
