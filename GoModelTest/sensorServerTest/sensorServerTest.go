/*
#Time      :  2019/12/17 11:05 AM 
#Author    :  chuangangshen@deepglint.com
#File      :  sensorServerTest.go
#Software  :  GoLand
*/
package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
)

const (
	SensorIp = "192.168.5.251"
)

var (
	Image = map[string]bool{
		"vibo2vibo":   true,
		"libra-cuda":  true,
		"bumble-bee":  true,
		"onvifserver": true,
		"nanomsg2nsq": true,
		"eventserver": true,
		"foundation":  true,
		"tunerd":      true,
		"adu":         true,
		"vulcand":     true,
		"flowservice": true,
		"nsqd-live":   true,
		"etcd":        true,
		"crtmpserver": true,
	}
)

func main() {
	SensorServerTest(SensorIp)
}

type Container struct {
	Id      string
	Image   string
	Command string
	Created int64
	Status  string
	Names   []string
}

func SensorServerTest(ip string) {
	url := "http://" + ip + ":8008/api/container/list"
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data []Container
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	for _, ret := range data {
		name := strings.Trim(ret.Names[0], "/")
		if strings.Contains(ret.Status, "Up") {
			// fmt.Println("正在运行")
		} else if strings.Contains(ret.Status, "Exited") {
			fmt.Printf("%s : 停止运行", name)
		} else {
			fmt.Printf("%s : 其他情况 : %s", name, ret.Status)
		}
		Image[name] = false
	}
	// fmt.Printf("%+v\n", Image)
	for key, value := range Image {
		if value {
			fmt.Printf("%s : 未运行", key)
		}
	}
}
