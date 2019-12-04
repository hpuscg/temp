package main

import (
	"fmt"
	"time"
)

func main()  {
	// noBuffer()
	// useBuffer()
	channelLength()
	time.Sleep(5 * time.Second)
	fmt.Println("ok")
}

func f1(in chan int)  {
	for {
		fmt.Println(<-in)
		fmt.Println("continue")
		time.Sleep(1 * time.Second)
	}
}

func noBuffer()  {
	out := make(chan int)
	go f1(out)
	out <- 2
}


func useBuffer()  {
	out := make(chan int, 1)
	out <- 3
	go f1(out)
}

func channelLength()  {
	e := make(chan int, 10)
	fmt.Println(len(e))
}

