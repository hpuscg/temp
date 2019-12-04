package main

import (
	"fmt"
	"time"
)

func main() {
	tc := time.After(5*time.Second)
	i := 1
	for i > 0 {
		time.Sleep(2*time.Second)
		i++
		fmt.Println(i)
	}
	<-tc
	fmt.Println("close")
	time.Now().UTC()
}

