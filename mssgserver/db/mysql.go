package db

import (
	"log"
	"mssgserver/config"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func TestDB() {
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		log.Println("数据库配置缺失", err)
		panic(err)
	}
	Engine, err = xorm.NewEngine("mysql", "root:123@/test?charset=utf8")

}
