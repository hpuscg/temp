/*
#Time      :  2020/12/18 5:58 下午
#Author    :  chuangangshen@deepglint.com
#File      :  casbin.go
#Software  :  GoLand
*/
package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"strings"
	"temp/GoModelTest/casbinTest/global"
)

func Casbin() *casbin.Enforcer {
	mysql := global.Server.Mysql
	a, _ := gormadapter.NewAdapter(global.Server.System.DbType, mysql.Username+":"+
		mysql.Password+"@("+mysql.Path+")/"+mysql.Dbname, true)
	e, _ := casbin.NewEnforcer(global.Server.Casbin.ModelPath, a)
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = e.LoadPolicy()
	return e
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
