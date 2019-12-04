package main

import (
	"github.com/segmentio/nsq-go"
	"fmt"
	"encoding/json"
	"github.com/deepglint/util/jsontime"
	"gopkg.in/mgo.v2/bson"
	"bufio"
	"strings"
	"os"
	"io"
	"github.com/goburrow/modbus"
	"time"
)

func main() {
	fmt.Println("1111")
	readConfig()
	fmt.Println("2222")
	go getNsqMessage()
	fmt.Println("33333")
	eventHandler()
}

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

// var msgChan chan nsq.Message
var (
	eventBox  = make([]Event, 5)
	LibraAddr string // 192.168.19.251
	IOAddr    string // 192.168.7.11
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
		if event.EventType == 741 || event.EventType == 742 || event.EventType == 740 || event.EventType == 743 {
			eventBox = append(eventBox, event)
			fmt.Println(event)
			// _, err = client.WriteSingleCoil(100, 0xFF00)

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
		fmt.Println("err is :", err)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		line := strings.TrimSpace(string(b))
		fmt.Println(line)
		if strings.Contains(line, "=") {
			caseStr := strings.Split(line, "=")[0]
			if caseStr == "libraaddr" {
				LibraAddr = strings.Split(line, "=")[1]
			}
			if caseStr == "ioaddr" {
				IOAddr = strings.Split(line, "=")[1]
			}
		}
	}
	f.Close()

	if LibraAddr == "" {
		LibraAddr = "192.168.100.191"
	}
	if IOAddr == "" {
		IOAddr = "192.168.100.252"
	}
	fmt.Println("libra ip is:", LibraAddr)
	fmt.Println("io addr is :", IOAddr)
}

func eventHandler() {
	client := modbus.TCPClient(IOAddr + ":502")
	for true {
		var event Event
		if len(eventBox) > 0 {
			event = eventBox[0]
			eventBox = eventBox[1:]
		} else {
			// fmt.Println("no event")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// fmt.Println("event type", event.EventType)
		if 741 == event.EventType || 742 == event.EventType {
			fmt.Println(event.EventType)
			fmt.Println("========linght up")
			_, err := client.WriteSingleCoil(100, 0xFF00)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if 740 == event.EventType || 743 == event.EventType {
			fmt.Println(event.EventType)
			fmt.Println("***********light down")
			_, err := client.WriteSingleCoil(100, 0x0000)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
