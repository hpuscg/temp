/*
#Time      :  2019/3/28 下午5:05 
#Author    :  chuangangshen@deepglint.com
#File      :  ipTest.go
#Software  :  GoLand
*/
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
)

func main() {
	// GetIpByInterfaceAddrs()
	// GetIpByInterface()
	// GetIpByDial()
	t := time.Now()
	weekday := t.Weekday()
	fmt.Println(1<<byte(weekday))
	// ByteTest()
}

func ByteTest() {
	a := "aaaa"
	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(b[0]))
}

func GetIpByDial() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	fmt.Println(localAddr.IP)
}

func GetIpByInterfaceAddrs() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Printf("%+v\n", addrs)
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}

		}
	}
}

func GetIpByInterface() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("ifaces: %+v\n", ifaces)
	for _, iface := range ifaces {
		// fmt.Printf("iface: %+v\n", iface)
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
		}
		if iface.Name != "eth0" {
			continue
		}
		fmt.Printf("addrs: %+v\n", addrs)
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			fmt.Println(ip.String())
		}
	}
	fmt.Println("Check network connection first.")
}


