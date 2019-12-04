/*
#Time      :  2019/7/15 上午10:15 
#Author    :  chuangangshen@deepglint.com
#File      :  marshalTest.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	marshalError()
}

func marshalError() {
	var err error
	_, err = strconv.Atoi("er")
	data, err1 := json.Marshal(err)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("26===", string(data))
	}
}

