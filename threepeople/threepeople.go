package main


import (
	"strings"
	"os"
	"bufio"
	"fmt"
	"time"
	"io"
	"strconv"
	"github.com/segmentio/nsq-go"
	"github.com/goburrow/modbus"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/deepglint/muses/util/jsontime"
)

type Event struct {
	Id bson.ObjectId "_id" //`json:"_id,omitempty"` //auto generate, don't fill it
	StartTime *jsontime.Timestamp //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	//OriginStartTime *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event start time if the event is splitted
	//OriginEndTime   *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event end time if the event is splitted
	TimeLength  int //ms
	SliceLength int //`bson:"-" json:"-"` //ms
	AlarmLevel           int     //compute by backend, don't fill it
	EventType            int     //see values in EventTypeXXX consts
	EventTypeProbability float32 //[0,1]
	PlanetId  string
	SceneId   string
	SensorId  string
	HotspotId string //可选
	UserId string
	PeopleId string //需要给出全局唯一的ID
	PeopleNum int
	Path []int32 `bson:",omitempty" json:",omitempty"` //mm, in sequence of [x1,y1,z1,x2,y2,z2...]
	// DetectionStatus indicates the detection confidence of the current path.
	// 1 is real, and 0 is ghost.
	DetectionStatus []int `bson:",omitempty" json:",omitempty"`

	// FrameRate indicates the FPS of the trajectory points.
	FrameRate int
	//TouchPoint []int `json:",omitempty"` //mm, in sequence of [x, y]

	//这组目前得到的数据不准，可以不填
	ColorPanel16x16 []float32 `bson:",omitempty" json:",omitempty"` //[0..1], ab space, splite by 16*16
	LightPanel16    []float32 `bson:",omitempty" json:",omitempty"` // L space, splite by 16
	Color           []float32 `bson:",omitempty" json:",omitempty"` //in Lab format, in sequence of [c1L,c1a,c1b,c2L,c2a,c2b,...]. Generate by backend, don't fill it
	Height          int       `bson:",omitempty" json:",omitempty"` //cm
	PlaneArea       int       `bson:",omitempty" json:",omitempty"` //cm*cm

	//only event of normal flow has a cutboard
	//PicFId    string //auto generate, don't fill it
	//PicBinary string `json:",omitempty"`

	// [x,y,height,width] of the cutboard
	CutboardBox []int `bson:",omitempty" json:",omitempty"`
	// the offset in millisecond of the cutboard time from StartTime
	CutboardTimeOffset []int `bson:",omitempty" json:",omitempty"`
	CutboardTime *jsontime.Timestamp `bson:",omitempty" json:",omitempty"`
	// Payload can store various information like ServerTimestamp, etc.
	// Debugging data:
	// "ServerTimestamp" stores the timestamp based on server time.
	// "Msg" shows the detailed readable debug message.
	//
	// Abnormal detection data:
	// "Value" stores the feature value.
	// "Score" stores the abnormal store for this event.
	// "ExceedsUpperBound" indicates whether the feature value exceeds upper bound.
	// "ExceedsLowerBound" indicates whether the feature value is under the lower bound.
	// "Position" is the current position for this person. <optional>
	// "Area" is the active area of this person. <optional>
	Payload map[string]interface{} `bson:",omitempty" json:",omitempty"`
}



// var msgChan chan nsq.Message
var (
	eventBox = make([]Event,100)
	LibraAddr string  // 192.168.19.251
	IOAddr string  // 192.168.7.11
)

// get event message from sensor nsq
func getNsqMessage(){
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:   "events",
		Channel: "aaascg",
		Address: LibraAddr + ":4150",
		MaxInFlight: 250,
	})
	for msg := range consumer.Messages() {
		var event Event
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			fmt.Printf("Data Unmarshal error:%s", err.Error())
			continue
		}
		// append event to eventbox
		eventBox = append(eventBox, event)
		fmt.Println(event.EventType)
		msg.Finish()
	}
}

