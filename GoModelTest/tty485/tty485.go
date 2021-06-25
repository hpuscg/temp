package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"os"
	"time"
)

var (
	gLog *log.Logger
)

func main() {
	logFile, err := os.OpenFile("./log", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Printf("open file error=%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	logWriters := io.MultiWriter(writers...)
	gLog = log.New(logWriters, "", log.Ldate|log.Ltime|log.Lshortfile)
	io.MultiWriter()
	// glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./"), glog.WithLevel("info"))
	// readDataFromTty()

	flag.Parse()
	getTtyData()
}

func getTtyData() {
	var (
		// preStatus  string
		realStatus string
	)
	config := &serial.Config{
		Name:        "/dev/ttyS0",
		Baud:        9600,
		ReadTimeout: 5 * time.Second,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		gLog.Println(err)
		return
	}
	defer s.Close()
	for {
		var strTemp string
		bufTemp := make([]byte, 11)
		num, err := s.Read(bufTemp)
		if err != nil {
			gLog.Println(err)
			continue
		}
		if num > 0 {
			strTemp = hex.EncodeToString(bufTemp)
			gLog.Println(strTemp)
		}
		// ZT-IRSXD011
		/*for key, value := range bufTemp {
			if 6 == key {
				realStatus = fmt.Sprintf("%x", value)
			}
		}
		if "55" == preStatus && "aa" == realStatus {
			// if ("55" == preStatus || "" == preStatus) && "aa" == realStatus {

			gLog.Println("begin")
		}*/

		for key, value := range strTemp {
			if 1 == key {
				if "0" == string(value) {
					goto L
				}
			}
			if 11 == key {
				realStatus = string(value)
			}
			// gLog.Printf("key:%d, value:%s", key, string(value))
		}
		/*if "0" == preStatus && "1" == realStatus {

			gLog.Println("begin")
		} else if "1" == preStatus && "0" == realStatus {

			gLog.Println("end")
		}*/
		time.Sleep(100 * time.Microsecond)
		gLog.Println(realStatus)
	L:
		continue
		// preStatus = realStatus
	}
}
