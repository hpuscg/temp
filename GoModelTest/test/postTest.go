/*
#Time      :  2019/2/25 下午2:47
#Author    :  chuangangshen@deepglint.com
#File      :  postTest.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"
)

/*func main() {
	// execCmd()
	// timeTest()
	// strTest()
	// stringsTest()
	// secondTest()
	// getSysInfo()
	sysTest()
}
*/
func main() {
	const (
		a = 3
		b = 4
	)
	// c := math.Sqrt(a*a + b*b)
	// c := math.Sqrt(math.Pow(3, 2) + math.Pow(4, 2))
	// fmt.Printf("%.1f", c)
	// userName()
	systemTest()
	var info = `{"cameras":[{"cameraid":"18","errornum":201},{"cameraid":"19","errornum":201},{"cameraid":"20","errornum":201},{"cameraid":"21","errornum":201},{"cameraid":"22","errornum":201},{"cameraid":"23","errornum":201},{"cameraid":"24","errornum":201},{"cameraid":"25","errornum":201},{"cameraid":"26","errornum":201},{"cameraid":"27","errornum":201},{"cameraid":"29","errornum":201},{"cameraid":"30","errornum":201}],"nodes":[{"errornum":0,"nodeid":"a1f20f445035323133000003004d00e9"},{"errornum":0,"nodeid":"a1f20f44503532313300000300990111"},{"errornum":0,"nodeid":"a1f20f44503532313300000500ee0019"}]}`
	fmt.Println(len(info))
}

func systemTest() {
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)
}

func userName() {
	a := fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(
			"chuangangshen:901010",
		)))
	fmt.Println(a)
}

func sysTest() {
	fmt.Println(runtime.GOOS)
	fmt.Println()
}

func execCmd() {
	str := "etcdctl get /config/global/server_addr"
	cmd := exec.Command("/bin/bash", "-c", str)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	serverAddr := out.String()
	if serverAddr == "\n" {
		fmt.Printf("no server addr")
	} else {
		fmt.Printf("yes\n")
		fmt.Println(reflect.TypeOf(serverAddr))
		fmt.Printf(serverAddr)
		fmt.Println("===|" + serverAddr + "|===")
	}
}

func timeTest() {
	fileName := time.Now().Format("20060102")
	fmt.Println(fileName)
}

func strTest() {
	year := "2019"
	path := year + "/234"
	fmt.Println(path)
}

func stringsTest() {
	str := "/Users/hpu_scg/gocode/src/temp/GoModelTest/test/slice.txt"
	fileNames := strings.Split(str, "/")
	fmt.Println(fileNames)
	fileName := fileNames[len(fileNames)-1]
	fmt.Println(fileName)
}

func secondTest() {
	fmt.Println(int(time.Minute / time.Second))
}

func getSysInfo() {
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)
	fmt.Println(runtime.Version())
	fmt.Printf("%+v\n", os.Environ())
}
