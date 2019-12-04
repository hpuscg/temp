/*
#Time      :  2019/6/11 上午11:23
#Author    :  chuangangshen@deepglint.com
#File      :  scpOniFile.go
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
	OniFileName string
	Uid string
	TarGetUrl string
)

const (
	DirName = "/libra/dataset/scene"
)

func main() {
	InitServerInfo()
	if TarGetUrl == "" {
		return
	}
	err := GetOniFileName()
	if err != nil {
		fmt.Println(err)
		return
	}
	isFileOk := CheckFileTime()
	if isFileOk != nil  {
		fmt.Println(isFileOk)
		return
	}
	TarOniFile()
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
}

// 获取/Libra/dataset/scence下的ONI文件名
func GetOniFileName() (err error) {
	files, _ := ioutil.ReadDir(DirName)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".oni") {
			OniFileName = f.Name()
			return nil
		}
	}
	return errors.New("no oni file")
}

// 校验ONI文件最后修改时间是否达到5分钟
func CheckFileTime() (err1 error) {
	OniFIlePath := DirName + "/" + OniFileName
	f, err := os.Open(OniFIlePath)
	if err != nil {
		err1 = err
	}
	fi, err := f.Stat()
	if err != nil {
		err1 = err
	}
	timeNow := time.Now().Unix()
	if timeNow - fi.ModTime().Unix() < 300 {
		err1 = errors.New("oni file modify less 5 minute")
	}
	return
}

// 将ONI文件和模型文件压缩并上传到网管服务器
func TarOniFile() {
	rmOldTarFile := "rm -rf /root/*.tar.gz"
	ExecCmd(rmOldTarFile)
	postFileName := Uid + "_" + strconv.Itoa(int(time.Now().Unix())) + ".tar.gz"
	tarFileStr := "tar -zcvf " + postFileName + " " + DirName
	ExecCmd(tarFileStr)
	err := PostBigFile(postFileName)
	fmt.Println(err)
	rmTarPicture := "rm " + postFileName
	ExecCmd(rmTarPicture)
	rmPicture := "rm " + DirName + "/" + OniFileName
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
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", fileName)
		return err
	}
	req, err := http.NewRequest("POST", TarGetUrl, requestReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp)
	return nil
}

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