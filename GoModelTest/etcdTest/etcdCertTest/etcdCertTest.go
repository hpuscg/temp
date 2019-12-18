/*
#Time      :  2019/12/18 3:52 PM 
#Author    :  chuangangshen@deepglint.com
#File      :  etcdCertTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/coreos/go-etcd/etcd"
	"fmt"
)

var (
	server = "192.168.100.235:4001"
	cert   = "/Users/hpu_scg/gocode/src/temp/GoModelTest/etcdTest/etcdCertTest/client.crt"
	key    = "/Users/hpu_scg/gocode/src/temp/GoModelTest/etcdTest/etcdCertTest/client.key.insecure"
	caCert = "/Users/hpu_scg/gocode/src/temp/GoModelTest/etcdTest/etcdCertTest/ca.crt"
)

func main() {
	// etcdTest()
	etcdCertTest()
}

func etcdTest() {
	cli := etcd.NewClient([]string{"http://" + server})
	value, err := cli.Get("/config/image/armhf-etcd/cur_version", false, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", value)
}

func etcdCertTest() {
	cli, err := etcd.NewTLSClient([]string{"https://" + server}, cert, key, caCert)
	if err != nil {
		fmt.Println(err)
	}
	value, err := cli.Get("/config/image", false, false)
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	for _, value1 := range value.Node.Nodes {
		count++
		fmt.Printf("%+v\n", value1.Key)
		value2, _ := cli.Get(value1.Key + "/cur_version", false, false)
		fmt.Println(value2.Node.Value)
		value3, _ := cli.Get(value1.Key + "/" + value2.Node.Value + "/image", false, false)
		fmt.Println(value3.Node.Value)
		value4, _ := cli.Get(value1.Key + "/" + value2.Node.Value + "/script", false, false)
		fmt.Println(value4.Node.Value)
		fmt.Println(count)
		fmt.Println("==============")
	}
}
