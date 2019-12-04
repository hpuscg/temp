package main

import (
	"fmt"
	"runtime"
)

func main() {
	systemInfo()
}

func systemInfo() {
	fmt.Println(runtime.GOOS)
}

