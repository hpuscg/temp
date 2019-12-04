package main

import (
	"fmt"
)

func main() {
	mapToStruct()
}

type StructDemo struct {
	name string
	age int
}

func mapToStruct() {
	v := make(map[string]interface{})
	v["name"] = "lisi"
	v["age"] = 13
	structTest := StructDemo{}
	for key, value := range v {
		fmt.Println(key)
		switch key {
		case "name":
			structTest.name = value.(string)
		case "age":
			structTest.age = value.(int)
		}

	}
	fmt.Println(structTest)

	/*vl, err := json.Marshal(v)
	fmt.Println("err1 is: ",err)
	err = json.Unmarshal(vl, &structTest)
	fmt.Println("err2 is: ", err)
	fmt.Println(structTest)*/
}





