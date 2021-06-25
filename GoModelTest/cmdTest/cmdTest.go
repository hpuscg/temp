package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.OpenFile("/Users/hpu_scg/gocode/src/temp/GoModelTest/cmdTest/test.txt",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.WriteString("0")
	if err != nil {
		fmt.Println(err)
	}

}
