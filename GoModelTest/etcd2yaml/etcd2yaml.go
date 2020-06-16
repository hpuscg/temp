/*
#Time      :  2020/6/15 10:25 上午
#Author    :  chuangangshen@deepglint.com
#File      :  etcd2yaml.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"temp/GoModelTest/etcd2yaml/yaml"
)

var (
	sensorInfos           = make(map[string]string)
	ruleInfos             = make(map[string]interface{})
	configPath            string
	bumbleYamlClient      *yaml.YamlConfig
	flowserviceYamlClient *yaml.YamlConfig
	sensorUid             string
)

const (
	sensorUidKey    = "/config/global/sensor_uid"
	sensorSnKey     = "/config/global/sensor_sn"
	ntpAddrKey      = "/config/global/ntp_addr"
	sensorDescKey   = "/config/global/sensor_desc"
	iotUrlKey       = "/config/iot/server"
	eventUrlKey     = "/config/eventserver/pub_http_url"
	bumbleFile      = "config/libraT.yaml"
	flowserviceFile = "config/flowservice.yaml"
)

const (
	sensorUidYamlKey  = "/libraT/sensor_uid"
	sensorSnYamlKey   = "/libraT/sensor_sn"
	ntpAddrYamlKey    = "/datesync/ntpaddr"
	sensorDescYamlKey = "/libraT/sensor_desc"
	iotUrlYamlKey     = "/iot/server"
	eventUrlYamlKey   = "/config/eventserver/pub_http_url"
	preRuleKey        = "/config/rule"
	preBasicKey       = "/config/basic"
	pubDbUrl          = "127.0.0.1:8880/api/db"
	pubViboUrl        = "127.0.0.1:8880/api/vibo"
)

type EventRule struct {
	Id           string `json:"Id"`
	Enable       bool   `json:"Enable"`
	TimeRange    [2]int `json:"TimeRange"`
	WeekdayRange byte   `json:"WeekdayRange"`
	UpperBound   float64
	LowerBound   float64
}

func main() {
	flag.StringVar(&configPath, "configPath", "./libraT", "配置文件目录")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("logs"), glog.WithLevel("info"))
	GetInfoFromCmd()
	GetEventRule()
	setInfoToConfig()
	mvLibraConfig()
}

func mvLibraConfig() {
	latchOldFile := "/libra/dataset/scene/latch.dat"
	latchNewFile := filepath.Join(configPath, "config/latch.dat")
	_, err := os.Stat(latchOldFile)
	if !os.IsNotExist(err) {
		cmd := exec.Command("cp", latchOldFile, latchNewFile)
		err = cmd.Run()
		if err != nil {
			glog.Infoln(err)
		}
	}
	bgModelOldFile := "/libra/dataset/scene/MoG_bgmodel.dat"
	bgModelNewFile := filepath.Join(configPath, "works/libra/scene/MoG_bgmodel.dat")
	_, err = os.Stat(bgModelOldFile)
	if !os.IsNotExist(err) {
		cmd := exec.Command("cp", bgModelOldFile, bgModelNewFile)
		err = cmd.Run()
		if err != nil {
			glog.Infoln(err)
		}
	}
	mvFloorModelFile()
}

func mvFloorModelFile() {
	floorOldFile := "/libra/dataset/scene/floor_calib.dat"
	floorModelOldFile := "/libra/dataset/scene/floor_model_m.dat"
	floorNewFile := filepath.Join(configPath, "works/libra/scene/floor_calib.dat")
	_, err1 := os.Stat(floorOldFile)
	_, err2 := os.Stat(floorModelOldFile)
	if os.IsNotExist(err2) {
		if !os.IsNotExist(err1) {
			cmd := exec.Command("cp", floorOldFile, floorNewFile)
			err := cmd.Run()
			if err != nil {
				glog.Infoln(err)
			}
		}
	} else {
		if os.IsNotExist(err1) {
			cmd := exec.Command("cp", floorModelOldFile, floorNewFile)
			err := cmd.Run()
			if err != nil {
				glog.Infoln(err)
			}
		} else {
			f1, err := os.Open(floorOldFile)
			if err != nil {
				glog.Infoln(err)
			}
			f2, err := os.Open(floorModelOldFile)
			if err != nil {
				glog.Infoln(err)
			}
			f1Info, err := f1.Stat()
			if err != nil {
				glog.Infoln(err)
			}
			f1ModifyTime := f1Info.ModTime().Unix()
			f2Info, err := f2.Stat()
			if err != nil {
				glog.Infoln(err)
			}
			f2ModifyTime := f2Info.ModTime().Unix()
			if f1ModifyTime > f2ModifyTime {
				cmd := exec.Command("cp", floorOldFile, floorNewFile)
				err := cmd.Run()
				if err != nil {
					glog.Infoln(err)
				}
			} else {
				cmd := exec.Command("cp", floorModelOldFile, floorNewFile)
				err := cmd.Run()
				if err != nil {
					glog.Infoln(err)
				}
			}
		}
	}
}

func setInfoToConfig() {
	initYamlClient()
	setBumbleConfig()
	setFlowserviceConfig()
}

func setFlowserviceConfig() {
	// event server url
	setKeyValue(flowserviceYamlClient, eventUrlYamlKey, sensorInfos[eventUrlKey])
	// basic
	pubDb := filepath.Join(preBasicKey, sensorUid, "pub_db_url")
	setKeyValue(flowserviceYamlClient, pubDb, pubDbUrl)
	pubVibo := filepath.Join(preBasicKey, sensorUid, "pub_vibo_url")
	setKeyValue(flowserviceYamlClient, pubVibo, pubViboUrl)
	// event rule
	for key, value := range ruleInfos {
		var data EventRule
		err := json.Unmarshal(value.([]byte), &data)
		if err != nil {
			glog.Infoln(err)
		}
		ret, err := json.Marshal(data)
		if err != nil {
			glog.Infoln(err)
		}
		setEventRuleInfo(key, string(ret))
		// setEventRuleInfo(key, string(value.([]byte)))
	}
}

func setEventRuleInfo(key, value string) {
	key = strings.Trim(key, "/")
	keys := strings.Split(key, "/")
	id := keys[len(keys)-1]
	ruleYamlKey := filepath.Join(preRuleKey, sensorUid, id)
	glog.Infof("%s: %s", ruleYamlKey, value)

	setKeyValue(flowserviceYamlClient, ruleYamlKey, value)
}

func setBumbleConfig() {
	// uid
	setKeyValue(bumbleYamlClient, sensorUidYamlKey, sensorInfos[sensorUidKey])
	// sn
	setKeyValue(bumbleYamlClient, sensorSnYamlKey, sensorInfos[sensorSnKey])
	// desc
	setKeyValue(bumbleYamlClient, sensorDescYamlKey, sensorInfos[sensorDescKey])
	// ntpaddr
	setKeyValue(bumbleYamlClient, ntpAddrYamlKey, sensorInfos[ntpAddrKey])
	// iot server
	setKeyValue(bumbleYamlClient, iotUrlYamlKey, sensorInfos[iotUrlKey])
}

func setKeyValue(cli *yaml.YamlConfig, key, value string) {
	if "" == value || "" == key {
		glog.Infof("key:%s; value:%s", key, value)
		return
	}
	_, err := cli.SetValue(key, value)
	if err != nil {
		glog.Infoln(err)
	}
}

func initYamlClient() {
	var err error
	bumbleYamlFilePath := filepath.Join(configPath, bumbleFile)
	flowserviceYamlFilePath := filepath.Join(configPath, flowserviceFile)
	bumbleYamlClient, err = yaml.NewYamlConfig(bumbleYamlFilePath)
	if err != nil {
		glog.Infoln(err)
	}
	flowserviceYamlClient, err = yaml.NewYamlConfig(flowserviceYamlFilePath)
	if err != nil {
		glog.Infoln(err)
	}
}

func GetInfoFromCmd() {
	UseCmd(sensorUidKey)
	UseCmd(eventUrlKey)
	UseCmd(sensorSnKey)
	UseCmd(ntpAddrKey)
	UseCmd(sensorDescKey)
	UseCmd(iotUrlKey)
	sensorUid = sensorInfos[sensorUidKey]
	glog.Infoln(sensorInfos)
}

func UseCmd(key string) {
	cmd := exec.Command("etcdctl", "get", key)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		glog.Infoln(err)
		return
	} else if "" == out.String() {
		glog.Infoln(key, ": no data")
		return
	}
	sensorInfos[key] = strings.TrimSpace(out.String())
}

func GetEventRule() {
	ruleDirKey := "http://127.0.0.1:4001/v2/keys/config/eventbrain/alertrule/" +
		sensorInfos[sensorUidKey] + "?recursive=true"
	cmd := exec.Command("curl", ruleDirKey)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		glog.Infoln(err)
		// return
	}
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(out.String()), &data)
	if err != nil {
		glog.Infoln(err)
	}
	nodes := data["node"].(map[string]interface{})["nodes"].([]interface{})
	for _, node := range nodes {
		realData := node.(map[string]interface{})
		if realData["dir"] == nil {
			ruleInfos[realData["key"].(string)] = []byte(realData["value"].(string))
		}
	}
	glog.Infoln(ruleInfos)
}
