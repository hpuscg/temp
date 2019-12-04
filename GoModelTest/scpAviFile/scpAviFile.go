/*
#Time      :  2019/10/04 上午11:23
#Author    :  chuangangshen@deepglint.com
#File      :  scpAviFile.go
#Software  :  GoLand
*/
package main

import (
	"io/ioutil"
	"strings"
	"errors"
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"io"
	"fmt"
	"time"
	"os/exec"
	"strconv"
)

var (
	AviFileName string
	Uid         string
	TarGetUrl   string
	LimitTime   int64 = 300
)

const (
	DirName = "/libra/dataset/scene"
)

func main() {
	InitServerInfo()
	if TarGetUrl == "" {
		return
	}
	err := GetAviFileName()
	if err != nil {
		fmt.Println(err)
		return
	}
	isFileOk := CheckFileTime()
	if isFileOk != nil {
		fmt.Println(isFileOk)
		return
	}
	TarAviFile()
}

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
	limitTimeStr := "etcdctl get /config/libra/sensor/3D_recorder_seconds"
	limitTime := ExecCmd(limitTimeStr)
	limitTimeInt, err := strconv.ParseInt(strings.TrimSpace(limitTime), 10, 64)
	if err == nil {
		LimitTime = limitTimeInt
	} else {
		fmt.Println(err)
	}
}

// 获取/Libra/dataset/scence下的AVI文件名
func GetAviFileName() (err error) {
	files, _ := ioutil.ReadDir(DirName)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".avi") {
			AviFileName = f.Name()
			return nil
		}
	}
	return errors.New("no avi file")
}

// 校验AVI文件最后修改时间是否达到5分钟
func CheckFileTime() (err1 error) {
	AviFIlePath := DirName + "/" + AviFileName
	f, err := os.Open(AviFIlePath)
	if err != nil {
		err1 = err
	}
	fi, err := f.Stat()
	if err != nil {
		err1 = err
	}
	timeNow := time.Now().Unix()
	fmt.Println(LimitTime, timeNow, fi.ModTime().Unix())
	if timeNow-fi.ModTime().Unix() < LimitTime {
		err1 = errors.New("avi file modify less 5 minute")
	}
	return
}

// 将ONI文件和模型文件压缩并上传到网管服务器
func TarAviFile() {
	rmOldTarFile := "rm -rf /root/*.tar.gz"
	ExecCmd(rmOldTarFile)
	postFileName := Uid + "_" + strconv.Itoa(int(time.Now().Unix())) + ".tar.gz"
	tarFileStr := "tar -zcvf " + postFileName + " " + DirName
	ExecCmd(tarFileStr)
	err := PostBigFile(postFileName)
	fmt.Println(err)
	rmTarPicture := "rm " + postFileName
	ExecCmd(rmTarPicture)
	rmPicture := "rm " + DirName + "/" + AviFileName
	ExecCmd(rmPicture)
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

// 上传压缩包文件到网管服务器
func PostFile(fileName string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// fileWriter, err := bodyWriter.CreateFormFile("uploadfile", fileName)
	_, err := bodyWriter.CreateFormFile("uploadfile", fileName)
	if err != nil {
		return err
	}
	fh, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fh.Close()

	boundary := bodyWriter.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(bodyBuf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", fileName)
		return err
	}
	req, err := http.NewRequest("POST", TarGetUrl, request_reader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(close_buf.Len())

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp)

	/*_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(TarGetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}*/
	return nil
}

// 上传大文件到网管服务器
func PostBigFile(fileName string) error {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("uploadfile", fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		if _, err := io.Copy(part, file); err != nil {
			return
		}
	}()
	resp, err := http.Post(TarGetUrl, m.FormDataContentType(), r)
	fmt.Println(resp.Body)
	return err
}
