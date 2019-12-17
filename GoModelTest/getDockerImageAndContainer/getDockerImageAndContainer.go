/*
#Time      :  2019/10/24 下午4:53 
#Author    :  chuangangshen@deepglint.com
#File      :  getDockerImageAndContainer.go
#Software  :  GoLand
*/
package main

import (
	"flag"
	"net/http"
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"strings"
)

var (
	SensorIp   string
	ServerIp   string
	SensorId   string
	GetIp      bool
	GetDocker  bool
	IpListFile string
	ServerTest bool
	Image      = map[string]bool{
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
	flag.StringVar(&SensorIp, "sensorIp", "192.168.5.250", "设备IP")
	flag.StringVar(&ServerIp, "serverIp", "192.168.100.235", "网管服务器IP")
	flag.BoolVar(&GetIp, "getIp", false, "是否从网管服务器获取其下挂载的设备IP")
	flag.BoolVar(&GetDocker, "getDocker", false, "是否获取设备端的docker image 和container")
	flag.BoolVar(&ServerTest, "serverTest", false, "是否判断设备端的服务是否正常")
	flag.StringVar(&IpListFile, "ipListFile", "./ip.txt", "存储设备IP的文件")
	flag.Parse()
	if GetIp {
		GetIpFromServer()
	}
	if GetDocker {
		GetImageAndContainer()
	}
	if ServerTest {
		SensorServerTest()
	}
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

type Container struct {
	Id      string
	Image   string
	Command string
	Created int64
	Status  string
	Names   []string
}

// 获取网管服务器下挂载的设备IP
func GetIpFromServer() {
	file, err := os.Open(IpListFile)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(IpListFile)
	}
	url := "http://" + ServerIp + ":8008/api/sensor_list"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data []Sensor
	json.Unmarshal(body, &data)
	for _, value := range data {
		file.WriteString(value.Host.HostIp)
		file.WriteString("\n")
	}
}

// 获取设备的image和container
func GetImageAndContainer() {
	err := getSensorId()
	if err != nil {
		fmt.Println("get sensor id err: ", SensorIp)
		return
	}
	createFile()
	getDockerImages()
	getDockerContainer()
	getVersion()
}

// 获取设备的sensorID
func getSensorId() (err error) {
	url := "http://" + SensorIp + ":8008/api/sensorid"
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}
	SensorId = string(body)
	return
}

// 创建文件
func createFile() {
	f, err := os.Create(SensorId + ".txt")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// 获取设备上的docker images
func getDockerImages() {
	url := "http://" + SensorIp + ":4243/v1.18/images/json"
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}
	f, err := os.OpenFile(SensorId+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 066)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString(SensorIp)
	f.WriteString("\n")
	f.Write(body)
	f.WriteString("\n")
}

// 获取设备上正在运行docker container
func getDockerContainer() {
	url := "http://" + SensorIp + ":8008/api/container/list"
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}
	f, err := os.OpenFile(SensorId+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 066)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.Write(body)
	f.WriteString("\n")
}

// 获取设备的版本号
func getVersion() {
	url := "http://" + SensorIp + ":8008/api/version"
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}
	f, err := os.OpenFile(SensorId+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 066)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.Write(body)
	f.WriteString("\n")
}

// 查看设备服务是否正常运行
func SensorServerTest() {
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
			// TODO
		} else if strings.Contains(ret.Status, "Exited") {
			fmt.Printf("%s : 停止运行", name)
		} else {
			fmt.Printf("%s : 其他情况 : %s", name, ret.Status)
		}
		Image[name] = false
	}
	for key, value := range Image {
		if value {
			fmt.Printf("%s : 未运行", key)
		}
	}
}
