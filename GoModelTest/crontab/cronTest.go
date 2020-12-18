/*
#Time      :  2020/10/22 11:37 上午
#Author    :  chuangangshen@deepglint.com
#File      :  cronTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	// SingletonTest()
	appMode := os.Getenv("PATH")
	fmt.Println("app mode is :", appMode)
}

func SingletonTest() {
	ins1 := GetInstance()
	ins2 := GetInstance()
	fmt.Println(ins1 == ins2)
	fmt.Println(&ins2, &ins1)
}

var task = func() {
	fmt.Println("hello world")
}

type Singleton struct {
}

var (
	singleton *Singleton
	once      sync.Once
)

func GetInstance() *Singleton {
	once.Do(func() {
		singleton = &Singleton{}
	})
	return singleton
}
