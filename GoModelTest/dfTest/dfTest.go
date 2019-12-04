/*
#Time      :  2019/3/20 上午10:59 
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
	GetFsFree()
}

func GetFsFree()  {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/Users/hpu_scg/gocode/src/temp", &fs)
	if err != nil {
		fmt.Println(err)
	}
	free := fs.Bfree * uint64(fs.Bsize)/1024/1024/1024
	fmt.Println(free)
}



