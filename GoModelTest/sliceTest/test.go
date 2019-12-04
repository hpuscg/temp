package main

import (
	"fmt"
	"encoding/json"
)

func main() {
	// sliceTest(t)
	jsonStringTest()
}


var t = [2]int{3,4}

func sliceTest(data interface{})  {
	switch vv := data.(type) {
	case interface{}:
		fmt.Println(vv)
		t2 := data.([2]int)
		fmt.Println(t2)
	}
}

type JsonStrong struct {
	Name string `json:"name"`
	Age int    `json:"age"`
}


func jsonStringTest()  {
	str1 := `{"age": 12, "name": "li"}`
	jr := JsonStrong{}
	err2 := json.Unmarshal([]byte(str1), &jr)
	fmt.Println("err2 is", err2)
	fmt.Println(jr)
}





