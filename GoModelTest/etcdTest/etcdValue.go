/*
#Time      :  2019/3/11 下午12:49 
#Author    :  chuangangshen@deepglint.com
#File      :  etcdValue.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	SetEtcdValue()
}

var PassWord = "nonghang!@#$"

func SetEtcdValue() {
	data, err := json.Marshal(PassWord)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
	/*etcdStr := "etcdctl set /config/iot/password " + string(data)
	cmd := exec.Command("/bin/bash", "-c", etcdStr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}*/
}
