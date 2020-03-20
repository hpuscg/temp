/*
#Time      :  2020/3/13 9:58 AM
#Author    :  chuangangshen@deepglint.com
#File      :  tansforConfig.go
#Software  :  GoLand
*/

// 该程序的目的是将现场已安装的T1设备内的程序升级为T3后，保留部分基础配置；
// 这里选用的中间配置存储方式为yaml文件；
// IOT和金砖平台登陆信息建议直接通过工具群发；
/*
具体实现步骤如下：
1、将要升级的T1的基础配置导出后存储；
2、设备进行升级，并清理系统垃圾及无关数据；
3、将之前导出的设备数据分别导入到原有的设备中；
*/

package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.deepglint.com/junkaicao/glog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"temp/GoModelTest/transforConfig/yaml"
	"time"
)

const (
	ContentType     = "application/json"
	DefaultUserName = "admin"
	DefaultPassWord = "Dg1304!@"
)

var (
	IsSet          bool
	IpListFile     string
	SensorId       string
	SensorIp       string
	FileName       string
	YamlClient     *yaml.YamlConfig = nil
	LogDir         string
	AlsoToFile     bool
	LogLevel       string
	Token          string
	TokenUrl       = "http://%s/api/a/login"
	T1SensorIdUrl  = "http://%s:8008/api/sensorid"
	T1EventRuleUrl = "http://%s:8008/api/iterate_values?key=/config/eventbrain/alertrule"
	T1DescUrl      = "http://%s:8008/api/name"
	T3SensorIdUrl  = "http://%s/api/sensorid"
	T3DescUrl      = "http://%s/api/name"
	T3EventRuleUrl = "http://%s/api/f/eventrule"
	UserName       string
	PassWord       string
	AsePassWord    string
	AesKey         = []byte("t3.deepglint.com")
	AesIv          = []byte("86-10-62950616-9")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			glog.Fatalf("panic err: %s", err)
			glog.Flush()
		}
	}()
	Run()
}

func Run() {
	flag.BoolVar(&IsSet, "isSet", true, "判断程序启动是获取配置还是设置配置，false为获取，true为设置")
	flag.StringVar(&IpListFile, "configFile", "ipList.txt", "配置文件的路径, default:ipList.txt")
	flag.StringVar(&SensorIp, "sensorIp", "", "配置迁移的设备IP")
	flag.StringVar(&UserName, "username", "", "设备web登陆用户名")
	flag.StringVar(&PassWord, "password", "", "设备web登陆的密码")
	flag.StringVar(&LogDir, "logDir", "log", "log的存放位置")
	flag.BoolVar(&AlsoToFile, "alsologtostderr", true, "log to stderr also to log file, default:true")
	flag.StringVar(&LogLevel, "log_level", "info", "log level, default:info")
	flag.Parse()
	// 初始化设备登陆用户名密码
	if UserName == "" {
		UserName = DefaultUserName
	}
	if PassWord == "" {
		PassWord = DefaultPassWord
	}
	// 判断是否已有历史log，如有进行移动
	timeStamp := time.Now().Unix()
	stringTimeStamp := strconv.Itoa(int(timeStamp))
	newLogFileName := filepath.Join(LogDir, stringTimeStamp+".log")
	fileNameArr := strings.Split(os.Args[0], "/")
	oldLogFileName := filepath.Join(LogDir, fileNameArr[len(fileNameArr)-1]+".log")
	_, err := os.Stat(oldLogFileName)
	if err == nil {
		cmd := exec.Command("mv", oldLogFileName, newLogFileName)
		_ = cmd.Run()
	}
	// init log config
	glog.Config(glog.WithAlsoToStd(AlsoToFile), glog.WithFilePath(LogDir), glog.WithLevel(LogLevel))
	err = EncodePassWord()
	if err != nil {
		glog.Infoln(err)
		return
	}
	if SensorIp == "" {
		// 逐行读取配置文件中的设备IP
		fi, err := os.Open(IpListFile)
		if err != nil {
			glog.Infof("Error: %s\n", err)
			return
		}
		defer fi.Close()

		br := bufio.NewReader(fi)
		for i := 0; i >= 0; i++ {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			SensorIp = string(a)
			glog.Infof("============%d============", i+1)
			err = tryPing(SensorIp)
			if err != nil {
				glog.Infof("%s 网络不通，请检查\n", SensorIp)
				return
			}
			if IsSet {
				SetConfig()
			} else {
				GetConfig()
			}
		}
	} else {
		err = tryPing(SensorIp)
		if err != nil {
			glog.Infof("%s 网络不通，请检查\n", SensorIp)
			return
		}
		if IsSet {
			SetConfig()
		} else {
			GetConfig()
		}
	}
}

