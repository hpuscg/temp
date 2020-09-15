/*
#Time      :  2020/9/9 2:22 下午
#Author    :  chuangangshen@deepglint.com
#File      :  rmlibraconfig.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"bytes"
	"errors"
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
	sensorIp string
	port int
	sensorList string
)

func main() {
	flag.StringVar(&sensorList, "sensorList", "ipList.txt", "sensor ip list file name")
	flag.IntVar(&port, "port", 22, "ssh port")
	flag.Parse()
	MvOldLog()
	// 初始化glog配置
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./logs"))
	if sensorList == "" {
		glog.Infoln("Please run program with sensor ip list file name !")
		return
	}
	runApp()
}

// 读取文件中的IP并进行事件规则配置
func runApp() {
	fi, err := os.Open(sensorList)
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
			glog.Infoln(err)
			continue
		}
		rmLibraLog()
	}
}

func rmLibraLog() {
	rmString := "rm /home/deepglint/AppData/libraT/config/libra_v3.json"
	RunCmd(rmString)
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
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if "windows" == sysInfo {
		for {
			line, err := out.ReadString('\n')
			if err != nil {
				break
			}
			if len(line) == 0 {
				continue
			}
			pos := strings.Index(line, "Reply")
			if pos >= 0 {
				pos = strings.Index(line, ":")
				if pos > 0 {
					if strings.Index(line, "TTL=") > 0 { // OK
						return nil
					} else {
						res := line[pos+1:]
						res = strings.Trim(res, " \n")
						return errors.New(res)
					}
				}
			}
		}
	} else {
		if err == nil {
			return err
		} else {
			for {
				line, err := out.ReadString('\n')
				if err != nil {
					break
				}
				if len(line) == 0 {
					continue
				}
				pos := strings.Index(line, "icmp_seq=1")
				if pos >= 0 {
					res := line[pos+len("icmp_seq=1"):]
					res = strings.Trim(res, " \n")
					return errors.New(res)
				}
			}
		}
	}
	return err
}

// 判断是否已有历史log，如有进行移动
func MvOldLog() {
	timeStamp := time.Now().Unix()
	stringTimeStamp := strconv.Itoa(int(timeStamp))
	newLogFileName := filepath.Join("./logs", stringTimeStamp+".log")
	_, fileName := filepath.Split(os.Args[0])
	oldLogFileName := filepath.Join("./logs", fileName+".log")
	_, err := os.Stat(oldLogFileName)
	if err == nil {
		cmd := exec.Command("mv", oldLogFileName, newLogFileName)
		_ = cmd.Run()
	}
}

// 运行命令
func RunCmd(cmd string) {
	session, err := connect("root", sensorIp, port)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", sensorIp, err)
		glog.Infoln(info)
	}
	_, err = session.CombinedOutput(cmd)
	if err != nil {
		glog.Infof("fail for (%s) \n", err)
	}
	_ = session.Close()
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
