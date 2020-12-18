/*
#Time      :  2020/12/18 6:42 下午
#Author    :  chuangangshen@deepglint.com
#File      :  config.go
#Software  :  GoLand
*/
package config

type Server struct {
	Mysql Mysql
	System System
	Casbin Casbin
}
