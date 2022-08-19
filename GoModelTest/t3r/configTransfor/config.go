package main

import "github.com/spf13/viper"

var (
	configViper    *viper.Viper
	configFileName = "config.yaml"
)

func initConfigViper() {
	configViper = viper.New()
	configViper.SetConfigFile(configFileName)
	if err := configViper.ReadInConfig(); err != nil {
		panic(err)
	}
}
