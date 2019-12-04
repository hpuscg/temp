package main

import (
	"fmt"
	"github.com/deepglint/eventserver2/util/io"
	"reflect"
)

func main() {
	// setKeyNil()
	ioTest()
}

func setKeyNil()  {
	ret := make(map[string]string)
	ret[""] = "hello"
	fmt.Println(ret[""])
}

type t struct {
	value string
}

func ioTest()  {
	var t1 t
	t1.value = "123"
	io.PrintStruct(reflect.TypeOf(t1), reflect.ValueOf(t1), 2)
}

