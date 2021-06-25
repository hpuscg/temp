package main

import (
	"gitlab.deepglint.com/junkaicao/glog"
	"os"
	"os/exec"
)

func main() {
	glog.Config(glog.WithAlsoToStd(true), glog.WithFilePath("./"), glog.WithLevel("info"))
	writeFile()
	// echoCmd()
}

func writeFile() {
	f, err := os.OpenFile("/Users/hpu_scg/gocode/src/temp/GoModelTest/echoTest/export", os.O_WRONLY, 0200)
	if err != nil {
		glog.Infoln(err)
	}
	_, err = f.Write([]byte("5"))
	if err != nil {
		glog.Infoln(err)
	}
}

func echoCmd() {
	cmd := exec.Command("", "echo 5 > /Users/hpu_scg/gocode/src/temp/GoModelTest/echoTest/export")
	// cmd := exec.Command("ls", "-a")
	outPut, err := cmd.Output()
	if err != nil {
		glog.Warningln(err)
	}
	glog.Infoln(string(outPut))
}
