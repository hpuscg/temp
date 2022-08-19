package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

func main() {
	DealWithFiles("/Users/hpu_scg/gocode/src/temp/GoModelTest/fileTest/removeOldFile/test/", "")
}

type ByModTime []os.FileInfo

func (fis ByModTime) Len() int {
	return len(fis)
}
func (fis ByModTime) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}
func (fis ByModTime) Less(i, j int) bool {
	return fis[i].ModTime().Before(fis[j].ModTime())
}

// 根目录下的文件按时间大小排序, 从远到近
func SortFile(path, name string) (files ByModTime) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fis, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	files = make(ByModTime, len(fis)+10)
	j := 0
	for i, v := range fis {
		if strings.Contains(fis[i].Name(), name) {
			files[j] = v
			j++
		}
	}
	files = files[:j]
	sort.Sort(ByModTime(files))
	return
}

// 返回当下时间的文件, 并删除大于 5 个的文件, 删除最早的, 如果目录下没有文件, 就自动创建
func DealWithFiles(pathdir, name string) {
	files := SortFile(pathdir, name)
	// fmt.Println(path + files[len(files)-1].Name())
	if len(files) > 5 {
		for k := range files[5:] {
			err := os.Remove(path.Join(pathdir, files[k].Name()))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