func EncodePassWord() error {
	passWordByte := []byte(PassWord)
	//获取block块
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return err
	}
	//补码
	passWordByte = PKCS7Padding(passWordByte, block.BlockSize())

	//创建明文长度的数组
	crypted := make([]byte, len(passWordByte))
	//加密模式，
	blockMode := cipher.NewCBCEncrypter(block, AesIv)
	//加密明文
	blockMode.CryptBlocks(crypted, passWordByte)
	AsePassWord = hex.EncodeToString(crypted)
	return nil
}

func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

// 测试设备IP能否ping通
func tryPing(ip string) error {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if "windows" == sysInfo {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}
	_, err := cmd.Output()
	return err
}

// 事件规则的结构
type EventRule struct {
	Id           string `json:"Id,omitempty"`
	TimeRange    [2]int `json:"TimeRange,omitempty"`
	WeekdayRange byte   `json:"WeekdayRange,omitempty"`
	Enabled      bool   `json:"Enabled,omitempty"`
	UpperBound   float64
	LowerBound   float64
}

// json api
type ResponseData struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

// 升级前从T1设备上获取已有配置
func GetConfig() {
	SensorId = ""
	GetSensorIdFromT1()
	if SensorId == "" {
		glog.Infoln("get sensor id err")
		return
	}
	CreateConfigFile()
	err := InitYamlClient()
	if err != nil {
		glog.Infoln(err)
		return
	}
	GetSensorConfig()
}

// 获取设备的配置信息
// desc 和 event rule
func GetSensorConfig() {
	_, err := YamlClient.SetValue("Ip", SensorIp)
	if err != nil {
		glog.Infoln(err)
	}
	_, err = YamlClient.SetValue("SensorId", SensorId)
	if err != nil {
		glog.Infoln(err)
	}
	url := fmt.Sprintf(T1EventRuleUrl, SensorIp)
	ret, err := http.Get(url)
	if err != nil {
		glog.Infoln(err)
	}
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		glog.Infoln(err)
	}
	value := make(map[string]string)
	val := make(map[string]EventRule)
	_ = json.Unmarshal(body, &value)
	for i, v := range value {
		key := strings.Split(i, "/")[len(strings.Split(i, "/"))-1]
		var eventRule EventRule
		_ = json.Unmarshal([]byte(v), &eventRule)
		val[key] = eventRule
	}
	byteVal, err := json.Marshal(val)
	if err != nil {
		glog.Infoln(err)
	}
	_, err = YamlClient.SetValue("EventRule", string(byteVal))
	if err != nil {
		glog.Infoln(err)
	}
	descUrl := fmt.Sprintf(T1DescUrl, SensorIp)
	descRet, err := http.Get(descUrl)
	if err != nil {
		glog.Infoln(err)
	}
	defer ret.Body.Close()
	descBody, err := ioutil.ReadAll(descRet.Body)
	if err != nil {
		glog.Infoln(err)
	}
	_, err = YamlClient.SetValue("Desc", string(descBody))
	if err != nil {
		glog.Infoln(err)
	}
	glog.Infof("======%s====Successful", SensorIp)
}

// init yaml client
func InitYamlClient() (err error) {
	YamlClient, err = yaml.NewYamlConfig(FileName)
	return
}

// 创建配置文件
func CreateConfigFile() {
	FileName = filepath.Join("config", SensorId+".yaml")
	f, err := os.Create(FileName)
	defer f.Close()
	if err != nil {
		glog.Infoln(err)
	}
}

