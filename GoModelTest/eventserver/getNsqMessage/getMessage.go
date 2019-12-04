package getNsqMessage

import (
	"github.com/segmentio/nsq-go"
	"encoding/json"
	"github.com/deepglint/util/jsontime"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	//
	Id bson.ObjectId "_id" //`json:"_id,omitempty"` //auto generate, don't fill it

	//
	StartTime *jsontime.Timestamp //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	//OriginStartTime *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event start time if the event is splitted
	//OriginEndTime   *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event end time if the event is splitted
	//
	TimeLength  int //ms
	SliceLength int //`bson:"-" json:"-"` //ms
	//
	AlarmLevel           int     //compute by backend, don't fill it
	EventType            int     //see values in EventTypeXXX consts
	EventTypeProbability float32 //[0,1]

	//
	PlanetId  string
	SceneId   string
	SensorId  string
	HotspotId string //可选

	//
	UserId string

	//
	PeopleId string //需要给出全局唯一的ID
	PeopleNum int

	//
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
	//
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


func main()  {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic: "events",
		Channel: "eventscg",
		Address: "192.168.19.251:4150",
		MaxInFlight: 250,
	})
	for msg := range consumer.Messages() {
		var event Event
		err := json.Unmarshal(msg.Body[:], &event)
		if err != nil {
			continue
		}
		fmt.Println("eventtype is :", event.EventType)
		msg.Finish()
	}
}
