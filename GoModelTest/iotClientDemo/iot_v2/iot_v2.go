package main

import (
	"gitlab.deepglint.com/idg/FGWare/FGWeb/common/jwt"
)

/* func getDeviceInfo(i iot_client.DeviceInfo) iot_client.DeviceInfo {
	i.ID = i.ID + "-T"
	i.ExtDevInfo = map[string]interface{}{"key": "blablabla"}
	return i
}

///////////////////////////////////////////////////////////////////

func startRecord(c *iot_client.IoTClient, session string, param map[string]interface{}) {
	fmt.Println("start record")
	c.PublishCmdResponse(session, map[string]interface{}{"status": "ok"})
}

func stopRecord(c *iot_client.IoTClient, session string, param map[string]interface{}) {
	fmt.Println("stop record")
	c.PublishCmdResponse(session, map[string]interface{}{"status": "ok"})
}

///////////////////////////////////////////////////////////////////

var configHash string = "0"

func onConfigChange(c *iot_client.IoTClient, config map[string]interface{}, hash string) bool {
	fmt.Printf("receive config %+#v", config)
	configHash = hash
	return true
}

func onGetConfigHash(*iot_client.IoTClient) string {
	return configHash
} */

func main() {
	/* client := iot_client.NewIotClinet()
	client.OnConfigChange = onConfigChange
	client.OnGetConfigHash = onGetConfigHash
	// 扩展支持ota
	deExt := iot_delinux.Extender{}
	deExt.DevInfoHook = getDeviceInfo
	client.Extend(deExt)
	// 扩展userCmd
	client.RegisterUserCmdCallback(10001, startRecord)
	client.RegisterUserCmdCallback(10001, stopRecord)

	client.SetIoTServer(iot_client.IoTConfig{
		Server:   "172.17.11.4:1883",
		Topic:    "topic/test",
		Username: "haomul",
		Password: "haomul!@#$",
	})

	client.Start()

	for {
		time.Sleep(time.Minute)
	} */
	jwt.GenToken("")

}
