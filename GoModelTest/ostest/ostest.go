package main

import (
	"io/ioutil"
	"fmt"
)

func main() {
	/*
	serverAddr := "192.168.4.42"
	conn, err := net.DialTimeout("tcp", serverAddr+ ":8008", 2*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if conn == nil {
		fmt.Println("can't conntect")
		return
	}else{
		resp, err := http.Get("http://" + serverAddr + ":8008/api/version")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		code := resp.StatusCode
		if code == 200 {
			fmt.Println("this is a librat")
		}else{
			fmt.Println("this is a server")
		}
	}
	*/

	// getDirList(`/`)

	// ifIn()
}


/*
func main() {
	a := 1
	b := "ij"
	fmt.Println(reflect.TypeOf(a))
	fmt.Println(reflect.TypeOf(b))
}
*/

func getDirList(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	fmt.Println(path)
}

func ifIn()  {
	sliceA := []string {"a", "b", "c"}
	sliceB := []string {"c", "d", "a"}
	var sliceC []string
	for _, a := range sliceA {
		for _, b := range sliceB {
			if a == b {
				sliceC = append(sliceC, a)
			}
		}
	}
	fmt.Println(sliceC)
}
