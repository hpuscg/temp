/*
#Time      :  2019/1/16 上午10:59
#Author    :  chuangangshen@deepglint.com
#File      :  getEvents.go
#Software  :  GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/deepglint/util/jsontime"
	"github.com/segmentio/nsq-go"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	Id                   bson.ObjectId       "_id" //`json:"_id,omitempty"` //auto generate, don't fill it
	StartTime            *jsontime.Timestamp //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	TimeLength           int                 //ms
	SliceLength          int                 //`bson:"-" json:"-"` //ms
	AlarmLevel           int                 //compute by backend, don't fill it
	EventType            int                 //see values in EventTypeXXX consts
	EventTypeProbability float32             //[0,1]
	PlanetId             string
	SceneId              string
	SensorId             string
	HotspotId            string //可选
	UserId               string
	PeopleId             string //需要给出全局唯一的ID
	PeopleNum            int
	Path                 []int32 `bson:",omitempty" json:",omitempty"` //mm, in sequence of [x1,y1,z1,x2,y2,z2...]
	DetectionStatus      []int   `bson:",omitempty" json:",omitempty"`
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
	UserData             interface{}            `json:",omitempty"`
}

var (
	LibraAddr string = "192.168.101.147"
	EventType int    = 0
)

func main() {
	/* readConfig()
	if "" == LibraAddr {
		LibraAddr = "10.147.40.19"
	} */
	// fmt.Println("sensor ip is :", LibraAddr)
	getNsqMessage()
}

// get event message from sensor nsq
func getNsqMessage() {
	fmt.Println("=======1======", LibraAddr)
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "vibo_events",
		Channel:     "abcbank4scg",
		Address:     LibraAddr + ":4150",
		MaxInFlight: 250,
	})
	fmt.Println("==========2==========")
	for msg := range consumer.Messages() {
		fmt.Println("=======3=========")
		var event Event
		// fmt.Println("connect to nsq success")
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			fmt.Printf("Data Unmarshal error:%s\n", err.Error())
			continue
		}
		if 0 == EventType {
			fmt.Printf("\n%+v\n", event)
		} else if event.EventType == EventType {
			fmt.Printf("\n%+v\n", event)
		}
		// fmt.Printf("\n%+v\n", event)
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
			case "libraaddr":
				LibraAddr = strings.Split(line, "=")[1]
			case "eventtype":
				EventType, _ = strconv.Atoi(strings.Split(line, "=")[1])
			default:
				fmt.Printf("nothing get")
			}
		}
	}
	f.Close()
	fmt.Println("libra ip is:", LibraAddr)
}
