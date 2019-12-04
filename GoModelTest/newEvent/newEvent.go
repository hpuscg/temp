/*
#Time      :  2019/12/3 1:39 PM 
#Author    :  chuangangshen@deepglint.com
#File      :  newEvent.go
#Software  :  GoLand
*/
package main

import (
	"time"
	"sync"
)

func main() {

}

// 接收到的libra发送的数据结构
type Event struct {
	SensorId  string   `json:"sensor_id"`
	TimeStamp uint64   `json:"time_stamp"`
	People    []Person `json:"people"`
}

// 人员数据结构
type Person struct {
	PeopleId        uint64    `json:"people_id"`
	PathPoint       PathPoint `json:"path_point"`
	DetectionStatus int8      `json:"detection_status"`
}

// 人员坐标
type PathPoint struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
}

// 用于计算的tracking数据结构
type Tracking struct {
	LastUpdateTime time.Time   `json:"last_update_time"`
	SensorId       string      `json:"sensor_id"`
	PeopleId       uint64      `json:"people_id"`
	IsGhost        bool        `json:"is_ghost"`
	Path           []PathPoint `json:"path"`
	TimeStamp      uint64      `json:"time_stamp"`
	StartTime      time.Time   `json:"start_time"`
}

// 人员池
type PeoplePool struct {
	sync.RWMutex
	P map[string]Tracking
}
