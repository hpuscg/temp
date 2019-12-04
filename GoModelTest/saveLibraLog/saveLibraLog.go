/*
#Time      :  2019/7/24 上午10:20 
#Author    :  chuangangshen@deepglint.com
#File      :  saveLibraLog.go
#Software  :  GoLand
*/
package main

import (
	"flag"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
)

func main() {
	var ip string
	flag.StringVar(&ip, "ip", "192.168.19.247", "sensor ip")
	flag.Parse()
	url := "http://" + ip + ":8008/api/libra/logs"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	var data []string
	ret, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(ret, &data)
	filePath := "/Users/hpu_scg/gocode/src/temp/GoModelTest/saveLibraLog/test.txt"
	_, err = os.Open(filePath)
	if os.IsNotExist(err) {
		os.Create(filePath)
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	osFileInfo, err := f.Stat()
	fmt.Printf("%+v\n", osFileInfo.Size()/1024)
	/*for _, line := range data {
		f.Write([]byte(line))
	}*/
}
