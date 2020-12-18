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
	IpList    string
	NtpServer string
)

const (
	flowserviceToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTg0OTQ4MzMsImlzcyI6ImRlZXBnbGludCIsIl" +
		"VzZXJJRCI6ImZsb3dzZXJ2aWNlIn0.ubMv0T3FTVURQG2E6YPForuq4ixX_sq5nI0JXn4Q6Io"
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
	if "windows" == sysInfo {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}
	_, err := cmd.Output()
	return err
}

type DateSource struct {
	Mode int
	Ntp  string
	Date string
}

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func GetNtpServer(ip string) {
	getNtpServerUrl := "http://" + ip + "/api/f/logininfo"
	client := &http.Client{}
	client.Timeout = 30 * time.Second
	req, err := http.NewRequest("GET", getNtpServerUrl, nil)
	req.Close = true
	if err != nil {
		glog.Warningln(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", flowserviceToken)
	resp, err := client.Do(req)
	var respData ResponseData
	if resp == nil || err != nil {
		glog.Warningln(err)
		return
	}
	defer resp.Body.Close()
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Warningln(err)
		return
	}
	err = json.Unmarshal(byteData, &respData)
	if err != nil {
		glog.Warningln(err)
		return
	}
	switch tt := respData.Data.(type) {
	case map[string]interface{}:
		for key, value := range tt {
			if "url" == key {
				switch value.(type) {
				case string:
					NtpServer = strings.TrimSuffix(strings.TrimPrefix(value.(string), "http://"), ":8888")
				}
			}
		}
	}
}

func SetNtpServer(ip string) {
	iotInfo := DateSource{
		Mode: 1,
		Ntp:  NtpServer,
		Date: "",
	}
	ret, err := json.Marshal(iotInfo)
	if err != nil {
		glog.Infoln(err)
		return
	}
	ntpSetUrl := "http://" + ip + "/api/synctime"
	client1 := &http.Client{}
	client1.Timeout = 30 * time.Second
	req1, err := http.NewRequest("PUT", ntpSetUrl, strings.NewReader(string(ret)))
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
	GetNtpServer(ip)
	SetNtpServer(ip)
}
