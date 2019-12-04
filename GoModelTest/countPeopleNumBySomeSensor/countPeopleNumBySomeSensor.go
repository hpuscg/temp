package main

import (
	"github.com/deepglint/util/jsontime"
	"gopkg.in/mgo.v2/bson"
	"github.com/segmentio/nsq-go"
	"fmt"
	"encoding/json"
	"os"
	"bufio"
	"strings"
	"io"
	"time"
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
	LibraIp              string                 `json:"libra_ip"`
}

func main() {
	readConfig()
	for _, libraInfo := range LibraConfig {
		go getNsqMessage(libraInfo.LibraAddr)
	}
	handleEvent()
}

// people num info
type PeopleNum struct {
	TotalNum int
	InNum    int
	OutNum   int
}

// libra info
type LibraInfo struct {
	LibraAddr string `json:"libra_addr"`
	LineIn    string `json:"line_in"`
	LineOut   string `json:"line_out"`
}

var (
	LibraConfig []LibraInfo
	peopleNum   PeopleNum
	eventBox = make([]Event, 5)
)

// get event message from sensor nsq
func getNsqMessage(LibraAddr string) {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "events",
		Channel:     "bank4scg00000",
		Address:     LibraAddr + ":4150",
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
			event.LibraIp = LibraAddr
			eventBox = append(eventBox, event)
			fmt.Printf("===94===%+v\n", event)
		}
		msg.Finish()
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
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		var libraInfo LibraInfo
		line := strings.TrimSpace(string(b))
		if strings.Contains(line, ";") && strings.HasPrefix(line, "libra") {
			sub := strings.Split(line, ";")
			if len(sub) == 3 {
				for _, value := range sub {
					caseStr := strings.Split(value, "=")[0]
					switch caseStr {
					case "libraaddr":
						libraInfo.LibraAddr = strings.Split(value, "=")[1]
					case "linein":
						libraInfo.LineIn = "door_" + strings.Split(value, "=")[1]
					case "lineout":
						libraInfo.LineOut = "door_" + strings.Split(value, "=")[1]
					}
				}
			}
			fmt.Println("===148===", libraInfo)
			LibraConfig = append(LibraConfig, libraInfo)
		} else {
			sub := strings.Split(line, ";")
			if len(sub) == 3 {
				for _, value := range sub {
					caseStr := strings.Split(value, "=")[0]
					switch caseStr {
					case "totalnum":
						peopleNum.TotalNum, _ = strconv.Atoi(strings.Split(value, "=")[1])
					case "innum":
						peopleNum.InNum, _ = strconv.Atoi(strings.Split(value, "=")[1])
					case "outnum":
						peopleNum.OutNum, _ = strconv.Atoi(strings.Split(value, "=")[1])
					}
				}
			}
		}
	}
	f.Close()
}

func handleEvent() {
	var LineIn, LineOut string
	for true {
		var event Event
		if len(eventBox) > 0 {
			event = eventBox[0]
			eventBox = eventBox[1:]
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		libraIp := event.LibraIp
		for _, libraInfo := range LibraConfig {
			if libraIp == libraInfo.LibraAddr {
				LineOut = libraInfo.LineOut
				LineIn = libraInfo.LineIn
			}
		}
		lineNum := event.HotspotId
		// EventType := event.EventType
		if "" != lineNum && 721 == event.EventType {
			if lineNum == LineIn {
				peopleNum.InNum += 1
				peopleNum.TotalNum += 1
			} else if lineNum == LineOut {
				if peopleNum.TotalNum > 0 {
					peopleNum.OutNum += 1
					peopleNum.TotalNum -= 1
				}
			} else {
				fmt.Println("Stop the program")
				fmt.Println("Please modify the linein and lineout value in config.txt file")
				time.Sleep(3 * time.Second)
				os.Exit(1)
			}
			fmt.Printf("people num is : %+v\n", peopleNum)
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
}
