package main

import (
	"log"
	"net"
)

func main() {
	/* t := sliceTest()
	fmt.Print(t) */
	socketClient()
}

/* func sliceTest() [][]int {
	var temp = [][]int{
		{1, 2, 3, 4},
		{3, 4, 5},
	}
	temp = append(temp, []int{1, 2, 3})
	return temp
} */

func socketClient() {
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go conn.Write([]byte("ok\n"))
	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(string(buf))
	}
}
