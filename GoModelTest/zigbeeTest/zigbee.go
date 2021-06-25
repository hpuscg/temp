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
	gLog       *log.Logger
	DataType   int64
	isContinue bool
	s          *serial.Port
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

	var (
		read  bool
		write bool
	)
	flag.BoolVar(&read, "read", false, "read")
	flag.BoolVar(&write, "write", false, "write")
	flag.Int64Var(&DataType, "data", 0, "data type")
	flag.Parse()
	if !read && !write {
		gLog.Println("need read or write")
	}
	/*go func() {
		time.Sleep(10 * time.Second)
		isContinue = true
		// err := s.Close()
		// gLog.Println(err)
		gLog.Println("999999")
	}()*/
	if read {
		getTtyData()
	}
	if write {
		zigbeeTest()
	}
}

func readDataFromTty() {
	config := &serial.Config{
		Name: "/dev/ttyS0",
		Baud: 9600,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		gLog.Println(err)
		return
	}
	defer s.Close()
	for {
		bufTemp := make([]byte, 11)
		num, err := s.Read(bufTemp)
		if err != nil {
			gLog.Println(err)
		}
		if num > 0 {
			strTemp := hex.EncodeToString(bufTemp)
			gLog.Println(strTemp)
		}

	}
}

const (
	Data1 = "\xFE\x0B\x02\x01\x02\x00\x00\x00\x1B\x3C\xD3"
	Data2 = "\xFE\x0B\x02\x01\x02\x00\x00\x00\x00\x3C\xC8"
	Data3 = "\xFE\x0B\x02\x03\x02\x00\x00\x00\x1B\x3C\xD3"
	Data4 = "\xFE\x0B\x02\x04\x02\x00\x00\x00\x1B\x3C\xD3"
	Data5 = "\xFE\x0B\x02\x07\x02\x00\x00\x00\x1B\x3C\xD3"
	Data7 = "\xFE\x0B\x09\x14\x02\x00\x00\x00\x1B\x3C\xD3"
	Data8 = "\xFE\x0B\x04\x03\x02\x00\x00\x00\x1B\x3C\xD3"
	Data9 = "\xFE\x0A\x01\x02\x02\x00\x00\x00\x88\x7D"
)

func zigbeeTest() {
	config := &serial.Config{
		Name: "/dev/ttyS4",
		Baud: 115200,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		gLog.Println(err)
		return
	}
	defer s.Close()
	data := Data1
	switch DataType {
	case 1:
		data = Data1
	case 2:
		data = Data2
	case 3:
		data = Data3
	case 4:
		data = Data4
	case 5:
		data = Data5
	case 7:
		data = Data7
	case 8:
		data = Data8
	case 9:
		data = Data9
	}
	n, err := s.Write([]byte(data))
	if err != nil {
		gLog.Println(err)
	}
	gLog.Println(n)
}

func getTtyData() {
	config := &serial.Config{
		Name:        "/dev/ttyS4",
		Baud:        115200,
		ReadTimeout: 5 * time.Second,
	}
	var err error
	s, err = serial.OpenPort(config)
	if err != nil {
		gLog.Println(err)
		return
	}
	defer s.Close()
	for !isContinue {
		gLog.Println("33333")
		bufTemp := make([]byte, 50)
		num, err := s.Read(bufTemp)
		if err != nil {
			gLog.Println(err)
		}
		if num > 0 {
			gLog.Println(bufTemp)
			length := bufTemp[1]
			tempData := bufTemp[:length]
			dataType := bufTemp[2:4]
			address := bufTemp[5:7]
			gLog.Printf("%x\n", dataType)
			gLog.Printf("%x\n", tempData)
			gLog.Println(bufTemp[9])
			gLog.Println(address)
		}

	}
}
