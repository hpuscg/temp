package main

import (
	// "unsafe"
	"fmt"
	// "github.com/deepglint/muses/eventserver/models"
	"time"
	// "os"
	// "os/exec"
)

/*
func main() {
	for {
		var timeout int
		go func() 
			for true {
				// var cmd *exec.Cmd
				switch timeout{
				case 2:
					pid := os.Getpid()
					fmt.Println(pid)
					os.Exit(1)
				case 1:

				}
			}
			time.Sleep(2*time.Second)
			timeout = 1
		}()
		go func() {
			time.Sleep(3*time.Second)
			timeout = 2
		}()
		Loop:
			for true {
				switch timeout {
				case 1:
					fmt.Println("time out")
					timeout = 0
					break Loop
				case 2:
					fmt.Println("sussess")
					timeout = 0
					break Loop
				default:
					ppid := os.Getpid()
					fmt.Println(ppid)
					fmt.Println(timeout)
					continue
				}
			}
	}
}
*/

/*
func main() {
	time.Sleep(5000000000)
	fmt.Println("sleep 5 seconds")
}
*/

/*
func main() {
	var sliceTest = make([]models.Event, 5)
	a := sliceTest[0]
	if a.SensorId == "" {
		fmt.Println("no thing")
	}else{
		fmt.Println(unsafe.Sizeof(a))
	}
	fmt.Println(a)
}
*/



func main() {
	fmt.Println("111")
	tc := time.After(3*time.Second)
	fmt.Println("222")
	fmt.Println("333")
	<-tc
	fmt.Println("444")
}


