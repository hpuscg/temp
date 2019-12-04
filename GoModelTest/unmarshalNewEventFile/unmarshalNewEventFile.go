/*
#Time      :  2019/6/28 下午1:10 
#Author    :  chuangangshen@deepglint.com
#File      :  unmarshalNewEventFile.go
#Software  :  GoLand
*/
package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"bytes"
	"os/exec"
	"os"
	"encoding/hex"
	"encoding/json"
	"time"
	"github.com/deepglint/muses/eventserver/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/deepglint/muses/eventserver/util/gkvlitehelper"
)

const (
	TarPath = "/Users/hpu_scg/gocode/src/temp/GoModelTest/unmarshalNewEventFile/"
)

func main() {
	ListTarEvent(TarPath)
	GetEvent(TarPath)
}

// 获取各个设备的event，存入json文件
func GetEvent(path string) {
	starttime := 1558504381
	query_starttime := time.Unix(int64(starttime), 0)
	endtime := 1561699916
	query_endtime := time.Unix(int64(endtime), 0)
	dateFiles, _ := ioutil.ReadDir(path)
	order := 1 //old to new
	if time.Time(query_endtime).Before(time.Time(query_starttime)) {
		order = -1 //new to old
	}
	for _, dateF := range dateFiles {
		if !strings.HasPrefix(dateF.Name(), "a1") {
			continue
		}
		uid := dateF.Name()
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
		collection := path + dateF.Name() + "/event"
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
						// time.Sleep(10 * time.Second)
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

// 解压压缩的event数据，归类到各自uid文件夹
func ListTarEvent(path string) {
	tarFiles, _ := ioutil.ReadDir(path)
	for _, dateF := range tarFiles {
		fmt.Println(dateF.Name())
		if !strings.HasSuffix(dateF.Name(), "tar.gz") {
			continue
		}
		uid := strings.Split(dateF.Name(), "_")[0]
		fmt.Println(uid)
		mkdirEventFile := "mkdir " + uid
		ExecCmd(mkdirEventFile)
		tarStr := "tar -zxvf " + dateF.Name()
		ExecCmd(tarStr)
		mvStr := "mv data/tf/eventserver/event " + uid
		ExecCmd(mvStr)
		rmDataStr := "rm -rf data"
		ExecCmd(rmDataStr)
		rmTarFile := "rm " + dateF.Name()
		ExecCmd(rmTarFile)
	}
}

// 执行命令行
func ExecCmd(str string) string {
	fmt.Println(str)
	cmd := exec.Command("/bin/bash", "-c", str)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return "\n"
	}
	ret := out.String()
	return ret
}
