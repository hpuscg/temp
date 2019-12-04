/*
#Time      :  2019/9/27 下午5:23 
#Author    :  chuangangshen@deepglint.com
#File      :  httpPostTest.go
#Software  :  GoLand
*/
package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
)

func main() {
	PostIotInfo()
}

type FoundationIotInfo struct {
	Iotserver string `json:"iotserver"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Topic     string `json:"topic"`
}


func PostIotInfo()  {
	var foundationIonInfo FoundationIotInfo
	foundationIonInfo.Iotserver = "tcp://"
	jsonData, _ := json.Marshal(foundationIonInfo)
	resp, _ := http.NewRequest("POST", "http://192.168.5.250:8180/api/iotInfo", strings.NewReader(string(jsonData)))
	// resp, _ := http.Post("http://192.168.5.250:8180/api/iotInfo", "", strings.NewReader(string(jsonData)))
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}

