package main

import (
	"time"
)

func main() {
	Lp := time.Duration(-1) * time.Millisecond
	Liva := time.Duration(-1) * time.Millisecond
	Lwhich := 1
	Yp := time.Duration(-1) * time.Millisecond
	Yiva := time.Duration(-1) * time.Millisecond
	Ywhich := 0
	InitLedMode()
	go LedTimer(60)
	// go InitButtonListener()
	changLedUpAndDown(Lwhich, Lp, Liva)
	changLedUpAndDown(Ywhich, Yp, Yiva)
}

func changLedUpAndDown(w int, p, i time.Duration) {
	var ledManager *LedManager
	ledManager = LedInstance()
	ledManager.SetLed(w, p, i)
}