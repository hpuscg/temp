/*
#Time      :  2019/11/27 上午11:23
#Author    :  chuangangshen@deepglint.com
#File      :  scpH264File.go
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
	H264FileName    string
	Uid             string
	TarGetUrl       string
	VideoTimeLength int64
)

const (
	DirName     = "/libra/dataset/scene"
	H264DirName = "/libra/dataset/scene/h264"
)

func main() {
	InitServerInfo()
	if TarGetUrl == "" {
		return
	}
	err := GetH264FileName()
	if err != nil {
		fmt.Println(err)
		return
	}
	GetVideoLength()
	isFileOk := CheckFileTime()
	if isFileOk != nil {
		fmt.Println(isFileOk)
		return
	}
	TarH264File()
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

// 获取/Libra/dataset/scence/h264下的h264文件名
func GetH264FileName() (err error) {
	files, _ := ioutil.ReadDir(H264DirName)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".h264") {
			H264FileName = f.Name()
			return nil
		}
	}
	return errors.New("no h264 file")
}

// 获取h264视频录制时长
func GetVideoLength() {
	var err error
	videoLengthStr := "etcdctl get /config/libra/sensor/3D_recorder_seconds"
	videoTime := ExecCmd(videoLengthStr)
	VideoTimeLength, err = strconv.ParseInt(videoTime, 10, 64)
	if err != nil {
		fmt.Println(err)
		VideoTimeLength = 900
	}
}

// 校验h264文件最后修改时间是否达到录制时间
func CheckFileTime() (err1 error) {
	OniFIlePath := H264DirName + "/" + H264FileName
	f, err := os.Open(OniFIlePath)
	if err != nil {
		err1 = err
	}
	fi, err := f.Stat()
	if err != nil {
		err1 = err
	}
	timeNow := time.Now().Unix()
	if timeNow-fi.ModTime().Unix() < VideoTimeLength {
		errStr := fmt.Sprintf("h264 file modify less %d seconds", VideoTimeLength)
		err1 = errors.New(errStr)
	}
	return
}

// 将h264文件和模型文件压缩并上传到网管服务器
func TarH264File() {
	rmOldTarFile := "rm -rf /root/*.tar.gz"
	ExecCmd(rmOldTarFile)
	postFileName := Uid + "_" + strconv.Itoa(int(time.Now().Unix())) + ".tar.gz"
	tarFileStr := "tar -zcvf " + postFileName + " " + DirName
	ExecCmd(tarFileStr)
	err := PostBigFile(postFileName)
	fmt.Println(err)
	rmTarPicture := "rm " + postFileName
	ExecCmd(rmTarPicture)
	rmPicture := "rm " + H264DirName + "/" + H264FileName
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
