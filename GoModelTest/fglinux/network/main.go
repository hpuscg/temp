package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
)

var (
	dev, ip, netmask, gateway string
)

func main() {
	flag.StringVar(&dev, "dev", "eth0", "dev info, default: eth0")
	flag.StringVar(&ip, "ip", "x.x.x.x", "dev info, default: x.x.x.x")
	flag.StringVar(&netmask, "netmask", "x.x.x.x", "dev info, default: x.x.x.x")
	flag.StringVar(&gateway, "gateway", "x.x.x.x", "dev info, default: x.x.x.x")
	flag.Parse()
	if ip == "" || netmask == "" || gateway == "" {
		fmt.Printf("ip: %s\n", ip)
		fmt.Printf("netmask: %s\n", netmask)
		fmt.Printf("gateway: %s\n", gateway)
		fmt.Println("please input right ip netmask gateway")
		return
	}
	cmd := exec.Command("set_net", "-dev", dev, "-ip", ip, "-netmask", netmask, "-gateway", gateway)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(out.String())
	fmt.Println(err)
}
