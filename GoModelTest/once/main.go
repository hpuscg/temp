package main

import (
	"flag"
	"fmt"
	"sync"
)

var (
	once sync.Once
	str  string
)

func main() {
	flag.StringVar(&str, "str", "default", "test flag")
	flag.Parse()
	for i := 0; i < 10; i++ {
		once.Do(onced)
		fmt.Println(i)
	}
}

func onced() {
	fmt.Println("onced")
}
