/*
#Time      :  2019/3/26 上午10:07 
#Author    :  chuangangshen@deepglint.com
#File      :  dfTest.go
#Software  :  GoLand
*/
package main

import (
	"syscall"
	"fmt"
)

func main() {
	dfTest()
}

func dfTest() {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(int64(fs.Bfree * uint64(fs.Bsize)))
	fmt.Println(int64(fs.Blocks * uint64(fs.Bsize)))
}