/*
#Time      :  2019/6/11 上午9:48 
#Author    :  chuangangshen@deepglint.com
#File      :  fileTest.go
#Software  :  GoLand
*/
package main

import (
	"syscall"
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"os"
)

func main() {
	GetFile()
	// GetDirFileNames()
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


