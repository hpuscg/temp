package main

import (
	"flag"
	"net"
	"golang.org/x/crypto/ssh"
	"time"
	"fmt"
	"github.com/pkg/sftp"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
)

type Config struct {
	Ip 		string
	dir 	string
}


func main() {
	// flag the cmd config
	var Config Config
	// flag.StringVar(&Config.Ip, "ip", "192.168.12.12", "the ip of LibraT")
	flag.StringVar(&Config.dir, "dir", "./", "the dir of auxiliary file")
	flag.Parse()
	// read ip form file
	Config.Ip = readIp()
	// test the ip is right?
	tryPing(Config.Ip)
	// try get ssh client
	client, err := connect("root", Config.Ip, 22)
	if err != nil {
		// info := fmt.Sprintf("connect to %s err : %s", Config.Ip, err)
		// fmt.Println(info)
		return
	}
	defer client.Close()
	// try get copy client
	copySession, err := copyConnect("root", Config.Ip, 22)
	if err != nil {
		// fmt.Println("connect sftp fail, ", err)
		return
	}
	defer copySession.Close()
	// read upgrade file from dir
	loadImage := readImageFiles(Config.dir)
	// copyFiles
	copyFile(copySession, client, loadImage, Config.dir)
	// deleteImage
	deleteImage(client, loadImage)
	// reboot
	session, _ := client.NewSession()
	session.Run(`reboot`)
	time.Sleep(120 * time.Second)
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

// connect for copy file
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

// test the ip is right
func tryPing(IP string)  {
	cmd := exec.Command("ping", IP)
	// cmd := exec.Command("ping", IP) //for windows
	_, err := cmd.Output()
	if err != nil {
		info := fmt.Sprintf("can't ping %s,the err is : %s, please checkout your config ip", IP, err)
		fmt.Println(info)
		os.Exit(1)
	}
	fmt.Println("this ip is ok !")
}

// read upgrade image file from dir
func readImageFiles(listDir string) []map[string]string {
	fmt.Println("start read image file")
	var deleteImages []map[string]string
	allImage := []string {"libra-init", "tunerd", "bumble-bee", "pioneer", "eventserver", "libra-cuda",
		"flowservice", "adu", "nanomsg2nsq", "vulcand", "crtmpserver", "vodserver", "etcd", "nsq", "onvifserver"}
	files, err := ioutil.ReadDir(listDir)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		// fmt.Println(file.Name())
		if strings.HasSuffix(file.Name(), ".tar") {
			fileName := strings.Split(file.Name(), ".")[0]
			fileTag := strings.Split(strings.Split(file.Name(), fileName + ".")[1], ".tar")[0]
			for _, tempImage := range allImage {
				if tempImage == fileName {
					fmt.Println("name", fileName)
					fmt.Println("newtag", fileTag)
					tempDict := make(map[string]string)
					tempDict["name"] = tempImage
					tempDict["NewTag"] = fileTag
					deleteImages = append(deleteImages, tempDict)
					fmt.Println(deleteImages)
				}
			}
		}
	}
	return deleteImages
}

