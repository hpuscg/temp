/*
#Time      :  2019/12/4 9:32 AM 
#Author    :  chuangangshen@deepglint.com
#File      :  getIpFromServer.go
#Software  :  GoLand
*/
package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"os"
)

func main() {
	GetIps()
}

type HostModel struct {
	Id            string `json:"id"`
	HostIp        string `json:"hostip"`
	Mac           string `json:"mac"`
	SensorId      string `json:"sensorid"`
	SN            string `json:"sn"`
	SubTopic      string `json:"subtopic"`
	Camera        string `json:"camera"`
	DevModel      string `json:"devmodel"`
	Chip          string `json:"model"`
	Configured    bool   `json:"configured"`
	Desc          string `json:"desc"` // used for etcd tree
	Version       string `json:"version"`
	DownloadState string `json:"downloadstate"`
	UpgradeState  string `json:"upgradestate"`
}

type Sensor struct {
	Host         HostModel
	LsReportTime time.Time
	IsInControl  bool
	Status       bool
}

func GetIps() {
	file, err := os.Open("ip.txt")
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create("ip.txt")
	}
	url := "http://192.168.100.235:8008/api/sensor_list"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(body))
	}
	var data []Sensor
	json.Unmarshal(body, &data)
	for _, value := range data {
		fmt.Printf("====%s=====\n", value.Host.HostIp)
		file.WriteString(value.Host.HostIp)
		file.WriteString("\n")
	}
}
