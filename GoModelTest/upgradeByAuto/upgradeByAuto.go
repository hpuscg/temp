/*
#Time      :  2020/1/14 2:59 PM 
#Author    :  chuangangshen@deepglint.com
#File      :  upgradeByAuto.go
#Software  :  GoLand
*/
package main

import (
	"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"time"
	"os/exec"
	"runtime"
	"bufio"
	"os"
	"io"
	httpS "github.com/deepglint/flowservice/util/http"
	"net"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"strconv"
	"path/filepath"
)

const (
	DONE = "DONE"
	ING  = "ING"
	TOBE = "TOBE"
)

var (
	ServerIp       string
	DownLoadStatus string
	UpgradeStatus  string
	IpListFile     string
	Images         = map[string]bool{
		"vibo2vibo":  true,
		"bumble-bee": true,
		"adu":        true,
	}
	CurSensor SensorInfo
	Port      int
)

type SensorInfo struct {
	Ip       string
	Port     int
	Username string
	Logfile  string
	Prompt   bool
}

func main() {
	flag.StringVar(&ServerIp, "serverIp", "192.168.100.235", "server addr")
	flag.StringVar(&IpListFile, "ipListFile", "./ip.txt", "sensor ip list")
	flag.IntVar(&Port, "port", 22, "ssh port")
	flag.Parse()
	// 逐行读取配置文件中的设备IP
	fi, err := os.Open(IpListFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		SensorIp := string(a)
		// 测试IP是否能ping通
		err := tryPing(SensorIp)
		if err != nil {
			fmt.Printf("%s 网络不通，请检查\n", SensorIp)
			continue
		}
		// 初始化sensor信息
		CurSensor.Ip = SensorIp
		CurSensor.Port = Port
		CurSensor.Username = "root"
		// 检查设备的磁盘空间
		dfStr := "df | awk 'NR==2{print}'|awk '{print $5}'"
		out, err := runLiveCommand(dfStr)
		var dfInt = 100
		if err != nil {
			fmt.Println(err)
		} else {
			dfInt, err = strconv.Atoi(strings.Split(out, "%")[0])
			if err != nil {
				fmt.Println(err)
			}
		}
		if dfInt > 85 {
			fmt.Printf("%s 磁盘使用率为:%d%", SensorIp, dfInt)
			continue
		}
		// 测试设备启动脚本是否损坏
		tegraSizeStr := "ls /usr/bin/tegra_init.sh -l |cut -d \" \" -f 5"
		tegraExistStr := "find /usr/bin/ -name tegra_init.sh"
		tegraSize, err := runLiveCommand(tegraSizeStr)
		if err != nil {
			fmt.Println(err)
		}
		tegraExist, err := runLiveCommand(tegraExistStr)
		if err != nil {
			fmt.Println(err)
		}
		tegraSizeInt, err := strconv.Atoi(strings.Trim(tegraSize, "\n"))
		if tegraSizeInt < 1754 || strings.Trim(tegraExist, "\n") == "" {
			fmt.Printf("%s 设备的启动脚本损坏，请注意！", SensorIp)
			// 将data下的启动脚本放置到/usr/bin下
			mvTegraStr := "cp /data/shell/_usrbin/tegra_init.sh /usr/bin/"
			runLiveCommand(mvTegraStr)
		}
		// 检查设备是否配置了网管服务器
		serverAddrStr := "etcdctl get /config/global/server_addr"
		serverAddr, err := runLiveCommand(serverAddrStr)
		if strings.Trim(serverAddr, "\n") == "" {
			// 未配置网管服务器
			NoServerAddr(SensorIp)
		} else {
			// 已配置网管服务器
			fmt.Println(strings.Trim(serverAddr, "\n"))
			// 检查设备升级脚本，是否需要预升级
			IsPreUpgrade()
			// 已配置网管服务器，升级设备
			HaveServerAddr(SensorIp)
		}
	}
}

// 设备预升级
func PreUpgrade() {
	copySession, err := copyConnect("root", CurSensor.Ip, 22)
	if err != nil {
		fmt.Println("connect sftp fail, ", err)
		return
	}
	// copy switch package
	err = copyFile(copySession, "/data/shell/service/hookshell/", "./switch_package.sh")
	if err != nil {
		fmt.Println("Copy switch file to sensor fail :", err)
		return
	}
	// copy download package
	err = copyFile(copySession, "/data/shell/service/hookshell/", "./download_package.sh")
	if err != nil {
		fmt.Println("Copy download file to sensor fail :", err)
		return
	}
	// copy load utils
	err = copyFile(copySession, "/tmp", "./load_libraT_utils.sh")
	if err != nil {
		fmt.Println("Copy utils file to sensor fail :", err)
		return
	}
	// copy change load
	err = copyFile(copySession, "/tmp", "./change_reload.sh")
	if err != nil {
		fmt.Println("Copy change file to sensor fail :", err)
		return
	}
	// 执行change_reload.sh
	execChangeStr := "bash /tmp/change_reload.sh"
	runLiveCommand(execChangeStr)
	// 删除change_reload.sh
	rmChangeStr := "rm /tmp/change_reload.sh"
	runLiveCommand(rmChangeStr)
}

