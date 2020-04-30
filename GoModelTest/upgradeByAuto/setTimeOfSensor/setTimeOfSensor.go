/*
#Time      :  2020/4/20 6:41 下午
#Author    :  chuangangshen@deepglint.com
#File      :  setTimeOfSensor.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"gitlab.deepglint.com/junkaicao/glog"
)

func main() {

}

// 设置NTP较时的结构
type DataSource struct {
	Mode int
	Ntp  string
}

// 给设备较时
func SetTimeOfSensor(ip string) (ok bool) {
	httpUrl := "http://" + ip + ":8008/api/synctime"
	data := DataSource{
		Mode: 1,
		Ntp:  ServerIp,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		glog.Infoln(err)
	}
	_, code, err := HttpPost(httpUrl, string(jsonData))
	if err != nil {
		glog.Infoln(err)
	}
	if code == 200 {
		ok = true
		return
	} else if code == 404 {
		httpUrl = "http://" + ip + ":8008/api/synctime/update"
		data := DataSource{
			Mode: 1,
			Ntp:  ServerIp,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			glog.Infoln(err)
		}
		_, code, err := HttpPost(httpUrl, string(jsonData))
		if err != nil {
			glog.Infoln(err)
		}
		if code == 200 {
			ok = true
			return
		} else {
			glog.Infoln(code)
		}
	}
	return
}