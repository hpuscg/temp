/*
#Time      :  2019/10/29 下午5:02 
#Author    :  chuangangshen@deepglint.com
#File      :  event.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/muses/util/jsontime"
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	Id                   bson.ObjectId "_id" //`json:"_id,omitempty"` //auto generate, don't fill it
	StartTime            *jsontime.Timestamp //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	TimeLength           int                 //ms
	SliceLength          int                 //`bson:"-" json:"-"` //ms
	AlarmLevel           int                 //compute by backend, don't fill it
	EventType            int                 //see values in EventTypeXXX consts
	EventTypeProbability float32             //[0,1]
	PlanetId             string
	SceneId              string
	SensorId             string
	HotspotId            string //可选，不需要libra提供
	UserId               string
	PeopleId             string                                              //需要给出全局唯一的ID
	Path                 []int32       `bson:",omitempty" json:",omitempty"` //mm, in sequence of [x1,y1,z1,x2,y2,z2...]
	DetectionStatus      []int         `bson:",omitempty" json:",omitempty"` // 1 is real, and 0 is ghost.
	FrameRate            int                                                 // FrameRate indicates the FPS of the trajectory points.
	UserData             string        `json:",omitempty"`                   //个性化数据
}
