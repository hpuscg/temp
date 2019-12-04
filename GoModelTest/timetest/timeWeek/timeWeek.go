/*
#Time      :  2019/8/1 上午11:03 
#Author    :  chuangangshen@deepglint.com
#File      :  timeWeek.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/util/jsontime"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// GetWeekFromTimeStamp()
	fmt.Println(time.Now().Unix())
}

func GetWeekFromTimeStamp() {
	var startTime *jsontime.Timestamp
	// timeStamp := 1564628988
	data, err := json.Marshal(time.Now())
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &startTime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(startTime)
	fmt.Println(startTime.UTC())
	fmt.Println(time.Now().UTC().Weekday())
}

