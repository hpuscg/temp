/*
#Time      :  2020/6/24 8:18 下午
#Author    :  chuangangshen@deepglint.com
#File      :  changeEventServerUrl.go
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
)

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTg0OTQ4MzMsImlzcyI6ImRlZXBnbGludCI" +
		"sIlVzZXJJRCI6ImZsb3dzZXJ2aWNlIn0.ubMv0T3FTVURQG2E6YPForuq4ixX_sq5nI0JXn4Q6Io"
)

var (
	serverIp string
)

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func main() {
	flag.StringVar(&serverIp, "serverIp", "10.147.152.244", "event server ip")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"), glog.WithLevel("info"))
	fi, err := os.Open("./ip.txt")
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
		glog.Infoln(sensorIp)
		// 测试IP是否能ping通
		err := tryPing(sensorIp)
		if err != nil {
			glog.Infof("%s 网络不通，请检查\n", sensorIp)
			continue
		}
		url := "http://" + sensorIp + "/api/f/logininfo"
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

func runApp(url string) {
	/*client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Close = true
	if err != nil {
		glog.Infoln(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", token)
	response, err := client.Do(req)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	var respData ResponseData
	err = json.Unmarshal(body, &respData)
	if err != nil {
		glog.Infoln(err)
	}
	var eventUrl string
	if respData.Data != nil {
		eventUrl = respData.Data.(map[string]interface{})["url"].(string)
	}
	if !strings.HasPrefix(eventUrl, "http://") && "" != eventUrl {
		eventUrl = "http://" + eventUrl
	} else if "" == eventUrl {
		eventUrl = "http://10.147.152.244:8888"
	}
	glog.Infoln(eventUrl)
	*/
	requestData := make(map[string]string)
	requestData["url"] = "http://" + serverIp + ":8888"
	requestData["user_name"] = "haomuT"
	requestData["pass_word"] = "abc@Dgsh"
	postData, err := json.Marshal(requestData)
	if err != nil {
		glog.Infoln(err)
	}
	client1 := &http.Client{}
	req1, err := http.NewRequest("POST", url, strings.NewReader(string(postData)))
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
	if respData1.Data != nil {
		glog.Infoln(respData1.Code)
	} else {
		glog.Infoln("post err")
	}
}
