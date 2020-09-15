/*
#Time      :  2020/7/24 4:17 下午
#Author    :  chuangangshen@deepglint.com
#File      :  demoForServer.go
#Software  :  GoLand
*/
package main

import (
	"context"
	"encoding/json"
	"gitlab.deepglint.com/junkaicao/glog"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	urlAddr  = "http://192.168.5.251/api/iotinfo"
	fileList = make([]string, 0)
)

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMTQ4MzIzODUsImlzcyI6ImRlZXBnbGludCIsIlVzZXJ" +
		"JRCI6ImJ1bWJsZSJ9.VWGZm5LkQDoyukekwg6KEG-BbAkP28lcpx8D32t5mLw"
)

func main() {
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"), glog.WithLevel("info"))
	glog.Fatalf("%s", "yes")
	glog.Flush()
	glog.Warningln("no")
	return
	GetFileList()
	glog.Infoln(fileList)
}

type Response struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

// 获取文件列表
func GetFileList() {
	timeStamp := int(time.Now().UnixNano() / 1000000)
	timeLength := 10000
	data := make(map[string]int)
	data["timeStamp"] = timeStamp - timeLength
	data["timeLength"] = timeLength
	urler := url.URL{}
	urlProxy, _ := urler.Parse(urlAddr)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				_ = c.SetDeadline(time.Now().Add(10 * time.Second))
				return c, nil
			},
			Proxy: http.ProxyURL(urlProxy),
		},
	}
	reqData, err := json.Marshal(data)
	if err != nil {
		glog.Warningln(err)
		return
	}
	req, err := http.NewRequest(http.MethodGet, urlAddr, strings.NewReader(string(reqData)))
	if err != nil {
		glog.Warningln(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Warningln(err)
		return
	}
	var respBody Response
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		glog.Warningln(err)
		return
	}
	respData := respBody.Data
	if respData != nil {
		switch respData.(type) {
		case []string:
			fileList = respData.([]string)
		}
		glog.Infoln(respData)
	}
}
