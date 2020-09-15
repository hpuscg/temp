package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	// timeTest()
	//NanoTimeTest()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sigChan:
			fmt.Println("single")
		case <-t.C:
			fmt.Println("yes")
		}
	}
	fmt.Println(time.Now().Unix())
}

func NanoTimeTest() {
	t1 := time.Now().UnixNano()
	time.Sleep(1 * time.Second)
	t2 := time.Now().UnixNano()
	fmt.Println(t2-t1)
	fmt.Println((t2-t1)/1e9)
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

