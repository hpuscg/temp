/*
#Time      :  2019/4/4 下午6:13 
#Author    :  chuangangshen@deepglint.com
#File      :  interfaceJson.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	InterfaceJson()
}

type JsonTest struct {
	age  int
	name string
}

func InterfaceJson() {
	data := makeStruct()
	ret, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	resp := string(ret)
	fmt.Println(resp)
	var value JsonTest
	err = json.Unmarshal([]byte(resp), &value)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)
}

func makeStruct() interface{} {
	data := JsonTest{
		age:12,
		name:"123",
	}
	return data
}


