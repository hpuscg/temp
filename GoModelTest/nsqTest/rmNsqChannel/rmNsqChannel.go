/*
#Time      :  2019/4/4 下午3:00 
#Author    :  chuangangshen@deepglint.com
#File      :  rmNsqChannel.go
#Software  :  GoLand
*/
package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// GetNsqChannel()
	DeleteNsqChannel()
}

func DeleteNsqChannel() {
	url := "http://192.168.100.223:4151/channel/delete?topic=" + "events" + "&channel=" + "bank4door1554374468"
	http.Post(url, "", nil)
}

func GetNsqChannel()  {
	ip := "http://192.168.100.223:4151/stats"
	resp, err := http.Get(ip)
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("===", err)
	}
	//fmt.Println(string(data))
	ret := strings.Split(string(data), "\n")
	// fmt.Println(ret)
	for _, value := range ret {
		// fmt.Println(value)
		channelRet := strings.Split(value, "[")
		if len(channelRet) > 1 {
			channel := strings.Split(channelRet[1], " ")[0]
			if strings.HasPrefix(channel, "bank4door") || strings.HasPrefix(channel, "bank4latch") {
				fmt.Println(channel)
				url := "http://192.168.100.223:4151/channel/delete?topic=" + "events" + "&channel=" + channel
				http.Post(url, "", nil)
			}
		}
	}
}




