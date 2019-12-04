/*
#Time      :  2019/6/27 下午4:25 
#Author    :  chuangangshen@deepglint.com
#File      :  scpEventFile.go
#Software  :  GoLand
*/
package main

import (
	"strings"
	"fmt"
	"bytes"
	"os/exec"
	"strconv"
	"time"
	"mime/multipart"
	"net/http"
	"io/ioutil"
	"os"
	"io"
)

func main() {
	InitServerInfo()
	TarEventFile()
}

var (
	Uid string
	TarGetUrl string
)

const (
	EventDirName = "/data/tf/eventserver/event"
)

// init网管服务器信息
func InitServerInfo() {
	serverStr := "etcdctl get /config/global/server_addr"
	serverAddr := ExecCmd(serverStr)
	if serverAddr == "\n" {
		return
	}
	serverAddr = strings.TrimSpace(serverAddr)
	TarGetUrl = "http://" + serverAddr + ":8008/api/upload"
	fmt.Println(TarGetUrl)
	uidStr := "etcdctl get /config/global/sensor_uid"
	uid := ExecCmd(uidStr)
	Uid = strings.TrimSpace(uid)
}

// 将event压缩并上传到网管服务器
func TarEventFile() {
	_, err := os.Stat(EventDirName)
	if os.IsNotExist(err) {
		return
	}
	postFileName := "/tmp/" + Uid + "_" + strconv.Itoa(int(time.Now().Unix())) + "_event.tar.gz"
	tarFileStr := "tar -zcvf " + postFileName + " " + EventDirName
	ExecCmd(tarFileStr)
	err = PostFile(postFileName)
	fmt.Println(err)
	rmTarPicture := "rm " + postFileName
	ExecCmd(rmTarPicture)
}

// 上传压缩包文件到网管服务器
func PostFile(fileName string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", fileName)
	if err != nil {
		return err
	}
	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	fmt.Println(TarGetUrl)
	fmt.Println(fileName)
	resp, err := http.Post(TarGetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// 执行命令行
func ExecCmd(str string) string {
	fmt.Println(str)
	cmd := exec.Command("/bin/bash", "-c", str)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return "\n"
	}
	ret := out.String()
	return ret
}
