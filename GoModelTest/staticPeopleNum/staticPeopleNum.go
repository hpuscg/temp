package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/segmentio/nsq-go"
	"gopkg.in/mgo.v2/bson"
	"github.com/deepglint/util/jsontime"
	"encoding/json"
	"time"
	"strconv"
	"github.com/deepglint/muses/eventserver/util/gkvlitehelper"
	"net/http"
)

type Event struct {
	Id        bson.ObjectId "_id" //`json:"_id,omitempty"` //auto generate, don't fill it
	StartTime *jsontime.Timestamp //`bson:"starttime,omitempty" json:"starttime,omitempty"`
	//OriginStartTime *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event start time if the event is splitted
	//OriginEndTime   *jsontime.Timestamp `bson:",omitempty" json:",omitempty"` //the origin event end time if the event is splitted
	TimeLength           int     //ms
	SliceLength          int     //`bson:"-" json:"-"` //ms
	AlarmLevel           int     //compute by backend, don't fill it
	EventType            int     //see values in EventTypeXXX consts
	EventTypeProbability float32 //[0,1]
	PlanetId             string
	SceneId              string
	SensorId             string
	HotspotId            string //可选
	UserId               string
	PeopleId             string //需要给出全局唯一的ID
	PeopleNum            int
	Path                 []int32 `bson:",omitempty" json:",omitempty"` //mm, in sequence of [x1,y1,z1,x2,y2,z2...]
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
	CutboardTimeOffset []int               `bson:",omitempty" json:",omitempty"`
	CutboardTime       *jsontime.Timestamp `bson:",omitempty" json:",omitempty"`
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

func main() {
	var ip string
	ip = readIpFromFile()
	go getNsqMessage(ip)
	go countPeopleNum()
	port := "4545"
	http.HandleFunc("/", PostCountPeopleNumber)
	http.ListenAndServe(":" + port, nil)

}

// get libra ip from config file
func readIpFromFile() string {
	ip := "192.168.12.12"
	contents, err := ioutil.ReadFile("ip.config")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		ip = strings.Replace(string(contents),"\n","",1)
	}
	return ip
}

// set a slice for get event
var eventBox = make(chan Event, 50)

// get event message from sensor nsq
func getNsqMessage(ip string) {
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "events",
		Channel:     "mark4scg",
		Address:     ip + ":4150",
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
			eventBox <- event
			// fmt.Println("line num is: ", event.HotspotId)
		}
		msg.Finish()
	}
}

var kv *gkvlitehelper.GkvliteHelper

// 删除历史数据
func deleteDailyPeopleOldData(nowTime time.Time, category string)  {
	oldTime := nowTime.AddDate(-1, 0,0)
	oldTimeStr := strconv.FormatInt(nowTime.AddDate(-1, 0, 0).Unix(), 10)
	gkvOldKey := oldTimeStr + category

	_, err := kv.ReadFromFile([]byte(gkvOldKey), &oldTime, "")
	if err == nil {
		kv.DeleteFromFile([]byte(gkvOldKey), &oldTime, "")
	}
}

// 根据天、小时、分钟统计客流量
func countPeopleNum() {
	kv = gkvlitehelper.NewGkvliteHelper("dbData", "peopleNum", gkvlitehelper.FileUnitHourly)
	for true {
		var countDay int
		var event Event
		// 获取当天零点的时间戳
		startTimeDay := time.Now().Format("2006-01-02")
		timeDay, _ := time.Parse("2006-01-02 15:04:05", startTimeDay + " 00:00:00")
		timeStampDay := timeDay.Unix()
		gkvKeyDay := strconv.FormatInt(timeStampDay, 10) + "day"
		bodyDay, err := kv.ReadFromFile([]byte(gkvKeyDay), &timeDay, "")
		if err != nil {
			countDay = 0
		} else {
			countDay, err = strconv.Atoi(string(bodyDay))
		}
		go deleteDailyPeopleOldData(timeDay, "day")
		for true {
			nowTimeDay := time.Now().Format("2006-01-02")
			var countHour int
			if startTimeDay != nowTimeDay {
				fmt.Println("The day of " + startTimeDay + " end")
				break
			}
			// 获取当前小时零分的时间戳
			startTimeHour := time.Now().Format("2006-01-02 15")
			timeHourFuture, _ := time.Parse("2006-01-02 15:04:05", startTimeHour + ":00:00")
			timeHour := timeHourFuture.Add(- 8 * 3600 * 1000000000)
			timeStampHour := timeHour.Unix()
			gkvKeyHour := strconv.FormatInt(timeStampHour, 10) + "hour"
			bodyHour, err := kv.ReadFromFile([]byte(gkvKeyHour), &timeHourFuture, "")
			if err != nil {
				countHour = 0
			} else {
				countHour, err = strconv.Atoi(string(bodyHour))
			}
			go deleteDailyPeopleOldData(timeHour, "hour")
			for true {
				nowTimeHour := time.Now().Format("2006-01-02 15")
				var countMinute int
				if startTimeHour != nowTimeHour {
					fmt.Println("The hour of " + startTimeHour + " end")
					break
				}
				// 获取当前分钟零秒的时间戳
				startTimeMinute := time.Now().Format("2006-01-02 15:04")
				timeMinuteFuture, _ := time.Parse("2006-01-02 15:04:05", startTimeMinute + ":00")
				timeMinute := timeMinuteFuture.Add(- 8 * 3600 * 1000000000)
				timeStampMinute := timeMinute.Unix()
				gkvKeyMinute := strconv.FormatInt(timeStampMinute, 10) + "minute"
				bodyMinute, err := kv.ReadFromFile([]byte(gkvKeyMinute), &timeMinute, "")
				if err != nil {
					countMinute = 0
				} else {
					countMinute, err = strconv.Atoi(string(bodyMinute))
				}
				go deleteDailyPeopleOldData(timeMinute, "minute")
				for true {
					nowTimeMinute := time.Now().Format("2006-01-02 15:04")
					if startTimeMinute == nowTimeMinute {
						if event.EventType == 721 {
							countDay++
							countHour++
							countMinute++
						}
					} else {
						fmt.Println("The minute of " + startTimeMinute + " end")
						break
					}
					// 实时更新当前分钟内的客流量
					setBodyMinute := []byte(strconv.Itoa(countMinute))
					_, err := kv.SetToFile(setBodyMinute, &timeMinuteFuture, []byte(gkvKeyMinute), "")
					if err != nil {
						fmt.Println(err)
					}
					// 实时更新当前小时段内的客流量
					setBodyHour := []byte(strconv.Itoa(countHour))
					_, err = kv.SetToFile(setBodyHour, &timeHourFuture, []byte(gkvKeyHour), "")
					if err != nil {
						fmt.Println(err)
					}
					// 实时更新当前天的客流量
					setBodyDay := []byte(strconv.Itoa(countDay))
					_, err = kv.SetToFile(setBodyDay, &timeDay, []byte(gkvKeyDay), "")
					if err != nil {
						fmt.Println(err)
					}
					for {
						select {
						case event = <- eventBox:
							break
						case <- time.After(5 * time.Second):
							event = Event{}
							break
						default:
							time.Sleep(5 * time.Second)
							continue
						}
					}
					event = <- eventBox
					time.Sleep(10000000)
				}
			}
		}
	}
}

