package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// osReadFile()
	ioutilReadFile()
}

func osReadFile()  {
	file, err := os.Open("config.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	bytesRead, _ := file.Read(buffer)
	ipTmp := string(bytesRead)
	fmt.Println(ipTmp)
	fmt.Println(fileSize)
	fmt.Println(bytesRead)
}

func ioutilReadFile()  {
	bytes, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println(err)
	}
	strFile := string(bytes)
	fmt.Println(strFile)
	ret := strings.Split(strFile, "=")
	fmt.Println(ret[1])
}

