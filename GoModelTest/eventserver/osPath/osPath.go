package main

import (
	"os"
	"fmt"
	"sort"
)

func main() {
	fileList("./")
}

func fileList(path string)  {
	info, err := os.Lstat(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info.IsDir())
	}
	fmt.Println(info.Name(), info.ModTime(), info.Mode(), info.Size())
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	fmt.Println(names)
	sort.Strings(names)
	fmt.Println(names)
}
