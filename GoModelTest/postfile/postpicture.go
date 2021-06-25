/*
#Time      :  2019/2/25 上午11:22
#Author    :  chuangangshen@deepglint.com
#File      :  postfile.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	PostFile("1.txt", "http://127.0.0.1:8789/event/upload/file")
}

func run() {
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
	fileWriter, err := bodyWriter.CreateFormFile("uploadFile", fileName)
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
	params := map[string]string{
		"eventId": "eventId",
	}
	for key, val := range params {
		fmt.Println(key, val)
		_ = bodyWriter.WriteField(key, val)
	}
	// resp, err := http.Post(targetUrl, contentType, bodyBuf)
	reqData := io.MultiReader(bodyBuf)
	resp, err := PostFileWithToken(targetUrl, contentType, reqData)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	/*defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}*/
	return nil
}

func PostFileWithToken(url, contentType string, data io.Reader) (body []byte, err error) {
	client := &http.Client{}
	var (
		req      *http.Request
		response *http.Response
	)
	req, err = http.NewRequest("POST", url, data)
	req.Close = true
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	// req.Header.Add("authorization", token)
	req.Header.Set("Connection", "close")
	req.Header.Set("eventId", "no-cache")
	// req.Form.Add("eventIds", "000")
	response, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func RemoveOldData() {
	for i := 1; i < 365; i++ {
		strTime := strconv.Itoa(-24*i) + "h"
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
