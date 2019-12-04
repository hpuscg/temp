/*
#Time      :  2019/5/17 下午10:38 
#Author    :  chuangangshen@deepglint.com
#File      :  glogTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/glog"
	"flag"
)

func main() {
	flag.Parse()
	GlogTest()
}

func GlogTest() {
	glog.Warningln(".......")
	glog.Infoln("333333")
	glog.V(0).Infoln("8888888")
	glog.V(3).Infoln("9999999")
	glog.Errorln("err")
}
