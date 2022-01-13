package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var (
	gLog       *log.Logger
	sensorType string
)

func main() {
	tty485()
	// pathTest()
}

func pathTest() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strings.Replace(dir, "\\", "/", -1))
}

func tty485() {
	flag.StringVar(&sensorType, "type", "sw-wb", "relay type")
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
	fmt.Println(sensorType)
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
		switch sensorType {
		case "d011": // ZT-IRSXD011
			for key, value := range bufTemp {
				if key == 6 {
					realStatus = fmt.Sprintf("%x", value)
				}
			}
		case "sw-wb": // sw-wb
			for key, value := range strTemp {
				if key == 1 {
					if string(value) == "0" {
						goto L
						// realStatus = string(value)
					}
				}
				if key == 11 {
					realStatus = string(value)
				}
			}
		case "d012": // ZT-IRSXD012
			for key, value := range strTemp {
				if key == 11 {
					realStatus = string(value)
				}
			}
		}
		gLog.Println(realStatus)
	L:
		continue
	}
}
