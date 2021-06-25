/*
#Time      :  2021/3/25 7:46 下午
#Author    :  chuangangshen@deepglint.com
#File      :  wsClient.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type AlarmEvent struct {
	DeviceId    string   `json:"device_id"`
	EventType   int      `json:"event_type"`
	EventTime   int      `json:"event_time"`
	DeviceType  int      `json:"device_type"` // 1-T，2-F，3-L
	VideoPath   string   `json:"video_path"`
	VideoStatus int      `json:"video_status"` // 0-合成中，1-完成，2-异常
	EventId     string   `json:"event_id"`
	EventLevel  int      `json:"event_level"` // 1-紧急，2-严重，3-一般，4-设备
	PicturePath []string `json:"picture_path" gorm:"type:VARCHAR(10000)"`
}

var addr = flag.String("addr", "192.168.100.238:8789", "http service address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// go timeWriter(conn)

	fmt.Println("begin......")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		var alarmEvent AlarmEvent
		json.Unmarshal(message, &alarmEvent)
		fmt.Printf("received: %+v\n", alarmEvent)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
