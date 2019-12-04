package main

import(
	"flag"
	"fmt"
)

var music_file *string = flag.String("file", "musicfile", "Use -file<filesource>")

func main() {
	flag.Parse()
	fmt.Println(*music_file)
}

