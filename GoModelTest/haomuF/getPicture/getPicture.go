/*
#Time      :  2020/9/23 11:09 上午
#Author    :  chuangangshen@deepglint.com
#File      :  getPicture.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func main() {
	GetFrame()
}

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA5MzM5NzksImlhdCI6MTYwMDg0NzU3OSwidXNlcm5hbWUiOiJhZG1pbiJ9.ffspD04fGGtQO7WrgmMMbqP5hQTfA69J3-g61aT-Wjk"
)

var (
	FrameId string
)

func GetFrame() {
	frameUrl := "http://192.168.100.94/api/triggerRealframe/0"
	client1 := &http.Client{}
	client1.Timeout = 30 * time.Second
	req1, err := http.NewRequest("GET", frameUrl, nil)
	req1.Close = true
	if err != nil {
		fmt.Println(err)
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Add("authorization", token)
	response1, err := client1.Do(req1)
	if err != nil || response1 == nil {
		fmt.Println(err)
	}
	defer response1.Body.Close()
	body, err := ioutil.ReadAll(response1.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	if 1 == data["result"].(float64) {
		switch vv := data["data"].(type) {
		case map[string]interface{}:
			fmt.Println(vv)
			FrameId = vv["FrameId"].(string)
		default:
			fmt.Println(reflect.TypeOf(data["data"]))
		}
		fmt.Println(FrameId)
	} else {
		fmt.Println(data["result"].(float64))
	}


}
