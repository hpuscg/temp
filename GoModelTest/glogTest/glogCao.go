/*
#Time      :  2019/11/13 上午10:14 
#Author    :  chuangangshen@deepglint.com
#File      :  glogCao.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/flowservice/models"
	"encoding/json"
	"gitlab.deepglint.com/junkaicao/glog"
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	TestJson()
}

func TestJson() {
	var err error
	tmpBuf := "{\"AlarmLevel\": 0,\"DetectionStatus\": [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1],\"EventType\": 222,\"EventTypeProbability\": 1.0,\"FrameRate\": 15,\"HotspotId\": \"door_1\",\"Path\": [-147, 1301, 1058, -146, 1300, 1058, -146, 1299, 1059, -147, 1298, 1058, -149, 1298, 1058, -152, 1299, 1058, -157, 1300, 1058, -161, 1301, 1057, -165, 1302, 1057, -170, 1302, 1057, -174, 1301, 1056, -179, 1299, 1056, -181, 1298, 1056, -181, 1298, 1056, -180, 1299, 1055],\"PeopleId\": \"0607dd3333fa41a98506de67f3112b4d\",\"PeopleNum\": 3,\"PicBinary\": \"\",\"PlanetId\": \"MF\",\"SceneId\": \"D1\",\"SensorId\": \"a1f211444554333435000003001600c1\",\"StartTime\": 1573610994957,\"TimeLength\": 1000,\"UserData\": \"\",\"UserId\": \"\" }"

	glog.Infoln(string(tmpBuf))
	var event *models.Event
	err = json.Unmarshal([]byte(tmpBuf), &event)
	if err != nil {
		glog.Errorln("unmarshal event err: ", err)
		return
	}
	glog.Infof("%+v", event)
	fmt.Printf("%+v", event)
}
