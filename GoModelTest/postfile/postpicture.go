/*
#Time      :  2019/2/25 上午11:22 
#Author    :  chuangangshen@deepglint.com
#File      :  postfile.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"mime/multipart"
	"os"
	"io"
	"net/http"
	"io/ioutil"
	"os/exec"
	"time"
	"strings"
	"strconv"
	"fmt"
)

func main() {
	serverStr := "etcdctl get /config/global/server_addr"
	serverAddr := ExecCmd(serverStr)
	if serverAddr == "\n" {
		return
	}
	serverAddr = strings.TrimSpace(serverAddr)

	targetUrl := "http://" + serverAddr + ":8008/api/upload"
	fmt.Println(targetUrl)

	uidStr := "etcdctl get /config/global/sensor_uid"
	uid := ExecCmd(uidStr)
	uid = strings.TrimSpace(uid)

	descStr := "etcdctl get /config/global/sensor_desc"
	desc := ExecCmd(descStr)
	desc = strings.TrimSpace(desc)

	d, _ := time.ParseDuration("-24h")
	yearMonthDay := time.Now().Add(d).Format("20060102")
	filePath := "/tmp/" + yearMonthDay + "/"
	_, err := os.Stat(filePath)
	if err == nil {
		tarFileStr := "tar zcvf " + desc + "_" + uid + "_" + yearMonthDay + ".tar.gz " + filePath
		ExecCmd(tarFileStr)
		rmPicture := "rm -rf /tmp/" + yearMonthDay
		ExecCmd(rmPicture)
		postFileName := desc + "_" + uid + "_" + yearMonthDay + ".tar.gz"
		PostFile(postFileName, targetUrl)
		rmTarPicture := "rm " + desc + "_" + uid + "_" + yearMonthDay + ".tar.gz"
		ExecCmd(rmTarPicture)
	}
	fallYearMonthDay := "fall_detection" + time.Now().Add(d).Format("20060102")
	fallFilePath := "/tmp/" + fallYearMonthDay + "/"
	_, err = os.Stat(fallFilePath)
	if err == nil {
		tarFileStr := "tar -zcvf " + desc + "_" + uid + "_" + fallYearMonthDay + ".tar.gz " + fallFilePath
		ExecCmd(tarFileStr)
		rmPicture := "rm -rf /tmp/" + fallYearMonthDay
		ExecCmd(rmPicture)
		postFileName := desc + "_" + uid + "_" + fallYearMonthDay + ".tar.gz"
		PostFile(postFileName, targetUrl)
		rmTarPicture := "rm " + desc + "_" + uid + "_" + fallYearMonthDay + ".tar.gz"
		ExecCmd(rmTarPicture)
	}
	RemoveOldData()
}

func ExecCmd(str string) string {
	cmd := exec.Command("/bin/bash", "-c", str)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "\n"
	}
	ret := out.String()
	return ret
}

func PostFile(fileName, targetUrl string) error {
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
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
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

func RemoveOldData() {
	for i := 1; i < 365 ; i++ {
		strTime := strconv.Itoa(-24 * i) + "h"
		d, _ := time.ParseDuration(strTime)
		yearMonthDay := time.Now().Add(d).Format("20060102")
		filePath := "/tmp/" + yearMonthDay + "/"
		_, err := os.Stat(filePath)
		if err == nil {
			rmPicture := "rm -rf " + filePath
			ExecCmd(rmPicture)
		}
	}
}
