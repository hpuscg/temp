package main

import (
	"log"
	"net"
)

func main() {
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
