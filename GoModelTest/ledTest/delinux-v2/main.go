package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

const (
	LedUrl = "http://127.0.0.1:9000/sys/devio"
)

func main() {
	// 打开led
	openLed()
	fmt.Println("打开led成功")
	// 打开红外补光灯
	openInfrared()
	fmt.Println("打开红外补光灯成功")
	// 毫米波传感器上电自动开启
	fmt.Println("毫米波传感器上电自动开启")
	// 打开激光发射器
	fmt.Println("请单独打开激光发射器")
}

func openInfrared() {
	if _, err := CommandOut("./io", "-4", "-w", "0xFF77e028", "0xFFFF115F"); err != nil {
		panic(err)
	}
	if _, err := CommandOut("./io", "-4", "-w", "0xFF42001c", "3"); err != nil {
		panic(err)
	}
	if _, err := CommandOut("./io", "-4", "-w", "0xFF420014", "48000"); err != nil {
		panic(err)
	}
	if _, err := CommandOut("./io", "-4", "-w", "0xFF420018", "24000"); err != nil {
		panic(err)
	}
}

func CommandOut(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func openLed() {
	if _, err := CommandOut("set_devio", "-gpio", "green", "on"); err != nil {
		panic(err)
	}
	if _, err := CommandOut("set_devio", "-gpio", "yellow", "on"); err != nil {
		panic(err)
	}
	if _, err := CommandOut("set_devio", "-gpio", "red", "on"); err != nil {
		panic(err)
	}
}
