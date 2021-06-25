package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	value = "0"
)

func main() {
	flag.StringVar(&value, "value", "0", "file value")
	flag.Parse()
	err := writeFile("/sys/class/gpio/gpio5/value", value, true)
	if err != nil {
		fmt.Println(err)
	}
}

func writeFile(fileName, value string, enable bool) error {
	if !enable {
		return nil
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	_, err = f.WriteString(value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
