/*
#Time      :  2020/6/17 10:50 上午
#Author    :  chuangangshen@deepglint.com
#File      :  saveH264.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"temp/GoModelTest/saveH264/yaml"
	"time"
)

var (
	serverIp         string
	sensorId         string
	bumbleYamlClient *yaml.YamlConfig
)

const (
	startFile    = "/tmp/START_VI"
	preH264Dir   = "/tmp/DC/Record/LibraT/"
	configFile   = "/home/deepglint/AppData/libraT/config/DCConfig.json"
	sensorIdFile = "/home/deepglint/AppData/libraT/config/libraT.yaml"
	sensorIdKey  = "/libraT/sensor_uid"
)

func main() {
	flag.StringVar(&serverIp, "serverIp", "", "文件服务器")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("logs"), glog.WithLevel("info"))
	getServerIp()
	getSensorId()
	// checkStartSave()
	saveH264()
}

func checkStartSave() {
	_, err := os.Stat(startFile)
	if os.IsNotExist(err) {
		f, err := os.Create(startFile)
		if err != nil {
			glog.Infoln(err)
		}
		defer f.Close()
	}
}

func getSensorId() {
	var err error
	bumbleYamlClient, err = yaml.NewYamlConfig(sensorIdFile)
	if err != nil {
		glog.Infoln(err)
		return
	}
	sensorId, err = bumbleYamlClient.GetString(sensorIdKey)
	if err != nil {
		glog.Infoln(err)
		return
	}
}

func getServerIp() {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		glog.Infoln("config file not found")
		return
	}
	fileData, err := ioutil.ReadFile(configFile)
	data := make(map[string]interface{})
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		glog.Infoln(err)
		return
	}
	serverIp = data["file_upload_server_ip"].(string)
}

func saveH264() {
	h264FileDir := filepath.Join(preH264Dir, sensorId)
	files, err := ioutil.ReadDir(h264FileDir)
	if err != nil {
		glog.Infoln(err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".h264") {
			nowTime := time.Now().Unix()
			fileModifyTime := f.ModTime().Unix()
			if nowTime - fileModifyTime > 60 {
				glog.Infoln(nowTime)
				glog.Infoln(fileModifyTime)
				rsyncStr := "-aopgtvp  /tmp/DC/Record root@" + serverIp + ":/home/ubuntu/"
				cmd := exec.Command("rsync", rsyncStr)
				cmd.Run()
				rmStr := filepath.Join(h264FileDir, f.Name())
				cmd2 := exec.Command("rm", "-rf", rmStr)
				cmd2.Run()
			}
		}
	}
}
