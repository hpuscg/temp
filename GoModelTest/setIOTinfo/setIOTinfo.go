/*
#Time      :  2019/4/12 下午3:39 
#Author    :  chuangangshen@deepglint.com
#File      :  setIOTinfo.go
#Software  :  GoLand
*/
package main

import (
	"net/http"
	"strings"
	"encoding/json"
	"fmt"
	"flag"
)

func main() {
	var Ip string
	flag.StringVar(&Ip, "ip", "", "sensor ip")
	flag.Parse()
	SetIotInfo(Ip)
}

type IotInfo struct {
	Server string `json:"server"`
	Topic string `json:"topic"`
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}

func SetIotInfo(ip string)  {
	fmt.Println(ip)
	iotInfo :=  IotInfo{
		Server:"10.147.8.55:1883",
		Topic:"topic/haomut",
		UserName:"haomut",
		PassWord:"haomut!@#$",
	}
	url := "http://" + ip + ":8181/api/bumble/iotInfo"
	contentType := "application/json"
	ret, err := json.Marshal(iotInfo)
	if err != nil {
		fmt.Println(err)
	}
	data := strings.NewReader(string(ret))
	resp, err := http.Post(url, contentType, data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}