// copy file and load docker image
func copyFile(copySession *sftp.Client, client *ssh.Client, copyImageFiles []map[string]string, dir string) {
	fmt.Println("start copy file")
	var session *ssh.Session
	copyFiles := []map[string]string {{"oldFile":"Release", "newFile":"/home/ubuntu/Release"},
		{"oldFile":"70-persistent-net.rules", "newFile":"/etc/udev/rules.d/70-persistent-net.rules"},
		{"oldFile":"reload_sensor.sh", "newFile":"/data/shell/_usrbin/reload_sensor.sh"}}
	// load docker image
	for _, files := range copyImageFiles {
		fmt.Println("image file is : ", files["name"] + "." + files["NewTag"] + ".tar")
		srcFile, err := os.Open(dir + files["name"] + "." + files["NewTag"] + ".tar")
		if err != nil {
			info := fmt.Sprintf("open %s failed :%s", files["name"], err)
			fmt.Println(info)
			return
		}
		defer srcFile.Close()
		dstFile, err := copySession.Create(`/home/ubuntu/` + files["name"] + ".tar")
		if err != nil {
			info := fmt.Sprintf("create %s failed :%s", files, err)
			fmt.Println(info)
			return
		}
		defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			// fmt.Println("readall err: ", err)
			return
		}
		dstFile.Write(ff)
		session, _ = client.NewSession()
		session.Run(`docker load --input ` + `/home/ubuntu/` + files["name"] + ".tar")
		time.Sleep(10 * time.Second)
		session, _ = client.NewSession()
		session.Run(`rm ` + `/home/ubuntu/` + files["name"] + ".tar")
	}
	// copy config file
	for _, files := range copyFiles {
		srcFile, err := os.Open(dir + files["oldFile"])
		fmt.Println("copy file is : ", files["oldFile"])
		if err != nil {
			// info := fmt.Sprintf("open %s failed :%s", files["oldFile"], err)
			// fmt.Println(info)
			return
		}
		defer srcFile.Close()
		session, _ = client.NewSession()
		session.Run(`rm ` + files["newFile"] + ``)
		dstFile, err := copySession.Create(`` + files["newFile"] + ``)
		if err != nil {
			// info := fmt.Sprintf("create %s failed :%s", files["newFile"], err)
			// fmt.Println(info)
			return
		}
		defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			// fmt.Println("readall err: ", err)
			return
		}
		dstFile.Write(ff)
	}
	time.Sleep(10 * time.Second)
	// fmt.Println("copy file successful")
}

// rm docker image which need upgrade
func deleteImage(client *ssh.Client, deleteImages []map[string]string) {
	// fmt.Println("start delete image")
	var session *ssh.Session
	session, _ = client.NewSession()
	bufRet, _ := session.CombinedOutput(`docker images`)
	imagesInfoRet := strings.Split(string(bufRet), "\n")
	// fmt.Println(string(bufRet))
	for _, imageLine := range imagesInfoRet {
		repo := strings.Split(strings.Split(imageLine, " ")[0], "armhf-")
		imageSlice := strings.Split(imageLine, " ")
		if len(repo) > 1 {
			imageName := repo[1]
			for _, delImage := range deleteImages {
				delImageName := strings.Split(delImage["name"], ".")[0]
				if imageName == delImageName {
					var imageTag, imageId string
					countTag := 0
					for i, imageTemp := range imageSlice {
						if "" != imageTemp {
							countTag++
						}

						if 2 == countTag {
							imageTag = imageSlice[i]
							break
						}
					}
					countId := 0
					for i, imageTemp := range imageSlice {
						if "" != imageTemp {
							countId++
						}

						if 3 == countId {
							imageId = imageSlice[i]
							break
						}
					}
					if imageTag == delImage["NewTag"] {
						continue
					}
					if "" != imageTag && "" != imageId {
						// fmt.Println("")
						// fmt.Println("start delete ", imageName)
						session, _ = client.NewSession()
						session.Run(`etcdctl rm /config/image/armhf-` + imageName + `/` + imageTag + `/image`)
						session, _ = client.NewSession()
						session.Run(`etcdctl rm /config/image/armhf-` + imageName + `/` + imageTag + `/script`)
						session, _ = client.NewSession()
						session.Run(`etcdctl rmdir /config/image/armhf-` + imageName + `/` + imageTag + `/`)
						session, _ = client.NewSession()
						session.Run(`docker rm -f ` + imageName)
						time.Sleep(5 * time.Second)
						session, _ = client.NewSession()
						session.Run(`docker rmi ` + imageId)
						fmt.Println("image name is : ", imageName, "image id is : ", imageId)
						time.Sleep(10 * time.Second)
					}
				}
			}
		}
	}
	session, _ = client.NewSession()
	session.Run(`etcdctl set /config/global/sensor_uid ""`)
	session, _ = client.NewSession()
	session.Run(`etcdctl set /config/global/sensor_sn ""`)
	session, _ = client.NewSession()
	session.Run(`rm /run/shm/tegra_init`)
}

// read ip form file
func readIp() string {
	var ip string
	bytes, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	strIp := string(bytes)
	ret := strings.Split(strIp, "=")
	if 2 == len(ret) {
		ip = ret[1]
	} else {
		fmt.Println("the content of config file is not right, please check id")
		os.Exit(1)
	}
	return ip
}