// 获取设备的SensorId
func GetSensorIdFromT1() {
	url := fmt.Sprintf(T1SensorIdUrl, SensorIp)
	result, err := http.Get(url)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	SensorId = string(body)
}

// 升级后将原有配置重新导入设备
func SetConfig() {
	// 根据用户名密码获取token
	Token = ""
	TokenUrl = fmt.Sprintf(TokenUrl, SensorIp)
	GetTokenByPostWithPassWord(TokenUrl, UserName, AsePassWord)
	if Token == "" {
		glog.Infoln("get token err")
		return
	}
	// 根据IP获取设备sensorID；
	SensorId = ""
	GetSensorIdFromT3()
	if SensorId == "" {
		glog.Infoln("get sensor id err")
		return
	}
	// 根据sensorID读取设备的配置文件，并校验IP；
	FileName = filepath.Join("config", SensorId+".yaml")
	err := InitYamlClient()
	if err != nil {
		glog.Infoln(err)
		return
	}
	ip, err := YamlClient.GetString("Ip")
	if err != nil || ip != SensorIp {
		glog.Infof("err is :%+v, get ip is: %s", err, ip)
		return
	}
	// 通过http根据配置文件将原配置导入设备；
	SetSensorConfig()
}

func SetSensorConfig() {
	desc, err := YamlClient.GetString("Desc")
	if err != nil {
		glog.Infoln("get sensor desc from yaml file err: ", desc)
		return
	}
	descUrl := fmt.Sprintf(T3DescUrl, SensorIp)
	var descData = make(map[string]string)
	descData["sensor_desc"] = desc
	descByte, err := json.Marshal(descData)
	if err != nil {
		glog.Infoln(err)
	}
	descRet, err := PostWithToken(descUrl, string(descByte), Token, ContentType)
	if err != nil {
		glog.Infoln(err)
	}
	var descResp ResponseData
	err = json.Unmarshal(descRet, &descResp)
	if err != nil || descResp.Code != 0 {
		glog.Infoln(err)
		glog.Infoln(descResp.Msg)
		return
	}

	/*// TODO post event rule to sensor
	eventRuleString, err := YamlClient.GetString("EventRule")
	if err != nil {
		glog.Infoln("get event rule from yaml file err: ", eventRuleString)
		return
	}
	eventRuleUrl := fmt.Sprintf(T3EventRuleUrl, SensorIp)
	_, err = PostWithToken(eventRuleUrl, eventRuleString, Token, ContentType)
	if err != nil {
		glog.Infoln(err)
	}
	var eventRuleResp ResponseData
	err = json.Unmarshal(descRet, &eventRuleResp)
	if err != nil || eventRuleResp.Code != 0 {
		glog.Infoln(err)
		glog.Infoln(eventRuleResp.Msg)
		return
	}*/
	glog.Infof("======%s====Successful", SensorIp)
}

func PostWithToken(url, data, token, contentType string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		glog.Infoln(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("authorization", token)
	response, err := client.Do(req)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	return
}

func GetSensorIdFromT3() {
	url := fmt.Sprintf(T3SensorIdUrl, SensorIp)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Infoln(err)
		return
	}
	req.Header.Set("Content-Type", ContentType)
	req.Header.Add("authorization", Token)
	response, err := client.Do(req)
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	var respSensorId ResponseData
	err = json.Unmarshal(body, &respSensorId)
	if err != nil {
		glog.Infoln(err)
	}
	SensorId = respSensorId.Data.(map[string]interface{})["sensor_id"].(string)
}

// 通过用户名密码获取token
func GetTokenByPostWithPassWord(url, username, password string) {
	data := make(map[string]string)
	data["username"] = username
	data["password"] = password
	byteData, err := json.Marshal(data)
	if err != nil {
		glog.Infoln(err)
		return
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(byteData)))
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infoln(err)
		return
	}
	var respToken ResponseData
	err = json.Unmarshal(respData, &respToken)
	if err != nil {
		glog.Infoln(err)
	}
	Token = respToken.Data.(map[string]interface{})["token"].(string)
}
