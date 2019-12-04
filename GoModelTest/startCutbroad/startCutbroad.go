/*
#Time      :  2019/4/12 下午1:13 
#Author    :  chuangangshen@deepglint.com
#File      :  startCutbroad.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"net/http"
	"strings"
	"flag"
)

func main() {
	var (
		Ip string
		isCut string
		bufferSize string
		interval string
	)
	flag.StringVar(&Ip, "ip", "", "")
	flag.StringVar(&isCut, "value", "0", "")
	flag.StringVar(&bufferSize, "buffer", "100", "picture in tmp size")
	flag.StringVar(&interval, "interval", "20000", "save cutboard time")
	flag.Parse()
	fmt.Printf("ip=%s,value=%s,buffersize=%s,interval=%s", Ip, isCut, bufferSize, interval)
	StartBodyCutbroad(Ip, isCut)
	StartFallCutbroad(Ip, isCut)
	SetBufferSize(Ip, bufferSize)
	SetInterval(Ip, interval)
}

func StartBodyCutbroad(Ip, isCut string) {
	fmt.Println(Ip, "===")
	bodyUrl := "http://" + Ip + ":8008/api/cutboard/body"
	bodyBuf := strings.NewReader("iscut=" + isCut)
	contentType := "application/x-www-form-urlencoded"
	http.Post(bodyUrl , contentType, bodyBuf)
}

func StartFallCutbroad(Ip, isCut string)  {
	fallUrl := "http://" + Ip + ":8008/api/cutboard/falldown"
	bodyBuf := strings.NewReader("iscut=" + isCut)
	contentType := "application/x-www-form-urlencoded"
	http.Post(fallUrl , contentType, bodyBuf)
}

func SetBufferSize(Ip, buffer string) {
	fallUrl := "http://" + Ip + ":8008/api/cutboard/buffersize"
	bodyBuf := strings.NewReader("buffersize=" + buffer)
	contentType := "application/x-www-form-urlencoded"
	http.Post(fallUrl , contentType, bodyBuf)
}

func SetInterval(Ip, interval string) {
	fallUrl := "http://" + Ip + ":8008/api/cutboard/interval"
	bodyBuf := strings.NewReader("interval=" + interval)
	contentType := "application/x-www-form-urlencoded"
	http.Post(fallUrl , contentType, bodyBuf)
}

