package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
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
	IotTopic    = "topic/haomu"
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
	MsgType  int    `json:"msg_type"`
	DeviceId string `json:"deviceid"`
	// 请参考https://confluence.deepglint.com/pages/viewpage.action?pageId=11078244
}

func main() {
	SetOpts()
	IotConnect()
	PublishMessage()
}

// 给iot-server发送消息
func PublishMessage() {
	iotMessage := IoTMessageInfo{
		MsgType:  1,
		DeviceId: "001test",
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
	for {
		select {
		case <-ticker.C:
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
				if token := MqttClient.Subscribe(IotTopic, 0, msgRecv); token.Wait() && token.Error() != nil {
					fmt.Println("Iot client subscribe error :", token.Error())
				}
				break
			}
		}
	}
}

// 设置iot-client参数
func SetOpts() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("192.168.100.238:1883")
	opts.SetUsername("styun")
	opts.SetPassword("styun!@#$")
	opts.SetClientID("test01")
	opts.KeepAlive = 20
	opts.AutoReconnect = false
	opts.SetCleanSession(true)
	will := WillJson{
		MsgType:  255,
		DeviceId: "test01",
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
	return
}

func PublishResponseInfo(data []byte) {
	if len(data) == 0 {
		return
	}
	token := MqttClient.Publish(IotTopic, 1, false, data)
	token.Wait()
}
