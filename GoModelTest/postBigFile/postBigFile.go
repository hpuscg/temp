/*
#Time      :  2019/7/9 下午3:49 
#Author    :  chuangangshen@deepglint.com
#File      :  postBigFile.go
#Software  :  GoLand
*/
package main

import (
	"crypto/rand"
	"io"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"os"
	"errors"
)

const (
	Url     = "http://192.168.101.238:8008/api/upload"
	FilPath = "/Users/hpu_scg/gocode/src/temp/GoModelTest/postBigFile/server_docker_images.tar"
)

func main() {
	upload(Url, FilPath)
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func upload(url, flpath string) {
	body := NewCircleByteBuffer(1024 * 2)
	boundary := randomBoundary()
	boundarybytes := []byte("\r\n--" + boundary + "\r\n")
	endbytes := []byte("\r\n--" + boundary + "--\r\n")

	reqest, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)
	go func() {
		//defer ruisRecovers("upload.run")
		f, err := os.OpenFile(flpath, os.O_RDONLY, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，但是因为windows下还有BUG，如果使用这个O_CREATE，就会直接清空文件，所以这里就不用了这个标志，你自己事先建立好文件。
		if err != nil {
			panic(err)
		}
		stat, err := f.Stat() //获取文件状态
		if err != nil {
			panic(err)
		}
		defer f.Close()

		header := fmt.Sprintf("Content-Disposition: form-data; name=\"upfile\"; filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", stat.Name())
		body.Write(boundarybytes)
		body.Write([]byte(header))

		fsz := float64(stat.Size())
		fupsz := float64(0)
		buf := make([]byte, 1024)
		for {
			time.Sleep(10 * time.Microsecond) //减缓上传速度，看进度效果
			n, err := f.Read(buf)
			if n > 0 {
				nz, _ := body.Write(buf[0:n])
				fupsz += float64(nz)
				progress := strconv.Itoa(int((fupsz/fsz)*100)) + "%"
				fmt.Println("upload:", progress, "|", strconv.FormatFloat(fupsz, 'f', 0, 64), "/", stat.Size())
			}
			if err == io.EOF {
				break
			}
		}
		body.Write(endbytes)
		body.Write(nil) //输入EOF,表示数据写完了
	}()
	resp, err := http.DefaultClient.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("上传成功")
	} else {
		fmt.Println("上传失败,StatusCode:", resp.StatusCode, resp)
	}
}

type CircleByteBuffer struct {
	io.Reader
	io.Writer
	io.Closer
	datas []byte

	start   int
	end     int
	size    int
	isClose bool
	isEnd   bool
}

func NewCircleByteBuffer(len int) *CircleByteBuffer {
	var e = new(CircleByteBuffer)
	e.datas = make([]byte, len)
	e.start = 0
	e.end = 0
	e.size = len
	e.isClose = false
	e.isEnd = false
	return e
}

func (e *CircleByteBuffer) getLen() int {
	if e.start == e.end {
		return 0
	} else if e.start < e.end {
		return e.end - e.start
	} else {
		return e.start - e.end
	}
}

func (e *CircleByteBuffer) getFree() int {
	return e.size - e.getLen()
}
func (e *CircleByteBuffer) putByte(b byte) error {
	if e.isClose {
		return io.EOF
	}
	e.datas[e.end] = b
	var pos = e.end + 1
	for pos == e.start {
		if e.isClose {
			return io.EOF
		}
		time.Sleep(time.Microsecond)
	}
	if pos == e.size {
		e.end = 0
	} else {
		e.end = pos
	}
	return nil
}

func (e *CircleByteBuffer) getByte() (byte, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if e.isEnd && e.getLen() <= 0 {
		return 0, io.EOF
	}
	if e.getLen() <= 0 {
		return 0, errors.New("no datas")
	}
	var ret = e.datas[e.start]
	e.start++
	if e.start == e.size {
		e.start = 0
	}
	return ret, nil
}

func (e *CircleByteBuffer) geti(i int) byte {
	if i >= e.getLen() {
		panic("out buffer")
	}
	var pos = e.start + i
	if pos >= e.size {
		pos -= e.size
	}
	return e.datas[pos]
}

func (e *CircleByteBuffer) Close() error {
	e.isClose = true
	return nil
}
func (e *CircleByteBuffer) Read(bts []byte) (int, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if bts == nil {
		return 0, errors.New("bts is nil")
	}
	var ret = 0
	for i := 0; i < len(bts); i++ {
		b, err := e.getByte()
		if err != nil {
			if err == io.EOF {
				return ret, err
			}
			return ret, nil
		}
		bts[i] = b
		ret++
	}
	if e.isClose {
		return ret, io.EOF
	}
	return ret, nil
}

func (e *CircleByteBuffer) Write(bts []byte) (int, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if bts == nil {
		e.isEnd = true
		return 0, io.EOF
	}
	var ret = 0
	for i := 0; i < len(bts); i++ {
		err := e.putByte(bts[i])
		if err != nil {
			fmt.Println("Write bts err:", err)
			return ret, err
		}
		ret++
	}
	if e.isClose {
		return ret, io.EOF
	}
	return ret, nil
}
