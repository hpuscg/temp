package main

import (
	"github.com/deepglint/muses/device/gpio"
	"os"
	"time"
	"fmt"
)

const (
	LED2J6     = 161
	LED3J4     = 160
	LED_ON     = 1
	LED_OFF    = 0
	SLEEP_TIME = 100
)

const (
	LED_YELLOW = 0
	LED_RED    = 1
)

type LedManager struct {
	leds            []int
	ledsFD          []*os.File
	ledsStatus      []int
	periods         []time.Duration
	intervals       []time.Duration
	periodTimers    []*time.Timer
	intervalTickers []*time.Ticker
}

var _instace *LedManager = nil

func LedInstance() *LedManager {
	if _instace == nil {
		_instace = newLedManager()
	}
	return _instace
}

// control the led
// Param which: which led light you want to control. Currently, only device.LED_YELLOW and device.LED_RED you can control
// Param period: how long the led is ON. 0 means always off, -1 or less than 0 means always on
// Param interval: interval time the led turns to On from Off. 0 means always off, -1 or less than 0 means always on
func (this *LedManager) SetLed(which int, period, interval time.Duration) {

	if which >= len(this.periods) {
		return
	}

	if this.periods[which] != period {
		this.periods[which] = period
	}

	if this.intervals[which] != interval {
		this.intervals[which] = interval

		if t := this.intervalTickers[which]; t != nil {
			t.Stop()
		}
		if interval == 0 || period == 0 {
			device.GpioSetValue(this.ledsFD[which], LED_OFF)
		} else if interval < 0 || period < 0 {
			device.GpioSetValue(this.ledsFD[which], LED_ON)
		} else {
			this.intervalTickers[which] = time.NewTicker(interval + period)
		}

	}

}

func newLedManager() *LedManager {

	ledManager := new(LedManager)
	ledManager.leds = []int{LED2J6, LED3J4}
	ledManager.ledsFD = []*os.File{nil, nil}
	ledManager.periods = []time.Duration{0, 0}
	ledManager.intervals = []time.Duration{0, 0}
	ledManager.periodTimers = []*time.Timer{nil, nil}
	ledManager.intervalTickers = []*time.Ticker{nil, nil}
	ledManager.ledsStatus = []int{LED_OFF, LED_OFF}

	for i, led := range ledManager.leds {
		if err := device.GpioExport(led); err != nil {
			fmt.Printf("Failed to export led %d: %s", led, err)
		}
		if err := device.GpioSetDirection(led, device.DIRECTION_OUT); err != nil {
			fmt.Printf("Failed to set led %d direction: %s", led, err)
		}
		fileFd, err := device.GpioOpen(led)
		if err != nil {
			fmt.Printf("Failed to open led %d: %s", led, err)
		}
		ledManager.ledsFD[i] = fileFd
	}

	return ledManager
}

func (this *LedManager) Start() {
	for i := 0; i < len(this.periods); i++ {
		go this.update(i)
	}
}

func (this *LedManager) update(which int) {
	for {

		if this.intervalTickers[which] != nil {
			select {
			case <-this.intervalTickers[which].C:
				device.GpioSetValue(this.ledsFD[which], LED_ON)
				// create a timer. When the timer escape,
				// turn off the led
				this.ledsStatus[which] = LED_ON
				if t := this.periodTimers[which]; t != nil {
					t.Stop()
				}
				t := time.AfterFunc(this.periods[which], func() {
					device.GpioSetValue(this.ledsFD[which], LED_OFF)
					this.ledsStatus[which] = LED_OFF
				})
				this.periodTimers[which] = t

			default:
				// do nothing
			}
		}

		time.Sleep(time.Duration(SLEEP_TIME) * time.Millisecond)
	}

}
