package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func loadConfig(deviceInfo DeviceInfo) {
	if err := setSensorName(deviceInfo.Ip, deviceInfo.Name); err != nil {
		logger.Println(err.Error())
	}
	if err := setIotInfo(deviceInfo.Ip, deviceInfo.IotInfo); err != nil {
		logger.Println(err.Error())
	}
	if err := setDateInfo(deviceInfo.Ip, deviceInfo.NtpInfo); err != nil {
		logger.Println(err.Error())
	}
	if err := setEventServerInfo(deviceInfo.Ip, deviceInfo.EventServerInfo); err != nil {
		logger.Println(err.Error())
	}
	if err := setLocalVedioSave(deviceInfo.Ip); err != nil {
		logger.Println(err.Error())
	}
}

// 简单恢复出厂
func cleanSensorData(ip string) (err error) {
	var data []byte
	cleanDataUrl := fmt.Sprintf("http://%s/api/cleandata", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPut).SetURL(cleanDataUrl).
		SetHeader("authorization", BumbleToken).Do(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	fmt.Println(string(data))
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 设置设备名称
func setSensorName(ip, name string) (err error) {
	var data, byteNameData []byte
	nameUrl := fmt.Sprintf("http://%s/api/name", ip)
	nameData := make(map[string]string)
	nameData["sensor_desc"] = name
	if byteNameData, err = json.Marshal(nameData); err != nil {
		return
	}
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPut).SetURL(nameUrl).
		SetHeader("authorization", BumbleToken).SetBody(string(byteNameData)).DoWithBody(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 设置设备iot信息
func setIotInfo(ip string, iotInfo Iot) (err error) {
	var data, byteIotData []byte
	iotUrl := fmt.Sprintf("http://%s/api/iotinfo", ip)
	if byteIotData, err = json.Marshal(iotInfo); err != nil {
		return
	}
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPut).SetURL(iotUrl).
		SetHeader("authorization", BumbleToken).SetBody(string(byteIotData)).DoWithBody(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 设置设备较时信息
func setDateInfo(ip string, ntpInfo Ntp) (err error) {
	var data, byteNtpData []byte
	ntpUrl := fmt.Sprintf("http://%s/api/synctime", ip)
	if byteNtpData, err = json.Marshal(ntpInfo); err != nil {
		return
	}
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPut).SetURL(ntpUrl).
		SetHeader("authorization", BumbleToken).SetBody(string(byteNtpData)).DoWithBody(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 设置设备事件上传地址
func setEventServerInfo(ip string, eventserverInfo EventServer) (err error) {
	var data, byteEventServerData []byte
	eventserverUrl := fmt.Sprintf("http://%s/api/f/logininfo", ip)
	if byteEventServerData, err = json.Marshal(eventserverInfo); err != nil {
		return
	}
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPost).SetURL(eventserverUrl).
		SetHeader("authorization", BumbleToken).SetBody(string(byteEventServerData)).DoWithBody(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 开启设备本地视频存储
func setLocalVedioSave(ip string) (err error) {
	var data, byteVedioData []byte
	vedioUrl := fmt.Sprintf("http://%s/api/l/ClipStorageEnable", ip)
	vedioData := make(map[string]bool)
	vedioData["sensor_desc"] = true
	if byteVedioData, err = json.Marshal(vedioData); err != nil {
		return
	}
	if _, data, err = NewHttpRequest().SetMethod(http.MethodPost).SetURL(vedioUrl).
		SetHeader("authorization", BumbleToken).SetBody(string(byteVedioData)).DoWithBody(); err != nil {
		return
	}
	resp := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Redirect string `json:"redirect"`
		Data     string `json:"data"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	if resp.Code != 0 {
		err = errors.New(resp.Msg)
	}
	return
}

// 设置设备事件配置
