package main

import (
	"encoding/json"
	"fmt"
)

type Pl struct {
	Weight int
	High   int
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

func struct2map() {
	p := Pl{
		Weight: 12,
		High:   24,
	}
	data, _ := json.Marshal(p)
	fmt.Println(data, p)
	fmt.Printf("%s", string(data))
	/*tmp := make(map[string]int)
	json.Unmarshal(data, &tmp)
	fmt.Println(tmp["weight"])*/
}

func structInit() {
	sshreq := struct {
		persist bool
		open    bool
		port    int
	}{
		false,
		true,
		22,
	}
	fmt.Printf("%+v", sshreq)
}

func main() {
	// initStruck()
	// initMap()
	// structTest()
	// struct2map()
	structInit()
}
