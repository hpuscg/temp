package main

import (
	"strconv"
	"time"
)

var (
	led_mode    = 0
	led_manager *LedManager
)

func InitLedMode() {
	led_manager = LedInstance()
	led_manager.Start()

	led_str := "1"

	led_mode, _ = strconv.Atoi(led_str)

	// glog.Infof("------current led mode: %d------", led_mode)
	showLedMode()
}

func LedTimer(interval int) {
	timer := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-timer.C:
			go updateLedMode()
		}
	}
}

func updateLedMode() {
	cur_led_mode := 0
	/*
		// default 0 / 1
		led_str, err := bumble.BumbleConfig.EtcdClient.GetValue(models.FACTORY_DEFAULT)
		if err != nil {
			glog.Errorln(err)
		}
		cur_led_mode, err = strconv.Atoi(led_str)
		if err != nil {
			glog.Errorln(err)
		}
	*/
	/////////////////////////////////////////
	// ap mode 2
	var CurrentWlanMode string
	if CurrentWlanMode == "ap" {
		cur_led_mode = 2
	}
	// software 3

	/*if !bumble.HostStatus.HostStatus {
		cur_led_mode = models.SOFTWARE_ERROR_WEIGHT
	}*/

	// config 4
	// TODO: check all config
	/*svraddr, err := bumble.BumbleConfig.EtcdClient.GetValue(models.SERVER_ADDRESS)
	if err != nil {
		// glog.Errorln(err)
		fmt.Println(err)
		// cur_led_mode = models.CONFIG_ERROR_WEIGHT
	}*/

	// ip 5
	/*if svraddr != "" && !ping.Ping(svraddr, 10) {
		cur_led_mode = models.IP_ERROR_WEIGHT
	}*/

	/////////////////////////////////////////
	// hardware 6
	// TODO: check hardware

	// usb 7
	/*if !funcs.CheckPrimeSenseExist() {
		cur_led_mode = models.USB_ERROR_WEIGHT
	}*/

	/*
	// tf 8
	// TODO: check tf card
	if TFManager == nil {
		cur_led_mode = models.TF_CARD_ERROR_WEIGHT
	} else {
		tf, _, _ := tfcard.CheckMount(TFManager.TfPartation)
		if !tf {
			cur_led_mode = models.TF_CARD_ERROR_WEIGHT
		}
	}
	*/

	// glog.Infof("------old led mode: %d------", led_mode)
	// glog.Infof("------new led mode: %d------", cur_led_mode)
	if led_mode != cur_led_mode {
		led_mode = cur_led_mode
		showLedMode()
	}
}

func showLedMode() {
	/*type LedControlModel struct {
		Which    int
		Period   time.Duration
		Interval time.Duration
	}

	type LedsControlModel struct {
		AlternateTime    int
		LedControlModels []LedControlModel
	}*/
	var lcms LedsControlModel
	switch led_mode {
	case FACTORY_DEFAULT_WEIGHT:
		lcms = FACTORY_DEFAULT_MODE
	case ALL_READY_WEIGHT:
		lcms = ALL_READY_MODE
	case AP_MODE_WEIGHT:
		lcms = AP_MODE
	case SOFTWARE_ERROR_WEIGHT:
		lcms = SOFTWARE_ERROR_MODE
	case CONFIG_ERROR_WEIGHT:
		lcms = CONFIG_ERROR_MODE
	case IP_ERROR_WEIGHT:
		lcms = IP_ERROR_MODE
	case HARDWARE_ERROR_WEIGHT:
		lcms = HARDWARE_ERROR_MODE
	case USB_ERROR_WEIGHT:
		lcms = USB_ERROR_MODE
	case TF_CARD_ERROR_WEIGHT:
		lcms = TF_CARD_ERROR_MODE
	default:
		break
	}
	alternate := lcms.AlternateTime
	for _, lcm := range lcms.LedControlModels {
		led_manager.SetLed(lcm.Which, lcm.Period, lcm.Interval)
		time.Sleep(time.Duration(alternate) * time.Millisecond)
	}
}
