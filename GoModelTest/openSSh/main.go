package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	ip, systemType, token string

	urlPaths = []struct {
		Type       string
		SkUrl      string
		OpenSShUrl string
	}{
		{
			Type:       "fg",
			SkUrl:      "http://%s:9000/api/GetDeviceInfo",
			OpenSShUrl: "http://%s:9000/api/SSH",
		},
		{
			Type:       "v1",
			SkUrl:      "http://%s:9000/sys/DeviceInfo",
			OpenSShUrl: "http://%s:9000/sys/SSH",
		},
		{
			Type:       "v2",
			SkUrl:      "http://%s:9000/sys/deviceinfo",
			OpenSShUrl: "http://%s:9000/sys/ssh",
		},
	}
)

func main() {
	flag.StringVar(&ip, "i", "x.x.x.x", "device ip")
	flag.Parse()
	if ip == "x.x.x.x" {
		fmt.Println("please input device ip: -i 192.168.12.12")
		return
	}

	if err := tryPing(ip); err != nil {
		fmt.Printf("%s 网络不通\n", ip)
		return
	}

	GetUUID()

	if err := OpenSSH(); err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 测试设备IP能否ping通
func tryPing(ip string) error {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if sysInfo == "windows" {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}
	_, err := cmd.Output()
	return err
}

func GetUUID() {
	var (
		resp *http.Response
		body []byte
		err  error
	)
	for _, urlData := range urlPaths {
		if resp, err = http.Get(fmt.Sprintf(urlData.SkUrl, ip)); err != nil {
			continue
		}
		if urlData.Type == "fg" {
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			var skReq struct {
				Code     int
				Msg      string
				Redirect string
				Data     map[string]interface{}
			}
			err = json.Unmarshal(body, &skReq)
			if err != nil {
				continue
			}
			if skReq.Data["DeviceModel"].(string) == "HC-BA101" {
				fmt.Printf("%s is TK1\n", ip)
				return
			}
			h := md5.New()
			h.Write([]byte(fmt.Sprintf("deepglint%s", skReq.Data["UUID"])))
			token = hex.EncodeToString(h.Sum(nil))
			systemType = urlData.Type
			return
		} else if urlData.Type == "v1" {
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			var skReq = make(map[string]interface{})
			err = json.Unmarshal(body, &skReq)
			if err != nil {
				continue
			}
			h := md5.New()
			h.Write([]byte(fmt.Sprintf("deepglint%s",
				skReq["data"].(map[string]interface{})["Hardware"].(map[string]interface{})["UUID"].(string))))
			token = hex.EncodeToString(h.Sum(nil))
			systemType = urlData.Type
			return
		} else if urlData.Type == "v2" {
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			var skReq = make(map[string]interface{})
			err = json.Unmarshal(body, &skReq)
			if err != nil {
				continue
			}
			h := md5.New()
			h.Write([]byte(fmt.Sprintf("deepglint%s",
				skReq["data"].(map[string]interface{})["id"].(string))))
			token = hex.EncodeToString(h.Sum(nil))
			systemType = urlData.Type
			return
		}
	}
}

func OpenSSH() (err error) {
	var (
		data, respData []byte
		req            *http.Request
		sshResp        *http.Response
		OpenSShUrl     string
	)
	for _, urlData := range urlPaths {
		if urlData.Type == systemType {
			OpenSShUrl = urlData.OpenSShUrl
		}
	}
	if systemType == "fg" {
		fullSshOpenUrl := fmt.Sprintf(OpenSShUrl, ip)
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
		requestData["SecretKey"] = token
		data, err = json.Marshal(requestData)
		if err != nil {
			return
		}
		req, err = http.NewRequest("POST", fullSshOpenUrl, strings.NewReader(string(data)))
		req.Close = true
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		sshResp, err = client.Do(req)
		if err != nil {
			return
		}
		defer sshResp.Body.Close()
		respData, err = ioutil.ReadAll(sshResp.Body)
		if err != nil {
			return
		}
		fmt.Println(string(respData))
		return
	} else if systemType == "v1" {
		fullSshOpenUrl := fmt.Sprintf(OpenSShUrl, ip)
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
		requestData["SecretKey"] = token
		data, err = json.Marshal(requestData)
		if err != nil {
			return
		}
		req, err = http.NewRequest("POST", fullSshOpenUrl, strings.NewReader(string(data)))
		req.Close = true
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		sshResp, err = client.Do(req)
		if err != nil {
			return
		}
		defer sshResp.Body.Close()
		respData, err = ioutil.ReadAll(sshResp.Body)
		if err != nil {
			return
		}
		fmt.Println(string(respData))
		return
	} else if systemType == "v2" {
		fullSshOpenUrl := fmt.Sprintf(OpenSShUrl, ip)
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
		requestData["open"] = true
		data, err = json.Marshal(requestData)
		if err != nil {
			return
		}
		req, err = http.NewRequest("POST", fullSshOpenUrl, strings.NewReader(string(data)))
		if err != nil {
			return
		}
		req.Close = true
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		sshResp, err = client.Do(req)
		if err != nil {
			return
		}
		defer sshResp.Body.Close()
		respData, err = ioutil.ReadAll(sshResp.Body)
		if err != nil {
			return
		}
		fmt.Println(string(respData))
		return
	} else {
		return fmt.Errorf("%s bad system type: %s", ip, systemType)
	}
}
