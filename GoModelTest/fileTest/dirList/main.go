package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	in_dir  string
	out_dir string
)

func main() {
	flag.StringVar(&in_dir, "in_dir", ".", "input file dir")
	flag.StringVar(&out_dir, "out_dir", "../temp", "output file dir")
	flag.Parse()
	files, err := GetFiles()
	/* fmt.Println(strings.HasSuffix("test.mp4", "mp4"))
	return */
	fmt.Println(in_dir)
	fmt.Println(out_dir)
	if err != nil {
		panic(err)
	}
	for _, inFile := range files {
		fmt.Println(inFile)
		if !strings.HasSuffix(inFile, "mp4") {
			fmt.Println(inFile)
			continue
		}
		oufFile := path.Join(out_dir, inFile)
		cmd := exec.Command("ffmpeg", "-i", inFile, "-vcodec", "libx264", "-bf", "0", "-vf", "scale=1920:1080", oufFile)
		outCmd, err := cmd.CombinedOutput()
		fmt.Println(string(outCmd))
		fmt.Println(err)
	}
	fmt.Println("end")
}

func GetFiles() (files []string, err error) {
	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}
	err = filepath.Walk(in_dir, walkFunc)
	return
}
