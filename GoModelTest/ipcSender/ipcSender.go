/*
#Time      :  2019/3/1 下午5:34 
#Author    :  chuangangshen@deepglint.com
#File      :  ipcSender.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/flowservice/glog"
	"github.com/op/go-nanomsg"
	"time"
)

func main() {
	var ipc *IpcSource
	ipc, err := NewIPCSource("ipc:///tmp/libra_cutboard.ipc")
	if err != nil {
		glog.Errorln("Fail to init IPC")
	}
	data := []byte("")
	for {
		SendEvent(ipc, data)
		time.Sleep(2 * time.Second)
	}
}

type IpcSource struct {
	addr string
	socket *nanomsg.PushSocket
	exitchan chan bool
}

func NewIPCSource(addr string) (*IpcSource, error) {
	var (
		ipc = new(IpcSource)
		err error
	)
	ipc.addr = addr
	ipc.socket, err = nanomsg.NewPushSocket()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	ipc.socket.Bind(addr)
	glog.Infoln("Finishing initializing ipc source:%s", addr)
	ipc.exitchan = make(chan bool, 1)
	return ipc, nil
}

func SendEvent(ipc *IpcSource, data []byte) error {
	_, err := ipc.socket.Send(data, 0)
	if err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}



