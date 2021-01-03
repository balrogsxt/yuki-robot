package robot

import "github.com/balrogsxt/xtbot-go/util/logger"

//定义机器人独立服务事件接口
type RobotEventHandle interface {
	Call()
}

func SetRecover() {
	if err := recover(); err != nil {
		logger.Error("[事件异常] %s", err)
	}
}
