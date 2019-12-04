package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"io"
)

func main() {
	/*
	// 方法二
	// router.StartTest(":5678")
	*/
	http.HandleFunc("/event", HandedRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func HandedRequest(rw http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	result, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		io.WriteString(rw, err.Error())
		return
	}
	ret := string(result)
	fmt.Println(ret)
	io.WriteString(rw, ret)
}
