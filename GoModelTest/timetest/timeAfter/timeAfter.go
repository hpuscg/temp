package main

import (
	"fmt"
	"time"
)

func main() {
	/*tc := time.After(5*time.Second)
	i := 1
	for i > 0 {
		time.Sleep(2*time.Second)
		i++
		fmt.Println(i)
	}
	<-tc
	fmt.Println("close")
	time.Now().UTC()*/
	timeTest()
}

func timeTest() {
	useTime := "4:12:14"
	todayRebootTime := time.Now().Format("2006-01-02 ") + useTime
	// fmt.Println(time.Now().Format(todayRebootTime), todayRebootTime)
	rebootTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayRebootTime, time.Local)
	nowTime := time.Now().UTC().Add(8 * time.Hour)
	fmt.Println(rebootTime, time.Now().Unix(), time.Now().UTC().Unix())
	if nowTime.After(rebootTime) {
		rebootTime = rebootTime.Add(24 * time.Hour)
	}
	fmt.Println(rebootTime, time.Now().After(rebootTime), time.Now().Before(rebootTime))
	// fmt.Println(nowStr)
	// fmt.Println(time.Parse("2006-01-02 15:04:05", nowStr))
	// rebootTime = rebootTime.AddDate(time.Now().Year(), int(time.Now().Month()), time.Now().Day())
	// fmt.Println(rebootTime)
}

