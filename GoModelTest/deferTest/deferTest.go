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
	// testDefer()
	defer_call()
}

func defer_call() {

	defer func() {
		if err := recover();err != nil {
			fmt.Println(err) //err 就是panic传入的参数
		}
		fmt.Println("打印前")
	}()

	defer func() { // 必须要先声明defer，否则recover()不能捕获到panic异常

		fmt.Println("打印中")
	}()


	defer func() {

		fmt.Println("打印后")
	}()
	panic("触发异常")
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



