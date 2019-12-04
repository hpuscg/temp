package main

import (
	"time"

	"github.com/deepglint/muses/device/led"
)

type LedsControlModel struct {
	AlternateTime    int
	LedControlModels []LedControlModel
}

type LedControlModel struct {
	Which    int
	Period   time.Duration
	Interval time.Duration
}

const (
	FACTORY_DEFAULT_WEIGHT = 0
	ALL_READY_WEIGHT       = 1
	AP_MODE_WEIGHT         = 2
	//
	SOFTWARE_ERROR_WEIGHT = 3
	CONFIG_ERROR_WEIGHT   = 4
	IP_ERROR_WEIGHT       = 5
	//
	HARDWARE_ERROR_WEIGHT = 6
	USB_ERROR_WEIGHT      = 7
	TF_CARD_ERROR_WEIGHT  = 8
	//
	CLIENT_SIGNAL_WEIGHT = 9
)

var (
	FACTORY_DEFAULT_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			},
		},
	}
	ALL_READY_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(1000) * time.Millisecond,
				time.Duration(1000) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(1000) * time.Millisecond,
				time.Duration(1000) * time.Millisecond,
			},
		},
	}
	//
	AP_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(-1) * time.Millisecond,
				time.Duration(-1) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(-1) * time.Millisecond,
				time.Duration(-1) * time.Millisecond,
			},
		},
	}
	//
	SOFTWARE_ERROR_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(100) * time.Millisecond,
				time.Duration(2000) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			},
		},
	}
	CONFIG_ERROR_MODE = LedsControlModel{
		100,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(100) * time.Millisecond,
				time.Duration(100) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			},
		},
	}
	IP_ERROR_MODE = LedsControlModel{
		800,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(-1) * time.Millisecond,
				time.Duration(-1) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			},
		},
	}
	//
	HARDWARE_ERROR_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(100) * time.Millisecond,
				time.Duration(2000) * time.Millisecond,
			},
		},
	}
	USB_ERROR_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(100) * time.Millisecond,
				time.Duration(100) * time.Millisecond,
			},
		},
	}
	TF_CARD_ERROR_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(0) * time.Millisecond,
				time.Duration(0) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(-1) * time.Millisecond,
				time.Duration(-1) * time.Millisecond,
			},
		},
	}
	CLIENT_SIGNAL_MODE = LedsControlModel{
		0,
		[]LedControlModel{
			LedControlModel{
				led.LED_YELLOW,
				time.Duration(200) * time.Millisecond,
				time.Duration(200) * time.Millisecond,
			}, LedControlModel{
				led.LED_RED,
				time.Duration(200) * time.Millisecond,
				time.Duration(200) * time.Millisecond,
			},
		},
	}
)
