package main

import (
	"github.com/goburrow/modbus"
	"fmt"
)

func main() {
	client := modbus.TCPClient("192.168.7.11:502")
	// Read input register 9
	// _, err := client.WriteSingleCoil(100, 0xFF00)
	_, err := client.WriteSingleCoil(100, 0x0000)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("--------")
	}
	results, err := client.ReadCoils(100, 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(results)
}

