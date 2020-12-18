/*
#Time      :  2020/11/16 4:40 下午
#Author    :  chuangangshen@deepglint.com
#File      :  testRpc.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
)


func main() {
	// joinString()
	PathTest()
}

func PathTest() {
	fmt.Println("???", filepath.Dir(os.Args[0]))
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println("???", dir)
}

func joinString() {
	/*LatchAddr := []string{
		"123.123",
		"234.234",
	}*/
	var LatchAddr []string
	var userData string
	for _, ip := range LatchAddr {
		if userData == "" {
			userData = ip
		} else  {
			userData = strings.Join(
				[]string{
					userData,
					ip,
				}, ";")
		}
	}
	fmt.Println(userData)
}

type RpcServer struct{}

type msgJson struct {
	Msg  string   `json:"msg"`
	To   []string `json:"to"`
	Type string   `json:"type"`
}

type response struct {
	code int
	msg  string
}

func RpcServerStart(RPCPort string) {
	rpc.Register(new(RpcServer))
	rpc.HandleHTTP()

	fmt.Println("rpc port :", RPCPort)
	lis, err := net.Listen("tcp", RPCPort)

	if err != nil {
		log.Fatalln("RPC listen error: ", err)
	}

	fmt.Fprintf(os.Stdout, "%s", "start rpc connection\n")
	http.Serve(lis, nil)
}

func (this *RpcServer) Message(req msgJson, res *response) {
	fmt.Println("req:\n")
	fmt.Println(req)
	fmt.Println("end\n")

}

func RpcClient(msg msgJson, ip string) {
	var res response

	conn, err := rpc.DialHTTP("tcp", ip)
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	err = conn.Call("RpcServer.Message", msg, &res)
	if err != nil {
		log.Fatalln("cell error: ", err)
	}
}