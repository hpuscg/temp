/*
#Time      :  2019/1/21 下午4:11 
#Author    :  chuangangshen@deepglint.com
#File      :  deferTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	testDefer()
}

func testDefer()  {
	defer func() {
		fmt.Println("a")
		if err := recover(); err != nil {
			fmt.Println("the err is: ", err)
			time.Sleep(10 * time.Second)
		}
		fmt.Println("c")
	}()

	a := 0
	b := 10
	c := b / a
	fmt.Println(c)
}



