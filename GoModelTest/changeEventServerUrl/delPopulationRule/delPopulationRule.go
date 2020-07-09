/*
#Time      :  2020/7/9 2:01 下午
#Author    :  chuangangshen@deepglint.com
#File      :  delPeopleEule.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
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

type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func main() {
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

// 删除人数过多的过滤规则
func runApp(ip string) {
	ids := getPopulationRuleId(ip)
	glog.Infoln(ids)
	delPopulationRule(ip, ids)
}

func delPopulationRule(ip string, ids []string) {
	url := "http://" + ip + "/api/f/eventrule"
	client := &http.Client{}
	postData, err := json.Marshal(ids)
	if err != nil {
		glog.Warningln(err)
		return
	}
	req, err := http.NewRequest("DELETE", url, strings.NewReader(string(postData)))
	if err != nil {
		glog.Warningln(err)
		return
	}
	req.Close = true
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
	glog.Infoln(respData.Msg)
}

func getPopulationRuleId(ip string) (ret []string) {
	url := "http://" + ip + "/api/f/eventrule"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Warningln(err)
		return
	}
	req.Close = true
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
		glog.Warningln(err)
		return
	}
	for id, data := range respData.Data.(map[string]interface{}) {
		if strings.HasPrefix(strings.TrimSpace(id), "population") {
			if 0 == data.(map[string]interface{})["UpperBound"].(float64) {
				ret = append(ret, id)
			}
		}
	}
	return
}
