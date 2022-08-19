package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gitlab.deepglint.com/junkaicao/glog"
)

func main() {
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"), glog.WithLevel("info"))
	ipListInfo := ReadInfoFromCsv()
	for i, ipInfo := range ipListInfo {
		glog.Infof("===========%d============", i+1)
		glog.Infoln(ipInfo.OldIp)
		// 测试IP是否能ping通
		err := tryPing(ipInfo.OldIp)
		if err != nil {
			glog.Warningf("%s 网络不通，请检查\n", ipInfo.OldIp)
			continue
		}
		SetSensorIp(ipInfo.OldIp, ipInfo.StaticIpConfig)
	}
}

type StaticIpConfig struct {
	Address string `json:"address"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type CsvData struct {
	OldIp string `json:"oldIp"`
	StaticIpConfig
}

func ReadInfoFromCsv() (resultDatas []CsvData) {
	fs, err := os.Open("iplist.csv")
	if err != nil {
		glog.Warningln(err.Error())
		return
	}
	defer fs.Close()

	r := csv.NewReader(fs)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		glog.Warningln(err.Error())
	}

	for index, item := range rows {
		if index == 0 {
			continue
		}
		if len(item) == 4 {
			tempData := CsvData{
				OldIp: item[0],
				StaticIpConfig: StaticIpConfig{
					Address: item[1],
					Netmask: item[2],
					Gateway: item[3],
				},
			}
			resultDatas = append(resultDatas, tempData)
		} else {
			glog.Warningf("line: %d, value: %+v ", index+1, item)
		}
	}
	glog.Infof("%+v", resultDatas)
	return
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

func SetSensorIp(ip string, addressInfo StaticIpConfig) {
	url := "http://" + ip + "/api/staticip"
	contentType := "application/json"
	ret, err := json.Marshal(addressInfo)
	if err != nil {
		glog.Warningln(err.Error())
	}
	data := strings.NewReader(string(ret))
	resp, err := http.Post(url, contentType, data)
	if err != nil {
		glog.Warningln(err.Error())
	}
	glog.Infoln(resp.Status)
}
