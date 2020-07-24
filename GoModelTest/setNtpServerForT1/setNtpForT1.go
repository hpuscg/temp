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
	"fmt"
	"gitlab.deepglint.com/junkaicao/glog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	IpList    string
	NtpServer string
)

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMTQ4MzIzODUsImlzcyI6ImRlZXBnbGludCIsIlVzZXJ" +
		"JRCI6ImJ1bWJsZSJ9.VWGZm5LkQDoyukekwg6KEG-BbAkP28lcpx8D32t5mLw"
)

func main() {
	flag.StringVar(&IpList, "ipList", "./ip.txt", "sensor ip list file")
	flag.StringVar(&NtpServer, "ntpServer", "", "ntp server")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"), glog.WithLevel("info"))
	if "" == NtpServer {
		glog.Warningln("please input ntp server")
		return
	}
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
		url := "http://" + sensorIp + "/api/synctime"
		runApp(url)
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

type DateSource struct {
	Mode   int
	Server string
	Ntp    string
	Date   string
}

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func runApp(url string) {
	iotInfo := DateSource{
		Server: "",
		Mode:   1,
		Ntp:    "192.168.101.8",
		Date:   "",
	}
	ret, err := json.Marshal(iotInfo)
	if err != nil {
		fmt.Println(err)
	}
	client1 := &http.Client{}
	req1, err := http.NewRequest("PUT", url, strings.NewReader(string(ret)))
	req1.Close = true
	if err != nil {
		glog.Infoln(err)
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Add("authorization", token)
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
