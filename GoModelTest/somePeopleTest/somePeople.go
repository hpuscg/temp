package main

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/deepglint/util/jsontime"
	"github.com/segmentio/nsq-go"
	"encoding/json"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"os"
	"io"
	"time"
	"net/http"
	"io/ioutil"
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
}

func main() {
	readConfig()
	go getNsqMessage()
	somePeopleEvent()
	time.Sleep(100 * time.Second)
}

// var msgChan chan nsq.Message
var (
	eventBox             = make([]*Event, 0)
	PeopleTotalNum             int
	PeopleInNum                int
	PeopleOutNum               int
	PeopleLimitNum             int
	LimitInTime          int64
	LimitOutTime         int64
	LibraAddr            string // 192.168.19.251
	PeopleLineIn               string
	PeopleLineOut              string
	ViboEventTopicServerAddr string
	topic                string
)

const (
	EventTypeComeInAndOut = 1020 // 同进同出
)

// get event message from sensor nsq
func getNsqMessage() {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "events",
		Channel:     "bank4scg",
		Address:     LibraAddr + ":4150",
		MaxInFlight: 250,
	})
	for msg := range consumer.Messages() {
		var event Event
		// fmt.Println("connect to nsq success")
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			fmt.Printf("Data Unmarshal error:%s", err.Error())
			continue
		}
		// append over line event to eventBox
		if event.EventType == 721 {
			eventBox = append(eventBox, &event)
			// fmt.Println("line num is: ", event.HotspotId)
		}
		msg.Finish()
	}
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
		respBody,_ := ioutil.ReadAll(resp.Body)
		fmt.Println("response data is: ", string(respBody))
	}
}

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
		// fmt.Println("err is :", err)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		line := strings.TrimSpace(string(b))
		// fmt.Println(line)
		if strings.Contains(line, "=") {
			caseStr := strings.Split(line, "=")[0]
			switch caseStr {
			case "peoplenum":
				PeopleLimitNum, err = strconv.Atoi(strings.Split(line, "=")[1])
			case "timelimitin":
				LimitInTime, err = strconv.ParseInt(strings.Split(line, "=")[1], 10, 64)
			case "timelimitout":
				LimitOutTime, err = strconv.ParseInt(strings.Split(line, "=")[1], 10, 64)
			case "libraaddr":
				LibraAddr = strings.Split(line, "=")[1]
			case "linein":
				PeopleLineIn = "door_" + strings.Split(line, "=")[1]
			case "lineout":
				PeopleLineOut = "door_" + strings.Split(line, "=")[1]
			case "viboserver":
				ViboEventTopicServerAddr = strings.Split(line, "=")[1]
			}
		}
	}
	f.Close()
	if PeopleLimitNum == 0 {
		PeopleLimitNum = 2
	}
	if LimitInTime == 0 {
		LimitInTime = 10
	}
	if LimitOutTime == 0 {
		LimitOutTime = 10
	}
	if PeopleLineIn == "" {
		PeopleLineIn = "1"
	}
	if PeopleLineOut == "" {
		PeopleLineOut = "2"
	}
	fmt.Println("libra ip is:", LibraAddr)
	fmt.Println("in time limit is :", LimitInTime)
	fmt.Println("out time limit is :", LimitOutTime)
	fmt.Println("limit people num is :", PeopleLimitNum)
	fmt.Println("in line is :", PeopleLineIn)
	fmt.Println("out line is :", PeopleLineOut)
	fmt.Println("vibo event server addr is :", ViboEventTopicServerAddr)
}

// 实际处理触发的事件
func somePeopleEvent() {
	for true {
		var event *Event
		if len(eventBox) > 0 {
			event = eventBox[0]
			eventBox = eventBox[1:]
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		lineNum := event.HotspotId
		// EventType := event.EventType
		if "" != lineNum {
			countPeopleNum(lineNum)
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if PeopleTotalNum < PeopleLimitNum && PeopleTotalNum > 0 {
			startTime := time.Now().Unix()
			for true {
				if PeopleTotalNum >= PeopleLimitNum || PeopleTotalNum <= 0 {
					break
				}
				if lineNum == PeopleLineIn {
					endTimeIn := time.Now().Unix()
					useTimeIn := endTimeIn - startTime
					if useTimeIn > LimitInTime {
						event.EventType = EventTypeComeInAndOut
						postMessage(event)
						fmt.Printf("the people num is: %d, come in use time upper %d s \n", PeopleTotalNum, LimitInTime)
						break
					}
					var eventIn *Event
					if len(eventBox) > 0 {
						eventIn = eventBox[0]
						eventBox = eventBox[1:]
						startTime = time.Now().Unix()
					} else {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					lineNum = eventIn.HotspotId
					// EventType = eventIn.EventType
					countPeopleNum(lineNum)
				} else if lineNum == PeopleLineOut {
					endTimeOut := time.Now().Unix()
					useTimeOut := endTimeOut - startTime
					if useTimeOut > LimitOutTime {
						event.EventType = EventTypeComeInAndOut
						postMessage(event)
						fmt.Printf("the home people num is: %d, lower %d\n", PeopleTotalNum, PeopleLimitNum)
						break
					}
					var eventOut *Event
					if len(eventBox) > 0 {
						eventOut = eventBox[0]
						eventBox = eventBox[1:]
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
		} else if PeopleTotalNum > PeopleLimitNum {
			event.EventType = EventTypeComeInAndOut
			postMessage(event)
			fmt.Printf("the home people num is: %d, upper %d\n", PeopleTotalNum, PeopleLimitNum)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// count the num of people
func countPeopleNum(line string) {
	// fmt.Printf("the line num is : %s\n", line)
	if line == PeopleLineIn {
		PeopleInNum += 1
		PeopleTotalNum += 1
	} else if line == PeopleLineOut {
		if PeopleTotalNum > 0 {
			PeopleOutNum += 1
			PeopleTotalNum -= 1
		} else {
			fmt.Println("total num is zero now")
		}
	} else {
		fmt.Println("Stop the program")
		fmt.Println("Please modify the linein and lineout value in config.txt file")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	fmt.Printf("total num is : %d,in num is : %d,out num is : %d\n", PeopleTotalNum, PeopleInNum, PeopleOutNum)
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
