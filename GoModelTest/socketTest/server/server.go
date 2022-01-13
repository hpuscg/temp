package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handlerConnection(conn)
	}
}

func handlerConnection(conn net.Conn) {
	defer conn.Close()
	data, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C
		_, err := conn.Write([]byte("ok\n"))
		if err != nil {
			log.Println(err)
			break
		}
	}
}
