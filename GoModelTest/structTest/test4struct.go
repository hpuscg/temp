package main

import (
	"fmt"
	"encoding/json"
)

type Pl struct {
	weight int
	high int
}

func structTest() {
	// pl1 := Pl{}
	fmt.Println("pl1.age")
}

/*
func initStruck()  {
	var p Person
	p.age = 12
	p.name = "ltt"
	fmt.Println(p)
}
*/

func initMap() {
	m := make(map[string]string)
	m["age"] = "12"
	m["name"] = "ltt"
	fmt.Println(m)
	for k, v := range m {
		fmt.Println("k is :", k)
		fmt.Println("v is :", v)
	}
}

func struct2map()  {
	p := Pl{
		weight:12,
		high:24,
	}
	data, _:= json.Marshal(p)
	tmp := make(map[string]int)
	json.Unmarshal(data, &tmp)
	fmt.Println(tmp["weight"])
}

func main() {
	// initStruck()
	// initMap()
	// structTest()
	struct2map()
}
