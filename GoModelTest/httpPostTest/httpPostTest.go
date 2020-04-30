/*
#Time      :  2019/9/27 下午5:23 
#Author    :  chuangangshen@deepglint.com
#File      :  httpPostTest.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	// PostIotInfo()
	for {
		GetPeopleNumTest()
		// time.Sleep(1 * time.Microsecond)
	}
}

type FoundationIotInfo struct {
	Iotserver string `json:"iotserver"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Topic     string `json:"topic"`
}

func GetPeopleNumTest() {
	url := "http://192.168.5.250:8008/api/people/num"
	result, err := HTTPGet(url)
	if err != nil {
		time.Sleep(1 * time.Second)
		fmt.Println(err)
	}
	peopleNum, err := strconv.Atoi(string(result))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(peopleNum)
}

func PostIotInfo()  {
	var foundationIonInfo FoundationIotInfo
	foundationIonInfo.Iotserver = "tcp://"
	jsonData, _ := json.Marshal(foundationIonInfo)
	resp, _ := http.NewRequest("POST", "http://192.168.5.250:8180/api/iotInfo", strings.NewReader(string(jsonData)))
	// resp, _ := http.Post("http://192.168.5.250:8180/api/iotInfo", "", strings.NewReader(string(jsonData)))
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}

func HTTPGet(url_addr string) (result []byte, err error) {
	urler := url.URL{}
	url_proxy, _ := urler.Parse(url_addr)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second))
				return c, nil
			},
			Proxy: http.ProxyURL(url_proxy),
		},
	}
	reqest, err := http.NewRequest("GET", url_addr, nil)
	if err != nil {
		return
	}
	response, err := client.Do(reqest)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	result = body
	return
}
