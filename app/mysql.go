package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func InitDatabase(config *RobotConfig) error {
	conf := config.Cache.Mysql

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.Name)
	engine, _ = xorm.NewEngine("mysql", dsn)
	err := engine.Ping()
	if err != nil {
		return err
	}
	return nil
}

func GetDb() *xorm.Engine {
	return engine
}
