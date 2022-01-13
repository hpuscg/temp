/*
#Time      :  2021/1/8 4:14 下午
#Author    :  chuangangshen@deepglint.com
#File      :  openSSH.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	sensorIp string
	sk       string
)

type Response struct {
	Code     int
	Msg      string
	Redirect string
	Data     map[string]interface{}
}

func main() {
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
		if err := get_sk(sensorIp); err == nil {
			OpenSSH(sensorIp)
		}
	}
}

func get_sk(ip string) error {
	skUrl := fmt.Sprintf("http://%s:9000/api/GetDeviceInfo", ip)
	resp, err := http.Get(skUrl)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	var skReq Response
	err = json.Unmarshal(body, &skReq)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("deepglint%s", skReq.Data["UUID"])))
	sk = hex.EncodeToString(h.Sum(nil))
	return nil
}

func OpenSSH(ip string) {
	fullSshOpenUrl := fmt.Sprintf("http://%s:9000/api/SSH", ip)
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
	requestData := make(map[string]interface{})
	requestData["On"] = 1
	requestData["SecretKey"] = sk
	data, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest("POST", fullSshOpenUrl, strings.NewReader(string(data)))
	req.Close = true
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respData))
}
