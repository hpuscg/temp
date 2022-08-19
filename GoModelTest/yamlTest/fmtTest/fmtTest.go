/*
#Time      :  2018/12/11 下午4:03
#Author    :  chuangangshen@deepglint.com
#File      :  fmtTest.go
#Software  :  GoLand
*/
package main

import "fmt"

func main() {
	fmtTest()
}

func fmtTest() {
	boundQuery := "((%s > %v)||(%s < %v))"
	b := fmt.Sprintf(boundQuery+"&&(%s == %v)", ".Duration", 4, ".Duration", 1, ".Population", 1)
	fmt.Println(b)
}
