package main

import "fmt"

func main() {
	// testInterface()
	intTest()
}

func testInterface() {
	var value interface{}
	fmt.Println(value)
}

func intTest() {
	fmt.Println(5 / 2)
}
