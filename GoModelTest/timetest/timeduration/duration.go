package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func main() {
	// fmt.Println(time.Duration(6666666600000))
	// a := 20
	// TimeSwitch(20)
	// DurationSwitchTest()
	// timeString()
	TimeNow()
}

func TimeNow() {
	fmt.Println(time.Now().Unix())
}

func timeString() {
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10) + "")
}

func TimeSwitch(t time.Duration) {

	b := int64(t / time.Second)
	fmt.Println(reflect.TypeOf(b))
}

func DurationSwitchTest() {
	fmt.Println("1111")
	time.Sleep(time.Duration(5 * int64(time.Second)))
	fmt.Println("2222")
}
