package event

import (
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
)
//消息撤回事件
type GroupMessageRecallEvent struct {
	api.GroupMessageRecallEventHandle
}

func (this *GroupMessageRecallEvent) Call()  {
	logger.Info("[消息撤回事件触发]")
}