/*
#Time      :  2020/10/9 5:09 下午
#Author    :  chuangangshen@deepglint.com
#File      :  setClip.go
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
	"runtime"
	"strings"
	"time"
)

var (
	IpList string
)

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

const (
	libraToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTg0OTQ4MjYsImlzcyI6ImRlZXBnbGludCIsIlVzZXJJ" +
		"RCI6ImxpYnJhIn0.EvValuL84BiVRAnwRIkMnjPPEOWgoDPxGjZvsbpE2bE"
)

func main() {
	flag.StringVar(&IpList, "ipList", "./ip.txt", "sensor ip list file")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"), glog.WithLevel("info"))
	fi, err := os.Open(IpList)
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
		sensorIp := string(a)
		glog.Infof("===========%d============", i+1)
		glog.Infoln(sensorIp)
		// 测试IP是否能ping通
		err := tryPing(sensorIp)
		if err != nil {
			glog.Infof("%s 网络不通，请检查\n", sensorIp)
			continue
		}
		setClip(sensorIp)
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

func setClip(ip string) {
	clipInfo := map[string]bool{
		"enable": true,
	}
	ret, err := json.Marshal(clipInfo)
	if err != nil {
		glog.Infoln(err)
		return
	}
	clipSetUrl := "http://" + ip + "/api/l/ClipStorageEnable"
	client1 := &http.Client{}
	client1.Timeout = 30 * time.Second
	req1, err := http.NewRequest("POST", clipSetUrl, strings.NewReader(string(ret)))
	req1.Close = true
	if err != nil {
		glog.Infoln(err)
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Add("authorization", libraToken)
	response1, err := client1.Do(req1)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer response1.Body.Close()
	body1, err := ioutil.ReadAll(response1.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	var respData1 ResponseData
	err = json.Unmarshal(body1, &respData1)
	if err != nil {
		glog.Infoln(err)
	}
	glog.Infoln(respData1.Msg)
}