// the main func to catch the crossLine event and handler it
func main() {
	// event begin
	var peopleNum int
	var timeLimit int64
	peopleNum, timeLimit, err := readConfig()
	if err != nil {
		peopleNum = 2
		timeLimit = 10
	}
	go func() {
		getNsqMessage()
	}()
	for true {
		var event Event
		if len(eventBox) > 0 {
			event = eventBox[0]
			eventBox = eventBox[1:]
		}else{
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// First Crossing Event
		// fmt.Println(event.EventType)
		if event.EventType == 721 {
			fmt.Println("======somepeople start come in or come out=======")
			HotspotId := event.HotspotId
			threePeopleEvent(HotspotId, peopleNum, timeLimit)
		}else{
			// sleep 0.1s
			// time.Sleep(100000000)
			continue
		}
	}
	
}

// read config data from config file
func readConfig() (int, int64, error) {
	f, err := os.Open("config.txt")
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, err
	}
	buf := bufio.NewReader(f)
	var peoplenum int
	var timelimit int64
	Loop:
		for {
			b, _, err := buf.ReadLine()
			if err != nil {
				if err == io.EOF {
					break Loop
				}
				return 0, 0, err
			}
			line := strings.TrimSpace(string(b))
			if line[:9] == "peoplenum" && len(line) >= 11 {
				peoplenum, err = strconv.Atoi(line[10:])
			}
			if line[:9] == "timelimit" && len(line) >= 11 {
				timelimit, err = strconv.ParseInt(line[10:], 10, 64)
			}
			if line[:9] == "libraaddr" && len(line) >= 11 {
				LibraAddr = line[10:]
			}
			if line[:6] == "ioaddr" && len(line) >= 8 {
				IOAddr = line[7:]
			}
 		}
	f.Close()
	if peoplenum == 0 {
		peoplenum = 3
	}
	if timelimit == 0 {
		timelimit = 10
	}
	return peoplenum, timelimit, nil
}

//
func threePeopleEvent(oldHotspotId string, peopleNum int, timeLimit int64) {
	// Determine whether the incident is compliant
	client := modbus.TCPClient(IOAddr + ":502")
	PeopleCount := 1
	startTime := time.Now().Unix()
	Loop:
		for true {
			endTime := time.Now().Unix()
			useTime := endTime - startTime
			if useTime > timeLimit {
				break Loop
			}
			var event Event
			if len(eventBox) >0 {
				event = eventBox[0]
				eventBox = eventBox[1:]
			}else{
				time.Sleep(100 * time.Millisecond)
				continue
			}
			EventType := event.EventType
			if EventType != 721 {
				continue
			}
			// Compliance events
			if EventType == 721 && event.HotspotId == oldHotspotId && PeopleCount < peopleNum {
				PeopleCount += 1
			}else{
				// Trigger alarm
				fmt.Println(EventType)
				fmt.Println(event.HotspotId)
				fmt.Println(oldHotspotId)
				fmt.Println(PeopleCount)
				fmt.Println("----This Operational is violation-----")
				// _, err := client.WriteSingleCoil(100, 0x0000)
				_, err := client.WriteSingleCoil(100, 0xFF00)
				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(1 * time.Second)
				// _, err = client.WriteSingleCoil(100, 0xFF00)
				_, err = client.WriteSingleCoil(100, 0x0000)
				if err != nil {
					fmt.Println(err.Error())
				}
				goto label
			}
		}
	// Overtime alarm
	if PeopleCount != peopleNum {
		fmt.Println("++++time out++++++")
		// _, err := client.WriteSingleCoil(100, 0x0000)
		_, err := client.WriteSingleCoil(100, 0xFF00)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(1 * time.Second)
		// _, err = client.WriteSingleCoil(100, 0xFF00)
		_, err = client.WriteSingleCoil(100, 0x0000)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(PeopleCount)
		fmt.Println(peopleNum)
	}else{
		fmt.Println("This active is success")
		_, err := client.WriteSingleCoil(101, 0xFF00)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(1 * time.Second)
		// _, err := client.WriteSingleCoil(100, 0xFF00)
		_, err = client.WriteSingleCoil(101, 0x0000)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	label: fmt.Println("")
}



