/*
#Time      :  2019/5/12 上午11:17 
#Author    :  chuangangshen@deepglint.com
#File      :  goCrontabTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/jiansoft/robin"
	"time"
	"fmt"
	"strconv"
	"reflect"
)

var CrontabFlag = true

func main() {
	CrontabTest(10)
	for CrontabFlag {
		time.Sleep(100 * time.Millisecond)
	}
	// strTest()
	// WatchEtcd()
}

func WatchEtcd()  {
	var event int
	fmt.Println(event)
}

func CrontabTest(second int)  {
	robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "every 1 Minutes")
}

func runCron(s string)  {
	fmt.Printf("I am %s CronTest %v\n", s, time.Now())
}

func strTest() {
	a := "00"
	b, err := strconv.Atoi(a)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reflect.TypeOf(b), reflect.ValueOf(b))
	}
}
