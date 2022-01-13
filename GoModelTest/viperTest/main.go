package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	readConfig()
}

func readConfig() {
	config := viper.New()
	config.AddConfigPath("./")
	config.SetConfigName("fgConfig")
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	config.SetDefault("deviceinfo.deviceName", "123")
	// fmt.Println(config.AllSettings())
	fmt.Println(config.Get("deviceinfo"))
	config.Set("deviceinfo.deviceName", "F57")
	if err := config.WriteConfig(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Get("deviceinfo"))
	fmt.Println(config.AllSettings())
}
