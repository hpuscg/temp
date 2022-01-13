/*
#Time      :  2020/7/15 4:12 下午
#Author    :  chuangangshen@deepglint.com
#File      :  setNtpForT1.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"gitlab.deepglint.com/junkaicao/glog"
)

var (
	IpList string
)

const (
	bumbleToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMTQ4MzIzODUsImlzcyI6ImRlZXBnbGludCIsIlVzZXJ" +
		"JRCI6ImJ1bWJsZSJ9.VWGZm5LkQDoyukekwg6KEG-BbAkP28lcpx8D32t5mLw"
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
		runApp(sensorIp)
	}
}

// 测试设备IP能否ping通
func tryPing(ip string) error {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if sysInfo == "windows" {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}
	_, err := cmd.Output()
	return err
}

type RebootInfo struct {
	Delay int
}

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func Reboot(ip string) {
	rebootInfo := RebootInfo{
		Delay: 5,
	}
	ret, err := json.Marshal(rebootInfo)
	if err != nil {
		glog.Infoln(err)
		return
	}
	rebootUrl := "http://" + ip + "/api/reboot"
	client1 := &http.Client{}
	client1.Timeout = 30 * time.Second
	req1, err := http.NewRequest(http.MethodPost, rebootUrl, strings.NewReader(string(ret)))
	req1.Close = true
	if err != nil {
		glog.Infoln(err)
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Add("authorization", bumbleToken)
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

func runApp(ip string) {
	Reboot(ip)
}
