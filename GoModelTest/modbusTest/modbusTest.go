package main

import (
	"flag"
	"temp/GoModelTest/modbusTest/modbusHelper"
	"time"

	"github.com/goburrow/modbus"
	"gitlab.deepglint.com/junkaicao/glog"
)

var (
	address  = "/dev/ttyS0"
	waitTime = 1
	selfDemo = true
)

func main() {
	flag.IntVar(&waitTime, "time", 1, "time out time")
	flag.BoolVar(&selfDemo, "demo", true, "is start self demo")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./"), glog.WithLevel("info"))
	if selfDemo {
		selfModbusDemo()
	} else {
		goBurrowModbus()
	}
}

func selfModbusDemo() {
	glog.Infoln(waitTime)
	handler := modbusHelper.NewRTUClientHandler(address)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = time.Duration(waitTime) * time.Second
	handler.FileName = "/sys/class/gpio/gpio5/value"
	handler.InValue = "0"
	handler.OutValue = "1"
	handler.FileEnable = true

	glog.Infof("%+v", handler)

	err := handler.Connect()
	if err != nil {
		glog.Infoln(err)
		return
	}
	defer handler.Close()
	client := modbusHelper.NewClient(handler)
	for {

		results, err := client.ReadHoldingRegisters(0, 10)
		if err != nil {
			glog.Infoln(err)
			// continue
		}
		for key, value := range results {
			if key == 1 {
				glog.Warningln(value)
			}
		}
		glog.Infoln(results)
	}
	/*results, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		glog.Infoln(err)
		// continue
	}
	time.Sleep(1 * time.Second)
	glog.Infoln(results)*/

}

func goBurrowModbus() {
	handler := modbus.NewRTUClientHandler(address)

	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	glog.Infoln(err)
	defer handler.Close()

	client := modbus.NewClient(handler)
	for {
		results, err := client.ReadHoldingRegisters(0, 10)
		if err != nil {
			glog.Info(err)
		}
		glog.Infoln(results)
	}
}
