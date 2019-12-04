/*
#Time      :  2019/4/23 上午11:14 
#Author    :  chuangangshen@deepglint.com
#File      :  ntpTest.go
#Software  :  GoLand
*/
package main

import (
	"time"
	"github.com/gpmgo/gopm/modules/log"
	"bytes"
	"os/exec"
	"fmt"
	"errors"
	"flag"
	"github.com/deepglint/muses/autobot/utils/strtool"
	"strings"
	"strconv"
)

func main() {
	var ip string
	flag.StringVar(&ip, "ip", "192.168.101.238", "ntp server ip")
	flag.Parse()
	// syncTime(time.Duration(30)*time.Second, "synctime.sh", ip)
	cmdNtp(ip)
}

func cmdNtp(ip string) {
	/*command := "ntpdate " + ip
	cmd := exec.Command("/bin/bash", "-c", command)
	respBytes, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	resp := string(respBytes)
	fmt.Println(resp)*/
	var ret time.Time
	query, err := CmdOutNoLn("ntpdate", "-q", ip)
	if err != nil {
		fmt.Println("===42===", err)
	}
	query, err = CmdOutNoLn("ntpdate", "-q", ip)
	if err != nil {
		fmt.Println("===46===", err)
	}
	offsetIndex := strings.LastIndex(query, "offset")
	secIndex := strings.LastIndex(query, "sec")

	duration := strings.TrimSpace(query[offsetIndex+6 : secIndex])

	if strings.HasPrefix(duration, "-") {
		mis := strings.TrimLeft(duration, "-")
		misFloat6, _ := strconv.ParseFloat(mis, 64)
		ret = time.Now().Add(time.Duration(-misFloat6) * time.Second)
	} else {
		misFloat64, _ := strconv.ParseFloat(duration, 64)
		ret = time.Now().Add(time.Duration(misFloat64) * time.Second)
	}
	fmt.Println("===57==", ret)

}

func syncTime(timeout time.Duration, name string, arg ...string) {
	fmt.Println("===command 40====name=", name, "arg=", arg, "timeout", timeout)
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		fmt.Println("=====command 47====", err)
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		//timeout
		if err = cmd.Process.Kill(); err != nil {
			log.Error("failed to kill: %s, error: %s", cmd.Path, err)
		}
		go func() {
			<-done // allow goroutine to exit
		}()
		msg := fmt.Sprintf("process:%s killed because of timeout", cmd.Path)
		err = errors.New(msg)
		fmt.Println("===command 67====err=", err)
		// return "", err, true
	case err = <-done:
		fmt.Println("===command 70===err=", err, "out=", out.String())
		// return out.String(), err, false
	}
}

func CmdOutNoLn(name string, arg ...string) (out string, err error) {
	out, err = CmdOut(name, arg...)
	if err != nil {
		return
	}

	return strtool.TrimRightSpace(string(out)), nil
}

func CmdOut(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

