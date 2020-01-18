/*
#Time      :  2020/1/16 2:29 PM 
#Author    :  chuangangshen@deepglint.com
#File      :  fgdevinfo.go
#Software  :  GoLand
*/
package main

import (
	"io/ioutil"
	"fmt"
)

const fileName = "/home/deepglint/AppData/libraT/config/devinfo.txt"

func main() {
	catDevInfo()
}

func catDevInfo() {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error : %s", err)
		return
	}
	fmt.Println(string(bytes))
}

