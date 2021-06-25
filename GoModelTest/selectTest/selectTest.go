/*
#Time      :  2019/5/15 上午10:21
#Author    :  chuangangshen@deepglint.com
#File      :  selectTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"time"
)

var (
	Test1 chan string
	Test2 = make(chan string)
	a     = "123"
)

func main() {
	// go SendChannel()
	// SelectTest()
	// fmt.Println(time.Now().UnixNano())
	// fmt.Println(time.Now().Unix())
	// test2()
	// selectChannel()
	switchTest(a)
}

func switchTest(a interface{}) {
	switch v := a.(type) {
	case string:
		fmt.Print(v)
	}
}

func selectChannel() {
	for {
		select {
		case a := <-Test2:
			fmt.Println(a)
			return
		default:
			fmt.Println("==yes===")
			time.Sleep(1 * time.Second)
		}
	}
}

func test2() {
	for {
		select {
		case <-Test2:
			fmt.Println(1111)
			goto L
		default:
			time.Sleep(1 * time.Second)
			fmt.Println(2222)
		}
	}
L:
	fmt.Println(33333)

}

func SelectTest() {
	for {
		select {
		case <-Test1:
			fmt.Println("111111")
		case msg := <-Test2:
			fmt.Println("222:", msg)
			fmt.Println("22222")
			goto Loop
		default:
			fmt.Println("333333")
			time.Sleep(1 * time.Second)
		}
	}
Loop:
	fmt.Println("444444444")
}

func SendChannel() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("yes")
		Test2 <- "yes"
	}
}
