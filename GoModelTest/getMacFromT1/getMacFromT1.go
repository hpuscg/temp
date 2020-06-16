/*
#Time      :  2020/5/28 5:41 下午
#Author    :  chuangangshen@deepglint.com
#File      :  getMacFromT1.go
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
)

var (
	sensorIp string
	ipFile string
	port int
	logDir string
)

func main() {
	flag.StringVar(&sensorIp, "sensorIp", "", "sensor ip")
	flag.StringVar(&ipFile, "ipFile", "ip.txt", "sensor ip file name")
	flag.IntVar(&port, "port", 22, "ssh port")
	flag.StringVar(&logDir, "logDir", "./", "log dir")
	flag.Parse()
	// 删除已有log
	_, fileName := filepath.Split(os.Args[0])
	LogFileName := filepath.Join(logDir, fileName+".log")
	_, err := os.Stat(LogFileName)
	if err == nil {
		cmd := exec.Command("rm ", LogFileName)
		_ = cmd.Run()
	}
	// 初始化glog配置
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath(logDir))
	// 逐行读取配置文件中的设备IP
	fi, err := os.Open(ipFile)
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
		sensorIp = string(a)
		glog.Infoln(sensorIp)
		// 测试IP是否能ping通
		err := tryPing(sensorIp)
		if err != nil {
			glog.Infof("%s 网络不通，请检查\n", sensorIp)
			continue
		}
		// 获取设备Mac信息
		getMacStr := "ifconfig"
		ret, err := runLiveCommand(getMacStr)
		if err != nil {
			glog.Infoln(err)
		}
		glog.Infoln(ret)
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

// run cmd
func runLiveCommand(cmd string) (log string, err error) {
	session, err := connect("root", sensorIp, port)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", sensorIp, err)
		glog.Infoln(info)
		return "", err
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


