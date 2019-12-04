package main

import (
	"net"
	"golang.org/x/crypto/ssh"
	"fmt"
	"time"
	"os"
	"github.com/pkg/sftp"
	"io/ioutil"
	"strings"
	"os/exec"
	"runtime"
)

/*func main() {
	Ip := "192.168.12.12"
	Dir := "./"

	tryPing(Ip)


	session, err := connect("root", Ip, 22)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", Ip, err)
		fmt.Println(info)
		return
	}
	// read upgrade image file from dir
	loadImage := readImageFiles(Dir)
	// copyFiles
	copySession, err := copyConnect("root", Ip, 22)
	if err != nil {
		fmt.Println("connect sftp fail, ", err)
		return
	}
	copyFile(copySession, session, loadImage)
	// deleteImage
	deleteImage(session, loadImage)
	// reboot
	// session.Run(`reboot`)
	fmt.Println("wait about five minutes")
	// test is upgrade successful
	ret := testUpgradeIsTrue(Ip)
	if !ret {
		fmt.Println("upgrade failed, please try again!")
	} else {
		fmt.Println("upgrade successfully")
	}
}*/

func main() {
	session, err := connect("root", "192.168.5.251", 22)
	if err != nil {
		fmt.Println(err)
	}
	ret, err := session.Output("ps -aux |grep flowservice |wc -l")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(ret))
}

// connect for exec cmd
func connect(user, host string, port int) (*ssh.Session, error) {
	var (
		auth 			[]ssh.AuthMethod
		addr 			string
		clientConfig 	*ssh.ClientConfig
		client 			*ssh.Client
		session 		*ssh.Session
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
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)

	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: 				user,
		Auth: 				auth,
		Timeout: 			60 * time.Second,
		HostKeyCallback: 	func(hostname string, remote net.Addr, key ssh.PublicKey) error {
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

// test the ip is right
func tryPing(IP string)  {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if "windows" == sysInfo {
		cmd = exec.Command("ping", IP)
	} else {
		cmd = exec.Command("ping", "-t", "1", IP)
	}
	_, err := cmd.Output()
	if err != nil {
		info := fmt.Sprintf("can't ping %s,the err is : %s, please checkout your config ip", IP, err)
		fmt.Println(info)
		os.Exit(1)
	} else {
		fmt.Println("this ip is ok")
	}
}

// test is upgrade success
func testUpgradeIsTrue(Ip string) bool {
	var count int
	for {
		if count < 100 {
			var cmd *exec.Cmd
			sysInfo := runtime.GOOS
			if "windows" == sysInfo {
				cmd = exec.Command("ping", Ip)
			} else {
				cmd = exec.Command("ping", "-t", "1", Ip)
			}
			// cmd := exec.Command("ping", IP)  // for windows
			_, err := cmd.Output()
			if err != nil {
				count++
				continue
			}	else {
				return true
			}
		} else {
			fmt.Println("upgrade failed, please try again!")
			return false
		}
	}
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

// rm docker image which need upgrade
func deleteImage(session *ssh.Session, deleteImages []map[string]string) {
	fmt.Println("start delete image")
	bufRet, _ := session.CombinedOutput(`docker images`)
	imagesInfoRet := strings.Split(string(bufRet), "\n")
	// fmt.Println(string(bufRet))
	session.Run("docker rm -f $(docker ps -a -q)")
	for _, imageLine := range imagesInfoRet {
		repo := strings.Split(strings.Split(imageLine, " ")[0], "armhf-")
		imageSlice := strings.Split(imageLine, " ")
		// fmt.Println(repo)
		// fmt.Println(imageLine)
		if len(repo) > 1 {
			imageName := repo[1]
			// fmt.Println(imageName)
			// fmt.Println("deleteImage is: ", deleteImages)
			for _, delImage := range deleteImages {
				// fmt.Println("delimage is :", delImage)
				// fmt.Println("imageName is :", imageName)
				if imageName == delImage["name"] {
					// imageTag := strings.Split(strings.Split(imageLine, " ")[1], " ")[0]
					// imageId := strings.Split(strings.Split(strings.Split(imageLine, " ")[1]))
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
						fmt.Println("")
						fmt.Println("start delete ", imageName, ": ", imageTag, ":", imageId)
						session.Run(`etcdctl rm /config/image/armhf-` + imageName + `/` + imageTag + `/image`)
						session.Run(`etcdctl rm /config/image/armhf-` + imageName + `/` + imageTag + `/script`)
						session.Run(`etcdctl rmdir /config/image/armhf-` + imageName + `/` + imageTag + `/`)
						session.Run(`docker rmi ` + imageId)
						time.Sleep(10 * time.Second)
						/*session.Start(`docker rmi ` + imageId)
						errWait := session.Wait()
						fmt.Println("err wait ", errWait)*/
						fmt.Println("delete " + imageName + " over")
						fmt.Println("")
					}
				}
			}
		}
	}
	session.Run(`etcdctl set /config/global/sensor_uid ""`)
	session.Run(`etcdctl set /config/global/sensor_sn ""`)
	session.Run(`rm /run/shm/tegra_init`)
	session.Run(`mkdir /libra/judicial`)
	fmt.Println("delete image successful")
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
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)
	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: 				user,
		Auth: 				auth,
		Timeout: 			60 * time.Second,
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

// copy file and load docker image
func copyFile(copySession *sftp.Client, session *ssh.Session, copyImageFiles []map[string]string) {
	fmt.Println("start copy files")
	copyFiles := []map[string]string {{"oldFile":"Release", "newFile":"/home/ubuntu/Release"},
		{"oldFile":"70-persistent-net.rules", "newFile":"/etc/udev/rules.d/70-persistent-net.rules"},
		{"oldFile":"reload_sensor.sh", "newFile":"/data/shell/_usrbin/reload_sensor.sh"},
		{"oldFile":"differCompany.config", "newFile":"/libra/judicial/differCompany.config"}}
	// load docker image
	for _, files := range copyImageFiles {
		srcFile, err := os.Open(files["name"] + "." + files["NewTag"] + `.tar`)
		fmt.Println(files)
		if err != nil {
			info := fmt.Sprintf("open %s failed :%s", files["name"] + "." + files["NewTag"] + `.tar`, err)
			fmt.Println(info)
		}
		// defer srcFile.Close()
		dstFile, err := copySession.Create(`/home/ubuntu/` + files["name"] + ``)
		if err != nil {
			info := fmt.Sprintf("create %s failed :%s", files["name"], err)
			fmt.Println(info)
		}
		// defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			fmt.Println("readall err: ", err)
		}
		dstFile.Write(ff)
		session.Run(`docker load --input ` + files["name"] + ``)
		time.Sleep(10 * time.Second)
		session.Run(`rm /home/ubuntu/` + files["name"] + ``)
	}
	// copy config file
	for _, files := range copyFiles {
		fmt.Println(files)
		srcFile, err := os.Open(files["oldFile"])
		if err != nil {
			info := fmt.Sprintf("open %s failed :%s", files["oldFile"], err)
			fmt.Println(info)
		}
		// defer srcFile.Close()
		session.Run(`rm ` + files["newFile"] + ``)
		dstFile, err := copySession.Create(`` + files["newFile"] + ``)
		if err != nil {
			info := fmt.Sprintf("create %s failed :%s", files["newFile"], err)
			fmt.Println(info)
		}
		// defer dstFile.Close()
		ff,err := ioutil.ReadAll(srcFile)
		if err != nil {
			fmt.Println("readall err: ", err)
		}
		dstFile.Write(ff)
	}
	fmt.Println("copy file successful")
}

