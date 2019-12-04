/*
#Time      :  2019/3/26 下午5:02 
#Author    :  chuangangshen@deepglint.com
#File      :  etcdWatch.go
#Software  :  GoLand
*/
package main

import (
	"strings"
	"github.com/coreos/go-etcd/etcd"
	"time"
	"fmt"
)

func main() {
	monitorEtcd()
	// fmt.Printf("%s", time.Now().Format("2006-01-02 15:04:05"))
}

func monitorEtcd() {
	path := "/config/iot"
	machines := strings.Split("http://127.0.0.1:4001", ",")
	c := etcd.NewClient(machines)
	var waitIndex uint64 = 0
	for {
		resp, err := c.Watch(path, waitIndex, true, nil, nil)
		fmt.Printf("%s===26==ing,", time.Now().Format("2006-01-02 15:04:05"))
		if resp == nil || err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		waitIndex = resp.Node.ModifiedIndex + 1

		k := resp.Node.Key
		v := resp.Node.Value
		requests := strings.Split(k, "/")
		request := requests[len(requests) -1]
		fmt.Printf("key=%s, value=%s\n", request, v)
	}
}

func oo()  {
	/*for {
		resp, err := global.BumbleConfig.EtcdClient.Watch(models.IOT_BASE_PATH, waitIndex, true, nil, nil)
		if resp == nil || err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		waitIndex = resp.Node.ModifiedIndex + 1
		keys := strings.SplitAfterN(resp.Node.Key, "/", -1)
		key := keys[len(keys)-1]
		value := resp.Node.Value
		switch key {
		case "server":
			models.IotServer = value
		case "username":
			models.IotUserName = value
		case "password":
			models.IotPassWord = value
		case "topic":
			models.IotTopic = value
		case "devmodel":
			models.IotDevModel = value
		default:
			continue
		}
		if "devmodel" == key {
			PublishDeviceInfo()
			continue
		}
		models.MqttClient.Disconnect(250)
		time.Sleep(5 * time.Second)
		InitIotServer()
	}*/
}