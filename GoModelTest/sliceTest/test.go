package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

func main() {
	// sliceTest(t)
	// jsonStringTest()
	compareSlice()
}

func compareSlice() {
	a := []string{}
	b := []string{"b", "c", "a"}
	fmt.Println(reflect.DeepEqual(a, b))
	sort.Strings(a)
	sort.Strings(b)
	fmt.Println(reflect.DeepEqual(a, b))
	fmt.Println(compareByFor(a, b))
}

func compareByFor(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for key, value := range a {
		if value != b[key] {
			return false
		}
	}

	return true
}

var t = [2]int{3, 4}

func sliceTest(data interface{}) {
	switch vv := data.(type) {
	case interface{}:
		fmt.Println(vv)
		t2 := data.([2]int)
		fmt.Println(t2)
	}
}

type JsonStrong struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func jsonStringTest() {
	str1 := `{"age": 12, "name": "li"}`
	jr := JsonStrong{}
	err2 := json.Unmarshal([]byte(str1), &jr)
	fmt.Println("err2 is", err2)
	fmt.Println(jr)
}
