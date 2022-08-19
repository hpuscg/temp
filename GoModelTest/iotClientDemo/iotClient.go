package main

import (
	"encoding/json"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	MqttClient            MQTT.Client
	DefaultPublishHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("%+v\n", msg)
	}
	OnConnectHandler = func(client MQTT.Client) {
		fmt.Println("iot is on connect")
	}
	ConnectionLostHandler = func(client MQTT.Client, err error) {
		fmt.Println("mqtt lost connected")
	}
	DeviceId    = "a1f20f445035323133000003004d00e1"
	IotTopic    = "topic/haomut"
	Topic       = "msg/" + DeviceId
	MqttOptions *MQTT.ClientOptions
)

type WillJson struct {
	MsgType  int    `json:"msg_type"`
	DeviceId string `json:"deviceid"`
}

type Dispatcher struct {
	MsgType int `json:"msg_type"`
}

type IoTMessageInfo struct {
	MsgType    int                    `json:"msg_type"`
	DeviceId   string                 `json:"deviceid"`
	SN         string                 `json:"sn"`
	Name       string                 `json:"devicename"`
	Model      string                 `json:"devmodel"`
	AppName    string                 `json:"appname"`
	AppVersion string                 `json:"appver"`
	OsName     string                 `json:"osname"`
	OsVersion  string                 `json:"osver"`
	IP         string                 `json:"deviceip"`
	ProtoVer   int                    `json:"protover"`
	AppsInfo   []struct{}             `json:"appsinfo"`
	ExtDevInfo map[string]interface{} `json:"extdevinfo"`
	// 请参考https://confluence.deepglint.com/pages/viewpage.action?pageId=11078244
}

func main() {
	SetOpts()
	fmt.Printf("%+v", MqttOptions)
	IotConnect()
	fmt.Println(MqttClient)
	InitPublishMessage()
	for {
		PublishMessage()
		time.Sleep(100 * time.Millisecond)
	}
}

// 给iot-server发送消息
func PublishMessage() {
	iotMessage := IoTMessageInfo{
		MsgType:  1,
		DeviceId: DeviceId,
		Name:     time.Now().String(),
	}
	data, err := json.Marshal(iotMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	PublishResponseInfo(data)
}

// 注册设备
func InitPublishMessage() {
	iotMessage := IoTMessageInfo{
		MsgType:    1,
		DeviceId:   DeviceId,
		Model:      "HC-BA311",
		AppName:    "libraT",
		AppVersion: "11",
	}
	data, err := json.Marshal(iotMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	PublishResponseInfo(data)
}

func IotConnect() {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for range ticker.C {
		if MqttClient != nil {
			MqttClient.Disconnect(250)
			fmt.Println("IOT MqttClient disconnect!")
		}
		// 建立连接
		MqttClient = MQTT.NewClient(MqttOptions)
		if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println("IOT client connect to server fail :", token.Error())
		} else {
			fmt.Println("IOT client connect to server ok")

			// 接收iot-server的消息
			if token := MqttClient.Subscribe(Topic, 0, msgRecv); token.Wait() && token.Error() != nil {
				fmt.Println("Iot client subscribe error :", token.Error())
			}
			return
		}
	}
}

// 设置iot-client参数
func SetOpts() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.100.238:1883")
	opts.SetUsername("haomut")
	opts.SetPassword("haomut!@#$")
	opts.SetClientID(DeviceId)
	opts.KeepAlive = 20
	opts.AutoReconnect = false
	opts.SetCleanSession(true)
	// opts.SetProtocolVersion(3)
	will := WillJson{
		MsgType:  255,
		DeviceId: DeviceId,
	}
	willData, err := json.Marshal(will)
	if err != nil {
		fmt.Println(err)
	}
	opts.SetWill(IotTopic, string(willData), 1, false)
	opts.SetDefaultPublishHandler(DefaultPublishHandler)
	opts.SetOnConnectHandler(OnConnectHandler)
	opts.SetConnectionLostHandler(ConnectionLostHandler)
	MqttOptions = opts
}

func msgRecv(client MQTT.Client, message MQTT.Message) {
	dispatcher := new(Dispatcher)
	err := json.Unmarshal(message.Payload(), dispatcher)
	// TODO message.Payload()
	if err != nil {
		fmt.Printf("Iot message unmarshal fail on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
		return
	}
	fmt.Printf("Iot Received message on topic: %s msg_type: %d\nMessage: %s", message.Topic(), dispatcher.MsgType, message.Payload())
}

func PublishResponseInfo(data []byte) {
	if len(data) == 0 {
		return
	}
	fmt.Println(string(data))
	token := MqttClient.Publish(IotTopic, 1, false, data)
	token.Wait()
}
