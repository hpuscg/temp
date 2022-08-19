package main

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

func initLog() {
	var (
		f   *os.File
		err error
	)
	logFile := "config.log"
	if _, err = os.Stat(logFile); os.IsNotExist(err) {
		if f, err = os.Create(logFile); err != nil {
			panic(err)
		}
	} else {
		if f, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766); err != nil {
			panic(err)
		}
	}
	logger = log.New(io.MultiWriter(os.Stdout, f), "", log.Lshortfile)
}
