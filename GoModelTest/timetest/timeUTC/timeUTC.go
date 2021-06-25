/*
#Time      :  2019/1/2 下午5:09
#Author    :  chuangangshen@deepglint.com
#File      :  timeUTC.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	// timeUTC()
	TimeUtcTest()
	// countSum()
}

func countSum() {
	var (
		i   = 0
		sum = 0
	)
	for i <= 120 {
		sum += i
		i += 5
	}
	fmt.Println(sum, sum/60)
}

func timeUTC() {
	year := 2019
	t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(t)
	fmt.Print(t.Unix())
}

func TimeUtcTest() {
	t1 := time.Now().UnixNano()
	time.Sleep(1 * time.Second)
	t2 := time.Now().UnixNano()
	fmt.Println(t2/1000000, t1/1000000)
}
