/*
#Time      :  2019/1/11 下午2:06 
#Author    :  chuangangshen@deepglint.com
#File      :  singlePeopleIn.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"strings"
	"os"
	"fmt"
	"io"
	"gopkg.in/mgo.v2/bson"
	"github.com/deepglint/util/jsontime"
	"github.com/segmentio/nsq-go"
	"encoding/json"
	"time"
	"net/http"
	"io/ioutil"
	"strconv"
)

type Event struct {
	Id                   bson.ObjectId          "_id" //`json:"_id,omitempty"` //auto generate, don't fill it
	StartTime            *jsontime.Timestamp          //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	TimeLength           int                          //ms
	SliceLength          int                          //`bson:"-" json:"-"` //ms
	AlarmLevel           int                          //compute by backend, don't fill it
	EventType            int                          //see values in EventTypeXXX consts
	EventTypeProbability float32                      //[0,1]
	PlanetId             string
	SceneId              string
	SensorId             string
	HotspotId            string //可选
	UserId               string
	PeopleId             string //需要给出全局唯一的ID
	PeopleNum            int
	Path                 []int32                `bson:",omitempty" json:",omitempty"` //mm, in sequence of [x1,y1,z1,x2,y2,z2...]
	DetectionStatus      []int                  `bson:",omitempty" json:",omitempty"`
	FrameRate            int
	ColorPanel16x16      []float32              `bson:",omitempty" json:",omitempty"` //[0..1], ab space, splite by 16*16
	LightPanel16         []float32              `bson:",omitempty" json:",omitempty"` // L space, splite by 16
	Color                []float32              `bson:",omitempty" json:",omitempty"` //in Lab format, in sequence of [c1L,c1a,c1b,c2L,c2a,c2b,...]. Generate by backend, don't fill it
	Height               int                    `bson:",omitempty" json:",omitempty"` //cm
	PlaneArea            int                    `bson:",omitempty" json:",omitempty"` //cm*cm
	CutboardBox          []int                  `bson:",omitempty" json:",omitempty"`
	CutboardTimeOffset   []int                  `bson:",omitempty" json:",omitempty"`
	CutboardTime         *jsontime.Timestamp    `bson:",omitempty" json:",omitempty"`
	Payload              map[string]interface{} `bson:",omitempty" json:",omitempty"`
	// LibraIp              string
}

func main() {
	readConfig()
	go getDoorEvent()
	for _, latchAddr := range LatchAddrs {
		go getLatchEvent(latchAddr)
	}
	go singleMoney()
	singlePeople()
}

var (
	LineEventBox             = make([]*Event, 5)
	SinglePeopleBox          = make([]*Event, 5)
	LimitNum                 int
	TotalNum                 int
	InNum                    int
	OutNum                   int
	SinglePeopleTime         int64
	CountNumAddr             string
	LatchAddrs               []string
	LineIn                   string
	LineOut                  string
	topic                    string
	ViboEventTopicServerAddr string
	SingleMoneyTime          int64
	SensorIds                []string
)

const (
	EventTypeSinglePeopleIn = 1030 // 单人进入
	EventTypeSingleMoney    = 1040 // 单人加钞
	EventTypeTooLong        = 1050 // 逗留过久
)

// read config data from config file
func readConfig() {
	f, err := os.Open("config.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		line := strings.TrimSpace(string(b))
		if strings.Contains(line, "=") {
			caseStr := strings.Split(line, "=")[0]
			switch caseStr {
			case "dooraddr":
				fmt.Println("===110===", strings.Split(line, ";")[0])
				CountNumAddr = strings.Split(strings.Split(line, ";")[0], "=")[1]
			case "linein":
				LineIn = "door_" + strings.Split(line, "=")[1]
			case "lineout":
				LineOut = "door_" + strings.Split(line, "=")[1]
			case "latchaddr":
				LatchAddrs = append(LatchAddrs, strings.Split(strings.Split(line, ";")[0], "=")[1])
				SensorIds = append(SensorIds, strings.Split(strings.Split(line, ";")[1], "=")[1])
			case "limitnum":
				LimitNum, _ = strconv.Atoi(strings.Split(line, "=")[1])
			case "viboserver":
				ViboEventTopicServerAddr = strings.Split(line, "=")[1]
			case "singlepeopletime":
				SinglePeopleTime, _ = strconv.ParseInt(strings.Split(line, "=")[1], 10, 64)
			case "singlemoneytime":
				SingleMoneyTime, _ = strconv.ParseInt(strings.Split(line, "=")[1], 10, 64)
			}
		}
	}
	f.Close()
}

// post message by http
func postMessage(event *Event) {
	if "vibo_events" != topic {
		topic = "vibo_events"
	}
	ei := eventToInterface(event, "json", "ms", "simple")
	postData, _ := json.Marshal(ei)
	data := strings.NewReader(string(postData))
	url := "http://" + ViboEventTopicServerAddr + ":4151/pub?topic=" + topic
	request, _ := http.NewRequest("POST", url, data)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("post data err: ", err)
	} else {
		fmt.Println("post data sussess")
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response data is: ", string(respBody))
	}
}

// 获取越线信息
func getDoorEvent() {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "events",
		Channel:     "bank4door",
		Address:     CountNumAddr + ":4150",
		MaxInFlight: 250,
	})
	for msg := range consumer.Messages() {
		var event Event
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			fmt.Printf("Data Unmarshal error:%s", err.Error())
			continue
		}
		// append over line event to eventBox
		if event.EventType == 721 {
			// event.LibraIp = CountNumAddr
			LineEventBox = append(LineEventBox, &event)
			fmt.Println("door event: ", event)
		}
		msg.Finish()
	}
	fmt.Println("===176===")
}

// 获取单人加钞信息
func getLatchEvent(LatchAddr string) {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "events",
		Channel:     "bank4latch",
		Address:     LatchAddr + ":4150",
		MaxInFlight: 250,
	})
	for msg := range consumer.Messages() {
		var event Event
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			fmt.Printf("Data Unmarshal error:%s", err.Error())
			continue
		}
		// append over line event to eventBox
		if 410 == event.EventType || 411 == event.EventType ||
			412 == event.EventType || 413 == event.EventType {
			// event.LibraIp = LatchAddr
			SinglePeopleBox = append(SinglePeopleBox, &event)
			fmt.Println("latch event: ", event)
		}
		msg.Finish()
	}
}

// 根据越线数人
// count the num of people
func countPeopleNum(line string) {
	// fmt.Printf("the line num is : %s\n", line)
	if line == LineIn {
		InNum += 1
		TotalNum += 1
	} else if line == LineOut {
		if TotalNum > 0 {
			OutNum += 1
			TotalNum -= 1
		} else {
			fmt.Println("total num is zero now")
		}
	} else {
		fmt.Println("Stop the program")
		fmt.Println("Please modify the linein and lineout value in config.txt file")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	fmt.Printf("total num is : %d,in num is : %d,out num is : %d\n", TotalNum, InNum, OutNum)
}

// 单人进入和逗留过久
func singlePeople() {
	for true {
		var doorEvent *Event
		if len(LineEventBox) > 0 {
			doorEvent = LineEventBox[0]
			LineEventBox = LineEventBox[1:]
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if doorEvent == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		lineNum := doorEvent.HotspotId
		if "" != lineNum {
			countPeopleNum(lineNum)
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if TotalNum > 0 && TotalNum < LimitNum {
			startTime := time.Now().Unix()
			for true {
				if TotalNum >= LimitNum || TotalNum <= 0 {
					break
				}
				if lineNum == LineIn {
					endTime := time.Now().Unix()
					useTime := endTime - startTime
					if useTime > SinglePeopleTime {
						// single people in
						// 单人进入
						doorEvent.EventType = EventTypeSinglePeopleIn
						postMessage(doorEvent)
						fmt.Printf("the people num is: %d, come in use time upper %d s \n", TotalNum, SinglePeopleTime)
					}
					var eventIn *Event
					if len(LineEventBox) > 0 {
						eventIn = LineEventBox[0]
						LineEventBox = LineEventBox[1:]
						startTime = time.Now().Unix()
					} else {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					lineNum = eventIn.HotspotId
					// EventType = eventIn.EventType
					countPeopleNum(lineNum)
				} else if lineNum == LineOut {
					endTime := time.Now().Unix()
					useTime := endTime - startTime
					if useTime > SinglePeopleTime {
						// some people state too long
						// 逗留过久
						doorEvent.EventType = EventTypeTooLong
						postMessage(doorEvent)
						fmt.Printf("the people num is: %d, come out use time upper %d s \n", TotalNum, SinglePeopleTime)
					}
					var eventOut *Event
					if len(LineEventBox) > 0 {
						eventOut = LineEventBox[0]
						LineEventBox = LineEventBox[1:]
						startTime = time.Now().Unix()
					} else {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					lineNum = eventOut.HotspotId
					// EventType = eventOut.EventType
					countPeopleNum(lineNum)
				}
			}
		}
	}
}

// 单人加钞
func singleMoney() {
	for true {
		var singleMoneyEvent *Event
		if len(SinglePeopleBox) > 0 {
			singleMoneyEvent = SinglePeopleBox[0]
			SinglePeopleBox = SinglePeopleBox[1:]
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if singleMoneyEvent == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		doorEventType := singleMoneyEvent.EventType
		startTime := time.Now().Unix()
		for true {
			if (411 != doorEventType && 412 != doorEventType) || 1 != TotalNum {
				break
			}
			endTime := time.Now().Unix()
			useTime := endTime - startTime
			if useTime > SingleMoneyTime {
				singleMoneyEvent.EventType = EventTypeSingleMoney
				postMessage(singleMoneyEvent)
				fmt.Printf("the people num is: %d, use time upper %d s \n", TotalNum, SingleMoneyTime)
			}
			var tempEvent *Event
			if len(SinglePeopleBox) > 0 {
				tempEvent = SinglePeopleBox[0]
				SinglePeopleBox = SinglePeopleBox[1:]
				if 411 != tempEvent.EventType && 412 != tempEvent.EventType {
					break
				}
			} else {
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
	}
}

// event struct ro interface
func eventToInterface(e *Event, format string, timeformat string, content string) interface{} {
	var output interface{}
	if content == "" {
		content = "default"
	}
	if timeformat == "" {
		timeformat = "ms"
	}
	var formattedtime interface{}
	if timeformat == "s" {
		formattedtime = e.StartTime.Unix()
	} else if timeformat == "ms" {
		formattedtime = e.StartTime.UnixMilli()
	} else if timeformat == "ns" {
		formattedtime = e.StartTime.UnixNano()
	} else if timeformat == "string" {
		formattedtime = e.StartTime.Format("2006-01-02T15:04:05.000MST")
	} else {
		formattedtime = e.StartTime.Format(timeformat)
	}
	if format == "line" {
		output = fmt.Sprintf("%v\n%d\n",
			formattedtime, e.TimeLength)
	} else {
		output = map[string]interface{}{
			"Id":                   e.Id,
			"StartTime":            formattedtime,
			"TimeLength":           e.TimeLength,
			"AlarmLevel":           e.AlarmLevel,
			"EventType":            e.EventType,
			"EventTypeProbability": e.EventTypeProbability,
			"PlanetId":             e.PlanetId,
			"SceneId":              e.SceneId,
			"SensorId":             e.SensorId,
			"UserId":               e.UserId,
			"HotspotId":            e.HotspotId,
			"PeopleId":             e.PeopleId,
		}
		if content == "complete" {
			// output.(map[string]interface{})["Path"] = e.Path
			// output.(map[string]interface{})["DetectionStatus"] = e.DetectionStatus
			output.(map[string]interface{})["FrameRate"] = e.FrameRate
			output.(map[string]interface{})["CutboardBox"] = e.CutboardBox
			output.(map[string]interface{})["CutboardTimeOffset"] = e.CutboardTimeOffset
			//output.(map[string]interface{})["Payload"] = e.Payload
			output.(map[string]interface{})["SliceLength"] = e.SliceLength
		}
		if content == "crowd_analysis" {
			output.(map[string]interface{})["Path"] = e.Path
			output.(map[string]interface{})["DetectionStatus"] = e.DetectionStatus
			output.(map[string]interface{})["FrameRate"] = e.FrameRate
			output.(map[string]interface{})["CutboardBox"] = e.CutboardBox
			output.(map[string]interface{})["CutboardTimeOffset"] = e.CutboardTimeOffset
			//output.(map[string]interface{})["Payload"] = e.Payload
			output.(map[string]interface{})["SliceLength"] = e.SliceLength
		}
	}
	return output
}
