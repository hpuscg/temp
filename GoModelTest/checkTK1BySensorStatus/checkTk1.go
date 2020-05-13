/*
#Time      :  2020/5/13 10:46 上午
#Author    :  chuangangshen@deepglint.com
#File      :  checkTk1.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	IpListFile string
	Port       int
	logDir     string
	alsoToFile bool
	logLevel   string
	serverIp   string
)

type SensorStatus struct {
	Usb           int
	Disk          int
	TfCard        int
	LocalNetwork  int
	RemoteNetwork int
	Memory        int
	Service       int
}

func main() {
	flag.StringVar(&IpListFile, "ipListFile", "./serverIp.txt", "sensor ip list")
	flag.StringVar(&logDir, "log_dir", "logs", "log dir, default /tmp")
	flag.BoolVar(&alsoToFile, "alsologtostderr", true, "log to stderr also to log file")
	flag.StringVar(&logLevel, "log_level", "info", "log level, default info")
	flag.IntVar(&Port, "port", 22, "ssh port")
	flag.Parse()
	// 判断是否已有历史log，如有进行移动
	timeStamp := time.Now().Unix()
	stringTimeStamp := strconv.Itoa(int(timeStamp))
	newLogFileName := filepath.Join(logDir, stringTimeStamp+".log")
	fileNameArr := strings.Split(os.Args[0], "/")
	oldLogFileName := filepath.Join(logDir, fileNameArr[len(fileNameArr)-1]+".log")
	_, err := os.Stat(oldLogFileName)
	if err == nil {
		cmd := exec.Command("mv", oldLogFileName, newLogFileName)
		_ = cmd.Run()
	}
	// 初始化glog配置
	glog.Config(glog.WithAlsoToStd(alsoToFile), glog.WithFilePath(logDir))
	// 逐行读取配置文件中的设备IP
	fi, err := os.Open(IpListFile)
	if err != nil {
		glog.Infof("Error: %s\n", err)
		return
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for i := 0; i >= 0; i++ {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		glog.Infof("***********************%d***********************", i+1)
		serverIp = string(a)
		glog.Infoln("serverIp:", serverIp)
		// 测试IP是否能ping通
		err := tryPing(serverIp)
		if err != nil {
			glog.Infof("网管服务器%s 网络不通，请检查", serverIp)
			continue
		}
		// GetSensorStatus(SensorIp)
		GetSensorList(i)
	}
}

// 获取sensor status
func GetSensorStatus(ip string) {
	var sS SensorStatus
	url := "http://" + ip + ":8008/api/sensorstatus"
	resp, err := http.Get(url)
	if err != nil {
		glog.Infoln(err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	err = json.Unmarshal(data, &sS)
	if err != nil {
		glog.Infoln(err)
	} else {
		// glog.Infof("%+v", sS)
		if sS.Disk == 0 && sS.Memory == 0 && sS.Service == 0 && sS.Usb == 0 {
			glog.Infof("%s正常:%+v", ip, sS)
		} else {
			glog.Infof("%s异常:Disk:%s,Memory:%s,Service:%s,Usb:%s", ip, sS.Disk, sS.Memory, sS.Service, sS.Usb)
		}
	}
}

// 测试设备IP能否ping通
func tryPing(ip string) error {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if "windows" == sysInfo {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}
	_, err := cmd.Output()
	return err
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
func GetSensorList(i int) {
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
	for j, sensor := range data {
		// glog.Infof("%+v", sensor)
		glog.Infof("============%d-%d===========", i+1, j+1)
		if sensor.Host.Version == "V2.15.200119A" {
			// 测试IP是否能ping通
			sensorIp := sensor.Host.HostIp
			err := tryPing(sensorIp)
			if err != nil {
				glog.Infof("设备%s 网络不通，请检查", sensorIp)
				continue
			}
			GetSensorStatus(sensorIp)
		}
	}
}