// 判读设备是否需要预升级
func IsPreUpgrade() {
	switchPackageSizeStr := "ls /data/shell/service/hookshell/switch_package.sh -l |cut -d \" \" -f 5"
	downLoadPackageSizeStr := "ls /data/shell/service/hookshell/download_package.sh -l |cut -d \" \" -f 5"
	switchPackageSize, err := runLiveCommand(switchPackageSizeStr)
	if err != nil {
		fmt.Println(err)
	}
	downLoadPackageSize, err := runLiveCommand(downLoadPackageSizeStr)
	if err != nil {
		fmt.Println(err)
	}
	switchPackageSizeInt, err := strconv.Atoi(strings.Trim(switchPackageSize, "\n"))
	if err != nil {
		fmt.Println(err)
	}
	downLoadPackageSizeInt, err := strconv.Atoi(strings.Trim(downLoadPackageSize, "\n"))
	if err != nil {
		fmt.Println(err)
	}
	if switchPackageSizeInt != 2220 || downLoadPackageSizeInt != 222 {
		// 设备预升级
		PreUpgrade()
	}
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

// 没有配置网管服务器
func NoServerAddr(ip string) {
	ok := CheckTimeOfSensor(ip)
	if ok {
		SetServerAddr(ip)
	}
}

// 设置NTP较时的结构
type DataSource struct {
	Mode int
	Ntp  string
}

// 给设备较时
func CheckTimeOfSensor(ip string) (ok bool) {
	url := "http://" + ip + ":8008/api/synctime"
	data := DataSource{
		Mode: 1,
		Ntp:  ServerIp,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	_, code, err := httpS.HTTPPost(url, string(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	if code == 200 {
		ok = true
	}
	return
}

// 配置网管服务器
func SetServerAddr(ip string) {
	var data = make(map[string]string)
	data["server_address"] = ServerIp
	url := "http://" + ip + ":8008/api/server_address"
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	_, _, err = httpS.HTTPPost(url, string(jsonData))
	if err != nil {
		fmt.Println(err)
	}
}

// 已经配置网管服务器
func HaveServerAddr(ip string) {
	isOk := CheckServerOfSensor(ip)
	if isOk {
		GetUpgradeStatus(ip)
	} else {
		CheckTimeOfSensor(ip)
	}
}

// docker容器的结构
type Container struct {
	Id      string
	Image   string
	Command string
	Created int64
	Status  string
	Names   []string
}

// 查看设备的vibo2vibo、bumble和adu服务
func CheckServerOfSensor(ip string) bool {
	var isOk = true
	url := "http://" + ip + ":8008/api/container/list"
	result, err := http.Get(url)
	if err != nil {
		isOk = false
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		isOk = false
	}
	var data []Container
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	for _, ret := range data {
		name := strings.Trim(ret.Names[0], "/")
		if strings.Contains(ret.Status, "Up") {
			Images[name] = false
		} else if strings.Contains(ret.Status, "Exited") {
			fmt.Printf("%s : 停止运行\n", name)
		} else {
			fmt.Printf("%s : 其他情况 : %s\n", name, ret.Status)
		}
	}
	for _, value := range Images {
		if value {
			isOk = false
		}
	}
	return isOk
}

// 设备状态结构
type HostModel struct {
	Id            string `json:"id"`
	HostIp        string `json:"hostip"`
	Mac           string `json:"mac"`
	SensorId      string `json:"sensorid"`
	SN            string `json:"sn"`
	SubTopic      string `json:"subtopic"`
	Camera        string `json:"camera"`
	DevModel      string `json:"devmodel"`
	Chip          string `json:"model"`
	Configured    bool   `json:"configured"`
	Desc          string `json:"desc"` // used for etcd tree
	Version       string `json:"version"`
	DownloadState string `json:"downloadstate"`
	UpgradeState  string `json:"upgradestate"`
}

// 设备全部运行状态
type Sensor struct {
	Host         HostModel
	LsReportTime time.Time
	IsInControl  bool
	Status       bool
}

// 获取设备的升级状态，并根据设备状态决定是否进行升级包下载和升级
func GetUpgradeStatus(ip string) {
	url := "http://" + ip + ":8008/api/sensorlist"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data []Sensor
	json.Unmarshal(body, &data)
	sensor := data[0]
	DownLoadStatus = sensor.Host.DownloadState
	UpgradeStatus = sensor.Host.UpgradeState
	if DownLoadStatus == TOBE {
		GetDownloadBegin(ip)
	} else if DownLoadStatus == DONE && UpgradeStatus == TOBE {
		GetUpgradeBegin(ip)
	}
}

// 开始下载镜像
func GetDownloadBegin(ip string) {
	url := "http://" + ip + ":8008/api/download/package"
	http.Get(url)
}

// 开始升级设备
func GetUpgradeBegin(ip string) {
	url := "http://" + ip + ":8008/api/switch/package"
	http.Get(url)
}

// connect for exec cmd
func connect(user, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	testPrivateKeys, err := ssh.ParseRawPrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApnIg4Q/g2thAR3vAUw6EPjqgWIEJ7+FZ+AQZtHUc7b920VJI
7JPmZ1xwdUArlCpQIMAt6mAwV9Z/C+Nc9qIpIeQwKaAd6YWVdr3jFyHLC9rqIT2g
VifCDnMkSnV7Lvuu5QTvgURGOYpyUhUDJBdBY4YAu9q1ITy35oB0xLh1vUCwuDxI
oM5lMc+HsPjf4/SyfyAacBuoD7BvAJsxJ6xuXBaIlmWcw8o76O/Y5PGcYKPS9/bI
rN8TrstuWILp2Nvi4WoxVMIQ98i1S6jM47arI+vNGlFrwolrCanH8GBj1NOBh4BF
JwJisi0Z3+RrtxOVRtgZ9S/tKdK73X6EpbN4hwIDAQABAoIBAAuBRAiKgm5eGENY
qHiVPkrW3pJ/iOJN31wnXGd+2NsOKvZZC7Vem8R1PUi9gMWjDxrUbdgPggfwSaPW
uWxK1TEEhte5u5eSpjwo7/N/YHuXTCu0CMsrwFwjVVTYPgWHXBV0e+GhiIEdsr09
upPaD6kDcDWL7o03lzaVlnyqi2jjXT6kUDyEFCbIAGtoxaYf3clT5e30FnyZhiCH
m8/Qqv5M1wcVIVdsItHqMsQXQF34eT/Lg3r/Ui1bQcUldc6yYjGpC08EdDNKhGT2
f2QwAv7UJ+GB8RNl12w3fAh3ReuiW8NEtDQ1nuSahkX5YlIWkqRDOd6Sjrg1ZkfW
u0/zPZECgYEA2m+w90vb3ui7M/Q0AYJivo88YKhT3ismQs2+CkkgWJ7IohJj3VSh
REljeAwEVEKv8G8lXgjTNKQ+B4sPFckIvIWGkwo7cuerIwn9n41K20oGb6gEl0jW
mVbhv0dy6yfp8deBCOZB4YgonXWsuv4lw8DaUoakGxZgFfChjH0VvbUCgYEAwxGj
rmq+RQWYYna9WWn2GPEJoX0SBU39pHQYBKfQ++pMIUrrryCjPvBNnIICng82RjTp
MU8BvudvDCJgj3J79TDetBnwVt8/nAGIkleyuWzDMQwF7khBS9/TqUUqmH88GmOt
40BPThCBx8YgKiPpmGYgPnUww1bqpvxKT9O0IssCgYEAjFH7qKD+mW9/8pwJXH7Z
1/hDnQQE/E9TwM5SKmFXehZmZFbT+DaJckiCsXdmwIomY5nCs2mP490uS8I06pW+
GvzbulF0ZxgTg+rDFl+5mq0u/UM9z8FmuhJp6mqHlDCLxGPf7EuePrctABm74FOr
Btk4ZpM/kHcLOozd+lXQRZECgYBipWr26zgpQ3kaYh3DN9iiKFLMfak9UYFxRtxW
jl8a5hN1yqOBPqoPTAqTmROlxt+VhXBf5Spm1jbMFh5qrGSPTBVzUqK968wJIqVk
DEFvj9bt2LyvEY8jxZ8OPNIbqExGtB3djEoOmj5nPoRJizu4O/0WWME+J5gmtfMG
h3LTHQKBgDlITGqdIM4Pp54X5ppOW9S55yaAMBJUUhgUsJ73vEcQsBCZ8xkJXg/Q
muPfcFzSD/IgeFoWxYrJIk0CBov3ah+14z5YV1JoKIXAlL7V18f7Omaav8/bozOP
x78MQ06CGEFRcD4LPMITxTDj6zDm1h7iPhG4m2c9Shy0rwpFmFdd
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)

	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		//		Timeout: 			60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

// connect for copy file
func copyConnect(user, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	testPrivateKeys, err := ssh.ParseRawPrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApnIg4Q/g2thAR3vAUw6EPjqgWIEJ7+FZ+AQZtHUc7b920VJI
7JPmZ1xwdUArlCpQIMAt6mAwV9Z/C+Nc9qIpIeQwKaAd6YWVdr3jFyHLC9rqIT2g
VifCDnMkSnV7Lvuu5QTvgURGOYpyUhUDJBdBY4YAu9q1ITy35oB0xLh1vUCwuDxI
oM5lMc+HsPjf4/SyfyAacBuoD7BvAJsxJ6xuXBaIlmWcw8o76O/Y5PGcYKPS9/bI
rN8TrstuWILp2Nvi4WoxVMIQ98i1S6jM47arI+vNGlFrwolrCanH8GBj1NOBh4BF
JwJisi0Z3+RrtxOVRtgZ9S/tKdK73X6EpbN4hwIDAQABAoIBAAuBRAiKgm5eGENY
qHiVPkrW3pJ/iOJN31wnXGd+2NsOKvZZC7Vem8R1PUi9gMWjDxrUbdgPggfwSaPW
uWxK1TEEhte5u5eSpjwo7/N/YHuXTCu0CMsrwFwjVVTYPgWHXBV0e+GhiIEdsr09
upPaD6kDcDWL7o03lzaVlnyqi2jjXT6kUDyEFCbIAGtoxaYf3clT5e30FnyZhiCH
m8/Qqv5M1wcVIVdsItHqMsQXQF34eT/Lg3r/Ui1bQcUldc6yYjGpC08EdDNKhGT2
f2QwAv7UJ+GB8RNl12w3fAh3ReuiW8NEtDQ1nuSahkX5YlIWkqRDOd6Sjrg1ZkfW
u0/zPZECgYEA2m+w90vb3ui7M/Q0AYJivo88YKhT3ismQs2+CkkgWJ7IohJj3VSh
REljeAwEVEKv8G8lXgjTNKQ+B4sPFckIvIWGkwo7cuerIwn9n41K20oGb6gEl0jW
mVbhv0dy6yfp8deBCOZB4YgonXWsuv4lw8DaUoakGxZgFfChjH0VvbUCgYEAwxGj
rmq+RQWYYna9WWn2GPEJoX0SBU39pHQYBKfQ++pMIUrrryCjPvBNnIICng82RjTp
MU8BvudvDCJgj3J79TDetBnwVt8/nAGIkleyuWzDMQwF7khBS9/TqUUqmH88GmOt
40BPThCBx8YgKiPpmGYgPnUww1bqpvxKT9O0IssCgYEAjFH7qKD+mW9/8pwJXH7Z
1/hDnQQE/E9TwM5SKmFXehZmZFbT+DaJckiCsXdmwIomY5nCs2mP490uS8I06pW+
GvzbulF0ZxgTg+rDFl+5mq0u/UM9z8FmuhJp6mqHlDCLxGPf7EuePrctABm74FOr
Btk4ZpM/kHcLOozd+lXQRZECgYBipWr26zgpQ3kaYh3DN9iiKFLMfak9UYFxRtxW
jl8a5hN1yqOBPqoPTAqTmROlxt+VhXBf5Spm1jbMFh5qrGSPTBVzUqK968wJIqVk
DEFvj9bt2LyvEY8jxZ8OPNIbqExGtB3djEoOmj5nPoRJizu4O/0WWME+J5gmtfMG
h3LTHQKBgDlITGqdIM4Pp54X5ppOW9S55yaAMBJUUhgUsJ73vEcQsBCZ8xkJXg/Q
muPfcFzSD/IgeFoWxYrJIk0CBov3ah+14z5YV1JoKIXAlL7V18f7Omaav8/bozOP
x78MQ06CGEFRcD4LPMITxTDj6zDm1h7iPhG4m2c9Shy0rwpFmFdd
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)
	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		//		Timeout: 			60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// rm docker image which need upgrade
func runLiveCommand(cmd string) (log string, err error) {
	session, err := connect(CurSensor.Username, CurSensor.Ip, CurSensor.Port)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", CurSensor.Ip, err)
		fmt.Println(info)
		return "", err
	}
	if CurSensor.Prompt {
		fmt.Println("root:~# ", cmd)
	}
	buf, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Printf("fail for (%s) \n", err)
		return "", err
	}
	log = string(buf)
	session.Close()
	return log, nil
}

// copy file
func copyFile(copySession *sftp.Client, destpath string, file string) error {
	srcFile, err := os.Open(file)
	if err != nil {
		info := fmt.Sprintf("open %s failed :%s", file, err)
		fmt.Println(info)
		return err
	}
	df := filepath.Base(file)
	dstFile, err := copySession.Create(destpath + df)
	if err != nil {
		info := fmt.Sprintf("create %s failed :%s", file, err)
		fmt.Println(info)
		return err
	}
	// defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("readall err: ", err)
		return err
	}
	dstFile.Write(ff)
	fmt.Printf("copy file (%s) to (%s) successful\n", file, destpath+df)
	return nil
}
