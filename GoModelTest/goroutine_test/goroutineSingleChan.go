package main

import (
	"fmt"
	"time"
)

func counter(out chan<- int)  {
	for x := 0; x <= 100; x++ {
		out <- x
	}
}

func squarer(out chan<- int, in <-chan int)  {
	for x := range in {
		out <- x * x
	}
}

func printer(in <-chan int)  {
	for x := range in {
		fmt.Println(x)
		time.Sleep(100 * time.Millisecond)
	}
}

func main()  {
	naturals := make(chan int)
	squares := make(chan int)
	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

// ["config/eventserver/disableevent"]="[320,329,540,541,542,543,550,551,552,553,740,741,742,743,1000,1001,1002,1003]"