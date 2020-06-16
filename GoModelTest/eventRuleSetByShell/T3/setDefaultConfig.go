/*
#Time      :  2020/6/9 4:27 下午
#Author    :  chuangangshen@deepglint.com
#File      :  setDefaultConfig.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"io/ioutil"
	"net/http"
	"strings"
	"temp/GoModelTest/eventRuleSetByShell/T3/yaml"
)

var (
	Token        string
	UserName     = "admin"
	PassWord     = "Dg1304!@"
	AsePassWord  string
	aeskey       = []byte("t3.deepglint.com")
	aesiv        = []byte("86-10-62950616-9")
	yamlFileName = "config.yaml"
	sensorIp     string
	yamlClient   *yaml.YamlConfig = nil
)

const ()

type Response struct {
	Code     int         `json:"code"`
	Msg      string      `json:"msg"`
	Redirect string      `json:"redirect"`
	Data     interface{} `json:"data"`
}

func main() {
	flag.StringVar(&sensorIp, "sensorIp", "192.168.100.92", "sensor ip")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("logs"), glog.WithLevel("info"))
	run()
}

func run() {
	GetToken()
}

// 初始化yaml client
func initYamlClient() (err error) {
	yamlClient, err = yaml.NewYamlConfig(yamlFileName)
	if err != nil {
		glog.Infoln(err)
	}
	return
}

// 获取设备登陆token
func GetToken() {
	tokenUrl := "http://" + sensorIp + "/api/a/login"
	EncodePassWord()
	body, err := GetTokenByPostWithPassWord(tokenUrl, UserName, AsePassWord)
	if err != nil {
		glog.Infoln(err)
		return
	}
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		glog.Infoln(err)
	}
	Token = resp.Data.(map[string]interface{})["token"].(string)
}

func HttpPut() {
	//TODO
}

func HttpGet() {
	//TODO
}

func HttpPost() {
	//TODO
}

// 通过用户名密码获取token
func GetTokenByPostWithPassWord(url, username, password string) (body []byte, err error) {
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
	body = respData
	return
}

// 密码加密
// 加密
func EncodePassWord() {
	ciphertext := AESEncrypt([]byte(PassWord), aeskey, aesiv)
	ret := hex.EncodeToString(ciphertext)
	AsePassWord = strings.ToUpper(ret)
}

//加密
func AESEncrypt(origData, key, iv []byte) []byte {
	//获取block块
	block, _ := aes.NewCipher(key)
	//补码
	origData = PKCS7Padding(origData, block.BlockSize())

	//创建明文长度的数组
	crypted := make([]byte, len(origData))
	//加密模式，
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密明文
	blockMode.CryptBlocks(crypted, origData)

	return crypted
}

//补码
func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}
