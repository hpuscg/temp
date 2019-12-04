/*
#Time      :  2019/10/31 下午5:03 
#Author    :  chuangangshen@deepglint.com
#File      :  uuidTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/satori/go.uuid"
	"fmt"
)

func main() {
	CreateUuid()
}

func CreateUuid() {
	tmpstr, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmpstr.String())
}

// 4dbb36e7-7b81-4392-a357-71d589781acf

