package api

import "time"

//缓存模块

//用于缓存的接口实现
type RobotCache interface {
	Set(string, interface{}, ...time.Duration) error //设置缓存数据
	Get(string) (string, error)                      //获取缓存数据
	Exists(string) bool                              //判断某个键是否存在
	GetMap(string, string) (string, error)           //获取map结构数据
	SetMap(string, string, interface{}) error        //设置map结构数据
	ExistsMap(string, string) bool                   //判断map指定字段是否存在
}
var cacheModule RobotCache = nil

func InitCache(c RobotCache)  {
	cacheModule = c
}

func GetCache() RobotCache {
	return cacheModule
}


