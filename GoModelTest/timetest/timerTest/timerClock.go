package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	// timerClock()
	// timeWeek()
	// nanoTest()
	TimeStamp()
}

func TimeStamp() {
	fmt.Println(strconv.Itoa(int(time.Now().UnixNano())))
}

func nanoTest() {
	t1 := time.Now().UnixNano() / 1000000
	time.Sleep(1 * time.Second)
	t2 := time.Now().UnixNano() / 1000000
	fmt.Println(t2 - t1)
}

func timeWeek() {
	t := time.Unix(200, 0)
	w := t.Weekday()
	fmt.Println(w)
}

func timerClock() {

	timer1 := make(chan bool, 2)
	go func() {
		num1 := 4
		for num1 > -1 {
			if num1 > 0 {
				timer1 <- false
			}
			num1 -= 1
			time.Sleep(1 * time.Second)
		}
		timer1 <- true
		return
	}()
	time1 := time.Now().Unix()
	for true {

		var num2 bool
		num2 = <-timer1
		if num2 {
			fmt.Println("yes")
			break
		} else {
			fmt.Println("no")
		}
		time.Sleep(100)
	}
	time2 := time.Now().Unix()
	time3 := time2 - time1
	fmt.Println(time3)
}
