/*
#Time      :  2020/4/16 2:40 下午
#Author    :  chuangangshen@deepglint.com
#File      :  rmLogs.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"gitlab.deepglint.com/junkaicao/glog"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	IpListFile string
	Port       int
	CurSensor  SensorInfo
	logDir     string
	alsoToFile bool
	logLevel   string
)

type SensorInfo struct {
	Ip       string
	Port     int
	Username string
	Logfile  string
	Prompt   bool
}

func main() {
	flag.StringVar(&IpListFile, "ipListFile", "./ip.txt", "sensor ip list")
	flag.StringVar(&logDir, "log_dir", "logs", "log dir, default /tmp")
	flag.BoolVar(&alsoToFile, "alsologtostderr", false, "log to stderr also to log file")
	flag.StringVar(&logLevel, "log_level", "info", "log level, default info")
	flag.IntVar(&Port, "port", 22, "ssh port")
	flag.Parse()
	// 判断是否已有历史log，如有进行移动
	timeStamp := time.Now().Unix()
	stringTimeStamp := strconv.Itoa(int(timeStamp))
	newLogFileName := filepath.Join(logDir, stringTimeStamp+".log")
	fileNameArr := strings.Split(os.Args[0], "/")
	oldLogFileName := filepath.Join(logDir, fileNameArr[len(fileNameArr)-1]+".log")
	_, err := os.Stat(oldLogFileName)
	if err == nil {
		cmd := exec.Command("mv", oldLogFileName, newLogFileName)
		_ = cmd.Run()
	}
	// 初始化glog配置
	glog.Config(glog.WithAlsoToStd(alsoToFile), glog.WithFilePath(logDir))
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
		glog.Infof("============%d===========", i)
		SensorIp := string(a)
		glog.Infoln(SensorIp)
		// 测试IP是否能ping通
		err := tryPing(SensorIp)
		if err != nil {
			glog.Infof("%s 网络不通，请检查\n", SensorIp)
			continue
		}
		// 初始化sensor信息
		CurSensor.Ip = SensorIp
		CurSensor.Port = Port
		CurSensor.Username = "root"
		// 删除root下的tar文件
		rmTar := "rm -rf /root/*.tar.gz'"
		rmTarOut, err := runLiveCommand(rmTar)
		glog.Infoln(rmTarOut)
		// 删除上传图片脚本
		rmPicture := "rm /usr/bin/postpicture"
		rmPictureOut, err := runLiveCommand(rmPicture)
		glog.Infoln(rmPictureOut)
		// 删除docker log
		rmDockerLog := "cd /var/lib/docker/containers && rm -rf `find ./ -name *json.log`"
		rmDockerLogOut, err := runLiveCommand(rmDockerLog)
		glog.Infoln(rmDockerLogOut)
		// 删除libra log
		rmLibraLog := "cd /libra/logs &&  rm -rf Libra.muxerlab.tegra-ubuntu*"
		rmLibraLogOut, err := runLiveCommand(rmLibraLog)
		glog.Infoln(rmLibraLogOut)
		// 删除bumble log
		rmBumbleLog := "cd /var/lib/docker/aufs/diff &&  rm -rf `find ./ -name *T00:00:00Z`"
		rmBumbleLogOut, err := runLiveCommand(rmBumbleLog)
		glog.Infoln(rmBumbleLogOut)
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
		glog.Infoln("Unable to parse test key :", err)
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

func runLiveCommand(cmd string) (log string, err error) {
	session, err := connect(CurSensor.Username, CurSensor.Ip, CurSensor.Port)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", CurSensor.Ip, err)
		glog.Infoln(info)
		return "", err
	}
	if CurSensor.Prompt {
		glog.Infoln("root:~# ", cmd)
	}
	buf, err := session.CombinedOutput(cmd)
	if err != nil {
		glog.Infof("fail for (%s) \n", err)
		return "", err
	}
	log = string(buf)
	_ = session.Close()
	return log, nil
}
