/*
#Time      :  2019/7/24 下午4:40 
#Author    :  chuangangshen@deepglint.com
#File      :  osPlayWav.go
#Software  :  GoLand
*/
package main

import (
	"flag"
	"os/exec"
	"fmt"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "spm.wav", "wav file name")
	flag.Parse()
	cmd := exec.Command("aplay", fileName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("run cmd err:", err)
	} else {
		fmt.Println("yes")
	}
}
