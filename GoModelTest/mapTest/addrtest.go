package main

import (
	"fmt"
)

func main() {
	// addrtest()
	MapTest()
}

type addr struct {
	ip string
}

var ip1 *addr

func addrtest() {
	var ip2 *addr
	ip2 = ip1
	fmt.Println("ip1 is:", &ip1)
	fmt.Println("ip2 is:", &ip2)
}

func MapTest()  {
	type value struct {
		name  string
	}
	var data = make(map[int]*value)
	ret, ok := data[23]
	fmt.Println(ret, ok)
}


