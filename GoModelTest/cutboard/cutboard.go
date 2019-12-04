/*
#Time      :  2019/1/15 上午11:33 
#Author    :  chuangangshen@deepglint.com
#File      :  cutboard.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"strings"
	"os"
	"fmt"
	"io"
	"net"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
	"io/ioutil"
)

var Ips []string

func main() {
	readConfig()
	for _, Ip := range Ips {
		handlers(Ip)
	}
}

// read libra ip
func readConfig() {
	f, err := os.Open("config.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		libraIp := strings.TrimSpace(string(b))
		if 0 != len(libraIp) {
			fmt.Println(libraIp)
			Ips = append(Ips, libraIp)
		}
	}
	f.Close()
}

// connect for exec cmd
func connect(user, host string, port int) (*ssh.Client, error) {
	var (
		auth 			[]ssh.AuthMethod
		addr 			string
		clientConfig 	*ssh.ClientConfig
		client 			*ssh.Client
		err 			error
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
		// fmt.Println("Unable to parse test key :", err)
		return nil, err
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)

	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: 				user,
		Auth: 				auth,
		Timeout: 			30 * time.Second,
		HostKeyCallback: 	func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return client, nil
}

// copy client
func copyConnect(user, host string, port int) (*sftp.Client, error) {
	var (
		auth 			[]ssh.AuthMethod
		addr 			string
		clientConfig 	*ssh.ClientConfig
		sshClient 		*ssh.Client
		sftpClient 		*sftp.Client
		err 			error
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
		// fmt.Println("Unable to parse test key :", err)
		return nil, err
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)
	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: 				user,
		Auth: 				auth,
		Timeout: 			30 * time.Second,
		HostKeyCallback: 	func(hostname string, remote net.Addr, key ssh.PublicKey) error {
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

// copy image file, copy reload file
func copyFile(copySession *sftp.Client, client *ssh.Client) {
	fmt.Println("start copy file")
	var session *ssh.Session
	// reset config
	session, _ = client.NewSession()
	session.Run(`etcdctl set /config/global/sensor_uid ""`)
	session, _ = client.NewSession()
	session.Run(`etcdctl set /config/global/sensor_sn ""`)
	fmt.Println("copy image file")
	copyImageFiles := []string{
		"bumble-bee_1.8.0.tar.gz", "eventserver_1.7.2.tar.gz", "libra-cuda_1.4.34.hf_20180329.tar.gz",
		"pioneer_1.5.4_hf.tar.gz", "flowservice_2.2.5.5.tar.gz",
	}
	for _, fileName := range copyImageFiles {
		fmt.Println("image file name is: ", fileName)
		srcFile, err := os.Open(fileName)
		if err != nil {
			info := fmt.Sprintf("open %s failed :%s", fileName, err)
			fmt.Println(info)
			return
		}
		defer srcFile.Close()
		dstFile, err := copySession.Create(`/home/ubuntu/` + fileName)
		if err != nil {
			info := fmt.Sprintf("create %s failed :%s", fileName, err)
			fmt.Println(info)
			return
		}
		defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			fmt.Println("readall err: ", err)
			return
		}
		dstFile.Write(ff)
		session, _ = client.NewSession()
		session.Run(`docker load --input ` + `/home/ubuntu/` + fileName)
		time.Sleep(10 * time.Second)
		session, _ = client.NewSession()
		session.Run(`rm ` + `/home/ubuntu/` + fileName)
	}
	// copy file
	copyFiles := []map[string]string {{"oldFile":"reload_sensor.sh", "newFile":"/usr/bin/reload_sensor.sh"}}

	// copy config file
	for _, files := range copyFiles {
		fmt.Println("copy file name is: ", files)
		srcFile, err := os.Open(files["oldFile"])
		if err != nil {
			info := fmt.Sprintf("open %s failed :%s", files["oldFile"], err)
			fmt.Println(info)
			return
		}
		defer srcFile.Close()
		session, _ = client.NewSession()
		session.Run(`rm ` + files["newFile"] + ``)
		dstFile, err := copySession.Create(`` + files["newFile"] + ``)
		if err != nil {
			info := fmt.Sprintf("create %s failed :%s", files["newFile"], err)
			fmt.Println(info)
			return
		}
		defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			fmt.Println("readall err: ", err)
			return
		}
		dstFile.Write(ff)
	}
	time.Sleep(10 * time.Second)
}

// handler cutboard
func handlers(Ip string) {
	// try get ssh client
	client, err := connect("root", Ip, 22)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", Ip, err)
		fmt.Println(info)
		return
	}
	defer client.Close()
	// try get copy client
	copySession, err := copyConnect("root", Ip, 22)
	if err != nil {
		fmt.Println("connect sftp fail, ", err)
		return
	}
	defer copySession.Close()

	// copyFiles
	copyFile(copySession, client)
	// 更新服务
	var session *ssh.Session
	session, _ = client.NewSession()
	session.Run(`chmod +x /usr/bin/reload_sensor.sh`)
	session, _ = client.NewSession()
	session.Run(`reload_sensor.sh`)
	fmt.Println(Ip, "is over!")
}



