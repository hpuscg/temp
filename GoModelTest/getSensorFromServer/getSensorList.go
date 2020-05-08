/*
#Time      :  2020/5/8 9:53 上午
#Author    :  chuangangshen@deepglint.com
#File      :  getSensorList.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	serverIp string
	logDir = "logs"
)

func main() {
	flag.StringVar(&serverIp, "serverIp", "192.168.12.12", "网管服务器IP")
	flag.Parse()
	_ = os.Mkdir(logDir, os.ModePerm)
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath(logDir))
	GetSensorList()
}

// 设备状态结构
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

// 设备全部运行状态
type Sensor struct {
	Host         HostModel
	LsReportTime time.Time
	IsInControl  bool
	Status       bool
}

// 获取网管服务器上挂载的sensor信息
func GetSensorList() {
	url := "http://" + serverIp + ":8008/api/sensor_list"
	resp, err := http.Get(url)
	if err != nil {
		glog.Infoln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infoln(err)
	}
	var data []Sensor
	_ = json.Unmarshal(body, &data)
	for _, sensor := range data {
		glog.Infof("%+v", sensor)
	}
}

