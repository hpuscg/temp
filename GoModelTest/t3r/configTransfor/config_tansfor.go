package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	useInfo           string
	deviceInfoList    []DeviceInfo
	ipFile            string
	deviceModel       string
	cleandataFaildIps []string
)

func main() {
	flag.StringVar(&useInfo, "useInfo", "load", "set tools function: save config or load config")
	flag.StringVar(&ipFile, "ipFile", "ip.txt", "ip list of device")
	flag.StringVar(&deviceModel, "model", "HC-BA421-T", "the model of device")
	flag.Parse()
	initLog()
	initHttp()
	initConfigViper()
	logger.Println("################ begin #################")
	if useInfo == "save" {
		fi, err := os.Open(ipFile)
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		for i := 0; i >= 0; i++ {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			sensorIp := string(a)
			logger.Printf("======%d======%s======\n", i+1, sensorIp)
			// 测试IP是否能ping通
			if err := tryPing(sensorIp); err != nil {
				logger.Printf("%s 网络不通，请检查\n", sensorIp)
				continue
			}
			if tempModel, err := checkDeviceModel(sensorIp); err != nil || tempModel != deviceModel {
				logger.Printf("check model err, model: %s, err: %+v\n", tempModel, err)
				continue
			}
			saveConfig(sensorIp)
		}
		fmt.Printf("%+v\n", deviceInfoList)
		configViper.Set("deviceInfos", deviceInfoList)
		if err := configViper.WriteConfig(); err != nil {
			logger.Println(err.Error())
		}
	} else if useInfo == "load" {
		if err := configViper.UnmarshalKey("deviceInfos", &deviceInfoList); err != nil {
			logger.Println(err.Error())
			return
		} else if len(deviceInfoList) == 0 {
			logger.Println("config no data")
			return
		}
		// 简单恢复出厂
		for _, deviceInfo := range deviceInfoList {
			if tempModel, err := checkDeviceModel(deviceInfo.Ip); err != nil || tempModel != deviceModel {
				logger.Printf("check model err, model: %s, err: %+v\n", tempModel, err)
				continue
			}
			if err := cleanSensorData(deviceInfo.Ip); err != nil {
				logger.Println(err.Error())
				cleandataFaildIps = append(cleandataFaildIps, deviceInfo.Ip)
			}
		}
		logger.Println("time wait 300s")
		time.Sleep(300 * time.Second)
		// 恢复配置
		for _, deviceInfo := range deviceInfoList {
			if tempModel, err := checkDeviceModel(deviceInfo.Ip); err != nil || tempModel != deviceModel {
				logger.Printf("check model err, model: %s, err: %+v\n", tempModel, err)
				continue
			}
			for _, tempIp := range cleandataFaildIps {
				if tempIp == deviceInfo.Ip {
					continue
				}
			}
			loadConfig(deviceInfo)
		}
		logger.Printf("faild recover config device : %+v\n", cleandataFaildIps)
	} else {
		fmt.Printf("use info: %s", useInfo)
		fmt.Println("please input tools function, save or load")
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

func checkDeviceModel(ip string) (model string, err error) {
	var data []byte
	modelUrl := fmt.Sprintf("http://%s/api/model", ip)
	if _, data, err = NewHttpRequest().SetMethod(http.MethodGet).SetURL(modelUrl).
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
	model = resp.Data["model"]
	return
}
