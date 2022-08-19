package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) {
	var disk DiskStatus
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	fmt.Printf("%+v", disk)
}

func GetCpuInfo() {
	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%.2f\n", cpuInfo[0])
}

func TimeInt() {
	t1 := time.Now().Unix()
	// t2 := t1 / 1e9
	// t3 := t1 - t2*1e9
	fmt.Println(t1)
	/* fmt.Println(t2)
	fmt.Println(t3) */

	fmt.Println(time.Unix(t1, 0))
}

func main() {
	/*
		DiskUsage("/tmp")
		GetCpuInfo()
	*/
	TimeInt()
}
