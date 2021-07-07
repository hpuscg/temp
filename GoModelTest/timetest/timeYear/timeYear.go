package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	timeYear()
}

func timeYear() {
	fmt.Println(time.Now().Year())
	if time.Now().Year() < 2021 {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}
}

func timeYearTest() {
	end := time.Now().UTC()
	// endYear := end.Year()
	endYearDay := end.YearDay()
	/*endDay := endYear*1000 + endYearDay
	begin := end.AddDate(0, 0, -3)
	oldest := time.Date(endDay/1000, 1, 1, 0, 0, 0, 0, time.UTC)
	oldestReset := oldest.AddDate(0, 0, endDay%1000-1)
	jsonOldest := jsontime.Timestamp(oldest)*/
	fmt.Println(endYearDay)
	a := 201907204
	fmt.Println(a / 100000)
	fmt.Println(a % 100000 / 100)
	tm := time.Unix(1552435200, 0)
	fmt.Println(tm)
	tn := tm.Unix()
	fmt.Println(tn)

	/*fmt.Println("end is : ", end)
	fmt.Println("endYear is : ", endYear)
	fmt.Println("endYearDay is : ", endYearDay)
	fmt.Println("endDay is : ", endDay)
	fmt.Println("begin is : ", begin)
	fmt.Println("oldest is : ", oldest)
	fmt.Println("oldestReset is : ", oldestReset)
	fmt.Println("jsonoldest is : ", jsonOldest)
	time_to := time.Now().UTC().Add(-time.Hour * time.Duration(24*1))
	fmt.Println(time_to)*/
	/*
		timestamp := end.Unix()
		var strtime string
		strtime = strconv.FormatInt(timestamp, 10)
		fmt.Println(timestamp)
		fmt.Println(strtime)
		tm := time.Unix(timestamp, 0).Format("2006-01-02") + " 00:00:01"
		// tm2 := tm.YearDay()
		tm2,_ := time.Parse("2006-01-02 15:04:05", tm)
		fmt.Println(tm2.Unix())

		timelimit := "2018-07-16 08:50:47 +0000 UTC"
		// timelimit2 := strconv.FormatInt(timelimit, 10)
		timefor, _ := time.Parse("2006-01-02 15:04:05", timelimit)
		fmt.Println(timefor)
		// var tmt string
		tmt := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
		fmt.Println(tmt)
	*/
	/*
		num := 2
		numstr := strconv.Itoa(num)
		fmt.Println(numstr)


		tm1 := time.Now().Format("2006-01-02")
		fmt.Println(tm1)
		time.Sleep(time.Second)
		tm2 := time.Now().Format("2006-01-02")
		if tm1 == tm2 {
			fmt.Println(tm2)
		} else {
			fmt.Println("no")
		}
	*/

	/*
		tim1 := time.Now().Format("2006-01-02")
		fmt.Println(tim1)
		tim, _ := time.Parse("2006-01-02 15:04:05", tim1 + " 00:00:00")


		timOldYear := tim.AddDate(-1, 0, 0)
		fmt.Println(timOldYear)
		tim3 := time.Now().Format("2006")

		tim2, _ := time.Parse("2006-01-02", tim3 + "-12-31")

		tim4 := tim2.AddDate(-2, 0, 0)
		maxCount := tim4.YearDay()
		fmt.Println(maxCount)

		timOldStamp := strconv.FormatInt(timOldYear.Unix(), 10)

		typetim := reflect.TypeOf(timOldStamp)

		fmt.Println(timOldStamp)

		fmt.Println(typetim)
	*/
	// oldDay()
	// timeSecond()
	// timeMinute()
	// timeTest()
	// utcTest()
}

func utcTest() {
	t0 := time.Now()
	t1 := time.Now().Format("2006-01-02 15:04:05")
	t2, _ := time.Parse("2006-01-02 15:04:05", t1)
	t3 := time.Now().Unix()
	t4 := time.Now().UTC()
	t5 := time.Unix(t3, 0)
	fmt.Println(t0)
	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)
	fmt.Println(t5)
}

func timeTest() {
	time1 := time.Now().Format("2006-01-02 15:04:05")
	time2, _ := time.Parse("2006-01-02 15:04:05", time1)
	time3 := time2.Add(-8 * 3600 * 1000000000)
	time4 := time3.Unix()
	time5 := time.Unix(time4, 0)
	time6 := time5.Format("2006-01-02 15:04:05")
	time7, _ := time.Parse("2006-01-02 15:04:05", time6)
	fmt.Println(time1)
	fmt.Println(time2)
	fmt.Println(time3)
	fmt.Println(time4)
	fmt.Println(time5)
	fmt.Println(time6)
	fmt.Println(time7)
}

func oldDay() {
	count := 0
	// 获取去年的天数
	nowTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02"+" 00:00:00"))
	countDay, _ := time.Parse("2006-01-02", time.Now().Format("2006")+"-12-31")
	maxCount := countDay.AddDate(-1, 0, 0).YearDay()
	for count < maxCount {
		oldTime := strconv.FormatInt(nowTime.AddDate(-1, 0, -count).Unix(), 10)
		fmt.Println(oldTime)
		count++
		// time.Sleep(2 * time.Second)
	}
}

func timeSecond() {
	timeS := time.Now().Unix()
	fmt.Println(timeS)
	timeMs := time.Now().UnixNano() / 1000000
	fmt.Println(timeMs)
	timeNs := time.Now().UnixNano()
	fmt.Println(timeNs)
}

func timeMinute() {
	timeM := time.Now().Format("2006-01-02 15")
	fmt.Println(timeM)
	time2, _ := time.Parse("2006-01-02 15:04:05", timeM+":00:00")
	fmt.Println(time2)
	time3 := time2.Unix()
	fmt.Println(time3)
}
