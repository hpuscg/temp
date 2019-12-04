/*
#Time      :  2019/6/25 上午11:32 
#Author    :  chuangangshen@deepglint.com
#File      :  unmarshalEventFile.go
#Software  :  GoLand
*/
package main

import (
	"os"
	"fmt"
	"github.com/steveyen/gkvlite"
	"strings"
	"bufio"
	"io"
	"errors"
	"strconv"
	"github.com/deepglint/muses/eventserver/util/gkvlitehelper"
	"time"
	"io/ioutil"
	"github.com/deepglint/muses/eventserver/models"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

const (
	Path = "/Users/hpu_scg/Desktop/event/"
)

func main() {
	GetIpUidMap()
	GetEvent(Path)

}

func GetEvent(path string) () {
	starttime := 1558504381
	query_starttime := time.Unix(int64(starttime), 0)
	endtime := 1561614781
	query_endtime := time.Unix(int64(endtime), 0)
	dateFiles, _ := ioutil.ReadDir(path)
	order := 1 //old to new
	if time.Time(query_endtime).Before(time.Time(query_starttime)) {
		order = -1 //new to old
	}
	for _, dateF := range dateFiles {
		if !strings.HasSuffix(dateF.Name(), "eventserver") {
			continue
		}
		ipNames := strings.Split(dateF.Name(), "_")
		fmt.Println(ipNames[0])
		uid := uidIpMap[ipNames[0]]
		fmt.Println(uid)
		if "" == uid {
			fmt.Println("===32==")
			continue
		}
		jsonFileName := path + uid + ".json"
		os.Create(jsonFileName)
		file6, err := os.OpenFile(jsonFileName, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
		}
		collection := path + dateF.Name()
		eventKv := gkvlitehelper.NewGkvliteHelper(collection, "events", gkvlitehelper.FileUnitHourly)
		tssKv := gkvlitehelper.NewGkvliteHelper(collection, "eventtss", gkvlitehelper.FileUnitDaily)
		tssOlap, err := models.LoadTSSOlapFile(tssKv, 170)
		if err != nil {
			fmt.Println(err)
		}
		var (
			keyset models.KeySet
			idmap = make(map[string]int)
		)
		keyset.SensorID = uid
		keyset.EventType = 0
		tss, ok := tssOlap.Data[keyset.String()]
		fmt.Println("===75=", ok, keyset.String())
		if ok {
			// timestamp, _ := time.Parse("2006-01-02 15:04:05", time.Unix(starttime, 0).Format("2006-01-02 03:04:05"))
			endtime := time.Time(query_endtime)
			fmt.Println(endtime)
			for {
				ids, search_endtime, partial := tss.FindResultPaging(endtime, 100, -1, true)
				fmt.Println("===80", partial, search_endtime)
				if 0 == len(ids) {
					fmt.Println("")
					break
				}
				for _, id := range ids {
					idstr := id.(map[string]interface{})["id"].(string)
					if len(idstr) == 0 {
						continue
					}
					_, ok := idmap[idstr]
					// glog.Infof("v: %v, ok: %v", v, ok)
					if !ok {
						id, err := hex.DecodeString(idstr)
						if err != nil {
							fmt.Println("invalid event id")
						}
						//glog.Infof("id hex: %v", id)
						//
						objectid := bson.ObjectIdHex(idstr)
						//glog.Infof("object id: %v", objectid)
						timestamp := objectid.Time()
						//glog.Infof("timestamp from id: %v", timestamp)
						body, err := eventKv.ReadFromFile(id, &timestamp, "")
						if err != nil {
							fmt.Println(err)
						}
						var event models.Event
						err = json.Unmarshal(body, &event)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Printf("==118%+v\n", event)
						if 0 == event.EventType {
							continue
						}
						interfaceevent := models.EventToInterface(&event, "json", "", "")
						byteevent, err := json.Marshal(interfaceevent)
						file6.WriteString(string(byteevent))
						file6.WriteString("\n")
						thisstarttime := time.Time(*event.StartTime)
						thisendtime := thisstarttime.Add(time.Millisecond * time.Duration(event.TimeLength))
						if (order == 1 && thisstarttime.After(time.Time(query_endtime))) ||
							(order == -1 && thisendtime.Before(time.Time(query_endtime))) {
							break
						}
					}
				}
				endtime = search_endtime
			}

		}
		file6.Close()

	}
	return
}

var uidIpMap = make(map[string]string)

func GetIpUidMap() {
	f, err := os.Open("./uidIp.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println(string(line))
		names := strings.Split(string(line), "=")
		uidIpMap[names[0]] = names[1]
	}
	fmt.Println(uidIpMap)
}

/*func ReadGkvFile(path, jsonFile string) {
	f, _ := os.Open(path)
	defer f.Close()
	s, err := gkvlite.NewStore(f)
	if err != nil {
		fmt.Println("28===", err)
	}
	c := s.GetCollection("events")
	//fmt.Println(c)
	data := c.GetValues(true)
	if err != nil {
		fmt.Println("54===", err)
	}
	file6, err := os.OpenFile(jsonFile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range data {
		fmt.Println(string(item.Val))
		file6.WriteString(string(item.Val))
		file6.WriteString("\n")
	}
	file6.Close()
}*/

func DeleteNotEventValue(path, gkvFile, datapath string) (err error) {
	f, _ := os.Open(path)
	defer f.Close()
	s, err := gkvlite.NewStore(f)
	if err != nil {
		fmt.Println("74===", err)
	}
	c := s.GetCollection("events")
	data, err := c.MaxItem(true)
	if err != nil {
		fmt.Println("79===", err)
	}
	if strings.HasSuffix(string(data.Key), "hour") || strings.HasSuffix(string(data.Key), "minute") ||
		strings.HasSuffix(string(data.Key), "day") {
		timestamp := TimeParse(string(data.Key))
		fmt.Println(timestamp)
		kv := gkvlitehelper.NewGkvliteHelper(datapath, "events", gkvlitehelper.FileUnitHourly)
		ret, err := kv.DeleteFromFile(data.Key, &timestamp, "")
		_, err = c.Delete(data.Key)
		if err != nil {
			fmt.Println("85===", err)
		}
		fmt.Println(string(data.Key))
		fmt.Println(ret, err)
		time.Sleep(5 * time.Second)
	} else {
		return errors.New("over")
	}
	s.Flush()
	return
}

func TimeParse(timeyear string) (timestamp time.Time) {
	if strings.HasSuffix(timeyear, "day") {
		timestr := strings.Split(timeyear, "day")
		timestampint, _ := strconv.Atoi(timestr[0])
		timestamp, _ = time.Parse("2006-01-02 15:04:05", time.Unix(int64(timestampint)- 8*3600, 0).Format("2006-01-02 03:04:05"))
	} else if strings.HasSuffix(timeyear, "hour") {
		timestr := strings.Split(timeyear, "hour")
		timestampint, _ := strconv.Atoi(timestr[0])
		timestamp, _ = time.Parse("2006-01-02 15:04:05", time.Unix(int64(timestampint)- 8*3600, 0).Format("2006-01-02 03:04:05 PM"))
	} else if strings.HasSuffix(timeyear, "minute") {
		timestr := strings.Split(timeyear, "minute")
		timestampint, _ := strconv.Atoi(timestr[0])
		timestamp, _ = time.Parse("2006-01-02 15:04:05", time.Unix(int64(timestampint)- 8*3600, 0).Format("2006-01-02 03:04:05 PM"))
	}
	return
}