// 根据用户需求统计人数
type timeLength struct {
	StartTime *jsontime.Timestamp
	EndTime *jsontime.Timestamp
	Category string
}




func PostCountPeopleNumber(rw http.ResponseWriter, req *http.Request)  {
	var ret []interface{}
	body, _ := ioutil.ReadAll(req.Body)
	var timeLimit timeLength
	err := json.Unmarshal(body, &timeLimit)
	if err != nil {
		info := "count people number query post body format error"
		// glog.Errorf("%v: %v", info, err)
		rw.WriteHeader(400)
		rw.Write([]byte(info))
		return
	}
	classification := timeLimit.Category
	switch classification {
	case "day":
		zeroTimeDay := timeLimit.StartTime.Format("2006-01-02") + " 00:00:00"
		zeroTimeDayStamp, _ := time.Parse("2006-01-02 15:04:05", zeroTimeDay)
		beginTime := zeroTimeDayStamp.Unix()
		afterTime := timeLimit.EndTime.Unix()
		for true {
			if beginTime + 24 * 3600 >= afterTime {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "day"
				timeStampDay, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime - 8 * 3600, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampDay, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				break
			} else {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "day"
				timeStampDay, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime - 8 * 3600, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampDay, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				beginTime = beginTime + 24 * 3600
			}
		}
	case "hour":
		zeroTimeHour := timeLimit.StartTime.Format("2006-01-02 15") + ":00:00"
		zeroTimeHourStamp, _ := time.Parse("2006-01-02 15:04:05", zeroTimeHour)
		beginTime := zeroTimeHourStamp.Unix()
		afterTime := timeLimit.EndTime.Unix()
		for true {
			if beginTime + 3600 >= afterTime {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "hour"
				timeStampHour, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampHour, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				break
			} else {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "hour"
				timeStampHour, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampHour, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				beginTime = beginTime + 3600
			}
		}
	case "minute":
		zeroTimeMinute := timeLimit.StartTime.Format("2006-01-02 15:04") + ":00"
		zeroTimeMinuteStamp, _ := time.Parse("2006-01-02 15:04:05", zeroTimeMinute)
		beginTime := zeroTimeMinuteStamp.Unix()
		afterTime := timeLimit.EndTime.Unix()
		for true {
			if beginTime + 60 >= afterTime {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "minute"
				timeStampMinute, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampMinute, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				break
			} else {
				retMember := make(map[string]string)
				gkvKey := strconv.FormatInt(beginTime, 10) + "minute"
				timeStampMinute, _ := time.Parse("2006-01-02 15:04:05", time.Unix(beginTime, 0).Format("2006-01-02 15:04:05"))
				num, err := kv.ReadFromFile([]byte(gkvKey), &timeStampMinute, "")
				realNum := string(num)
				var retValue string
				if err != nil || "" == realNum {
					retValue = "0"
				} else {
					retValue = string(num)
				}
				retKey := strconv.FormatInt(beginTime, 10)
				retMember["daytime"] = retKey
				retMember["peoplenum"] = retValue
				ret = append(ret, retMember)
				beginTime = beginTime + 60
			}
		}
	default:
		info := "the category of search format error"
		// glog.Errorf("%v", info)
		rw.WriteHeader(400)
		rw.Write([]byte(info))
		return
	}
	rw.WriteHeader(200)
	res, _ := json.Marshal(ret)
	rw.Write([]byte(res))
}

