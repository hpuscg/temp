/*
#Time      :  2020/12/18 6:04 下午
#Author    :  chuangangshen@deepglint.com
#File      :  global.go
#Software  :  GoLand
*/
package global

import (
	"go.uber.org/zap"
	"temp/GoModelTest/casbinTest/config"
)

var (
	Server config.Server
	Logger *zap.Logger
)
