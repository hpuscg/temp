/*
#Time      :  2018/11/30 下午4:05 
#Author    :  chuangangshen@deepglint.com
#File      :  saveEventToDB.go
#Software  :  GoLand
*/
package main

import (
	"github.com/go-pg/pg"
	"github.com/deepglint/util/jsontime"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-pg/pg/orm"
	"fmt"
)

const (
	addr     = "192.168.100.235:5432"
	user     = "postgres"
	passWard = "deepglint"
	dbName   = "libraT"
)

func main() {
	db := connect()
	err := CreateTable(db)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db)
	}
}

func connect() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     user,
		Addr:     addr,
		Password: passWard,
		Database: dbName,
	})
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		panic(err)
	}
	return db
}

type Events struct {
	Id        int64
	TimeStamp int64
	SensorId  string
	EventType int      `sql:"type:smallint"`
	EventInfo []byte   `sql:"type:text[]"`
	tableName struct{} `sql:"events"`
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
	PeopleId             string                                                       //需要给出全局唯一的ID
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

func CreateTable(db *pg.DB) error {
	for _, model := range []interface{}{&Events{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
