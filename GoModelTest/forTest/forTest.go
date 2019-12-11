/*
#Time      :  2019/12/11 9:47 AM 
#Author    :  chuangangshen@deepglint.com
#File      :  forTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	ForTest()
}

func ForTest() {
	for i := 0; i <= 10; i++ {
		fmt.Println(i)
		time.Sleep(200 * time.Millisecond)
	}
}

