/*
#Time      :  2019/3/27 下午1:39 
#Author    :  chuangangshen@deepglint.com
#File      :  jsonTest.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// byte2string()
	StructToJson()
}

func StructToJson() {
	data := HttpEvent{
		DeviceId:"123",
		DeviceType:123,
	}
	ret, err := json.Marshal(data)
	if err != nil {
		fmt.Println("===25", err)
	}
	var data2 Http2Event
	json.Unmarshal(ret, &data2)
	fmt.Printf("%+v\n", data2)
}


type Http2Event struct {
	DeviceId    string  // `json:"device_id"`
	DeviceType  int
	EventType   int
	EventTime   string
	EventDetail interface{}
}

type HttpEvent struct {
	DeviceId    string      `json:"deviceid"`
	DeviceType  int         `json:"deviceType"`
	EventType   int         `json:"eventType"`
	EventTime   string      `json:"eventTime"`
	EventDetail interface{} `json:"eventDetail"`
}


func byte2string() {
	httpevent := HttpEvent{
		DeviceType: 123,
		DeviceId: "1234",
	}
	data, err := json.Marshal(httpevent)
	if err != nil {
		fmt.Println(err)
	}
	var resp HttpEvent
	err = json.Unmarshal([]byte(string(data)), &resp)
	fmt.Println(err)
	fmt.Println(resp)
}