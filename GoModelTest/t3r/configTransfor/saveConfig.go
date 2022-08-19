package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func saveConfig(ip string) {
	var (
		deviceInfo DeviceInfo
		err        error
	)
	deviceInfo.Ip = ip
	if deviceInfo.Name, err = getDeviceName(ip); err != nil {
		logger.Println(err.Error())
		return
	}
	if deviceInfo.IotInfo, err = getIotInfo(ip); err != nil {
		logger.Println(err.Error())
		return
	}
	if deviceInfo.EventServerInfo, err = getEventServerInfo(ip); err != nil {
		logger.Println(err.Error())
		return
	}
	if deviceInfo.NtpInfo, err = getNtpInfo(ip); err != nil {
		logger.Println(err.Error())
		return
	}

	deviceInfoList = append(deviceInfoList, deviceInfo)
}

// 获取设备事件配置

// 获取设备较时信息
func getNtpInfo(ip string) (ntpInfo Ntp, err error) {
	var data []byte
	ntpUrl := fmt.Sprintf("http://%s/api/synctime", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodGet).SetURL(ntpUrl).
		SetHeader("authorization", BumbleToken).Do(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     Ntp    `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	ntpInfo = resp.Data
	return
}

// 获取事件上传地址
func getEventServerInfo(ip string) (eventserver EventServer, err error) {
	var data []byte
	eventServerUrl := fmt.Sprintf("http://%s/api/f/logininfo", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodGet).SetURL(eventServerUrl).
		SetHeader("authorization", BumbleToken).Do(); err != nil {
		return
	}
	resp := struct {
		Code     int         `json:"code"`
		Msg      string      `json:"msg"`
		Redirect string      `json:"redirect"`
		Data     EventServer `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	eventserver = resp.Data
	return
}

// 获取iot信息
func getIotInfo(ip string) (iotInfo Iot, err error) {
	var data []byte
	iotUrl := fmt.Sprintf("http://%s/api/iotinfo", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodGet).SetURL(iotUrl).
		SetHeader("authorization", BumbleToken).Do(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     Iot    `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	iotInfo = resp.Data
	return
}

// 获取设备名称
func getDeviceName(ip string) (name string, err error) {
	var data []byte
	nameUrl := fmt.Sprintf("http://%s/api/name", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodGet).SetURL(nameUrl).
		SetHeader("authorization", BumbleToken).Do(); err != nil {
		return
	}
	resp := struct {
		Code     int               `json:"code"`
		Msg      string            `json:"msg"`
		Redirect string            `json:"redirect"`
		Data     map[string]string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	name = resp.Data["sensor_desc"]
	return
}
