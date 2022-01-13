package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	bumbleToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMTQ4MzIzODUsImlzcyI6ImRlZXBnbGludCIsIlVzZXJ" +
		"JRCI6ImJ1bWJsZSJ9.VWGZm5LkQDoyukekwg6KEG-BbAkP28lcpx8D32t5mLw"
)

type SensorInfo struct {
	TimeStamp string
	Usb       int
	Memory    int
	Service   int
	HomeDisk  int
}

func main() {
	// time.Sleep(10 * time.Second)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	getSensorInfo()
}

func getSensorInfo() {
	getNtpServerUrl := "http://127.0.0.1/api/sensorstatus"
	client := &http.Client{}
	client.Timeout = 30 * time.Second
	req, err := http.NewRequest("GET", getNtpServerUrl, nil)
	req.Close = true
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", bumbleToken)
	resp, err := client.Do(req)
	var respData struct {
		Code     int        `json:"code"`
		Msg      string     `json:"msg"`
		Redirect string     `json:"redirect"`
		Data     SensorInfo `json:"data"`
	}
	if resp == nil || err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteData, &respData)
	if err != nil {
		panic(err)
	}
	respData.Data.TimeStamp = time.Now().Local().String()
	fmt.Printf("%+v\n", respData.Data)
}
