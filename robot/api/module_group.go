package api

import (
	"github.com/balrogsxt/xtbot-go/app"
)

//针对群组模块的接口
//获取允许接收的群组
func GetGroupAllow() []int64 {
	conf := app.GetRobotConfig()
	return conf.Group.Allow
}
//获取禁止接收的群组
func GetGroupDeny() []int64 {
	conf := app.GetRobotConfig()
	return conf.Group.Deny
}

//群组消息接收事件模块接口
type GroupMessageModule interface {
	Command() string //正则触发命令规则
	Call() bool  //调用触发,如果返回false则继续匹配下一个模块,如果为true则命中模块目标
}


