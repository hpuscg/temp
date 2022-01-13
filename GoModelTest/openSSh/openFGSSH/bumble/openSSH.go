/*
#Time      :  2021/1/8 4:14 下午
#Author    :  chuangangshen@deepglint.com
#File      :  openSSH.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	sensorIp    string
	loginUrl    = "/api/a/login"
	sshOpenUrl  = "/api/ssh/enable"
	sshCloseUrl = "/api/ssh/disable"
	token       string
	isOpen      bool
	logFile     = "./openSSH.log"
)

func main() {
	flag.BoolVar(&isOpen, "open", true, "open or close,true is open,false is close")
	flag.Parse()
	if err := LogInit(); err != nil {
		fmt.Println(err)
		return
	}
	fi, err := os.Open("./ip.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for i := 0; i >= 0; i++ {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		sensorIp = string(a)
		if err := GetTokenByPostWithPassWord(); err != nil {
			fmt.Println(err)
			return
		}
		if isOpen {
			OpenSSH()
		} else {
			CloseSSH()
		}
	}
}

type LoginIn struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type Response struct {
	Code     int
	Msg      string
	Redirect string
	Data     interface{}
}

func CloseSSH() {
	fullSshCloseUrl := "http://" + sensorIp + sshCloseUrl
	urler := url.URL{}
	url_proxy, _ := urler.Parse(fullSshCloseUrl)
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
	req, err := http.NewRequest("PUT", fullSshCloseUrl, nil)
	req.Close = true
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data Response
	if err := json.Unmarshal(respData, &data); err != nil {
		fmt.Println(err)
	}
	if 0 == data.Code {
		fmt.Println(sensorIp, "close ssh success")
	} else {
		fmt.Println(sensorIp, "close ssh fail")
		fmt.Println(data.Msg)
	}
}

func OpenSSH() {
	fullSshOpenUrl := "http://" + sensorIp + sshOpenUrl
	urler := url.URL{}
	url_proxy, _ := urler.Parse(fullSshOpenUrl)
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
	req, err := http.NewRequest("PUT", fullSshOpenUrl, nil)
	req.Close = true
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data Response
	if err := json.Unmarshal(respData, &data); err != nil {
		fmt.Println(err)
	}
	if 0 == data.Code {
		fmt.Println(sensorIp, "open ssh success")
	} else {
		fmt.Println(sensorIp, "open ssh fail")
		fmt.Println(data.Msg)
	}
}

// 通过用户名密码获取token
func GetTokenByPostWithPassWord() (err error) {
	var loginData = LoginIn{UserName: "admin", PassWord: "B32FA6018771638F277F0BE418708C10"}
	byteData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fullLoginUrl := "http://" + sensorIp + loginUrl
	urler := url.URL{}
	url_proxy, _ := urler.Parse(fullLoginUrl)
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
	req, err := http.NewRequest("POST", fullLoginUrl, strings.NewReader(string(byteData)))
	req.Close = true
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Add("authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data Response
	if err := json.Unmarshal(respData, &data); err != nil {
		fmt.Println(err)
	}
	token = data.Data.(map[string]interface{})["token"].(string)
	// fmt.Println(token)
	return
}

func LogInit() error {
	logFileFd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed", err)
		return err
	}
	log.SetOutput(logFileFd)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
	return nil
}
