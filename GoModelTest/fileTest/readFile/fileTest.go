/*
#Time      :  2019/6/11 上午9:48
#Author    :  chuangangshen@deepglint.com
#File      :  fileTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"gitlab.deepglint.com/junkaicao/glog"
)

func main() {
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./"), glog.WithLevel("info"))
	// GetFile()
	// GetDirFileNames()
	/* for {
		ReadTtyFile()
		time.Sleep(1 * time.Second)
	} */
	ReadFile()
}

func ReadFile() {
	f, err := os.Open("/Users/hpu_scg/Desktop/a1f20f44503536363700001600a5011d.reset_pass.lic")
	if err != nil {
		fmt.Println(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("1111")
	fmt.Println(string(data))
}

func ReadTtyFile() {
	fileName := "/dev/ttyS0"
	// fileName := "./file.txt"
	glog.Infoln("1")
	f, err := os.Open(fileName)
	if err != nil {
		glog.Infoln("2")
		glog.Infoln(err)
		return
	}
	defer f.Close()

	// var stringsData string
	for {
		buf := make([]byte, 8)
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		data, err := strconv.ParseUint(string(buf), 16, 8)
		if err != nil {
			glog.Infoln("Parse Error", err)
			return
		}
		n2 := uint8(data)
		fn := int(*(*int8)(unsafe.Pointer(&n2)))
		fmt.Println(fn)
		/*stringsData += string(buf)
		fmt.Println(stringsData)*/
	}

	/*
		glog.Infoln("3")
		fileData, err := ioutil.ReadAll(f)
		glog.Infoln("4")
		if err != nil {
			glog.Infoln("5")
			glog.Infoln(err)
			return
		}
		glog.Infoln("6")
		glog.Infoln(string(fileData))
	*/

	/*n, err := strconv.ParseUint(string(data), 16, 8)
	if err != nil {
		glog.Infoln("Parse Error")
		return
	}
	n2 := uint8(n)
	fn := int(*(*int8)(unsafe.Pointer(&n2)))
	fmt.Println(fn)*/

}

func GetFile() {
	f, _ := os.Open("file.txt")
	fi, _ := f.Stat()
	fmt.Println(fi.ModTime().Unix())
}

func timeSpecToTime(ts syscall.Timespec) time.Time {
	fmt.Println(int64(ts.Sec))
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func GetDirFileNames() {
	files, _ := ioutil.ReadDir("../fileTest")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".txt") {
			fmt.Println(f.Name())
		}
	}
}
