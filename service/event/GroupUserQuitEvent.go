package event

import (
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
)

//群成员退出事件
type GroupUserQuitEvent struct {
	api.GroupUserQuitEventHandle
}
func (this *GroupUserQuitEvent) Call() {
	logger.Info("[群成员退出] [%s](%d) -> %s(%d)",this.Group.Name,this.Group.Id,this.QQ.GetDisplayName(),this.QQ.Uin)
}