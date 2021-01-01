package event

import (
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
)

//群成员加入事件
type GroupUserJoinEvent struct {
	api.GroupUserJoinEventHandle
}
func (this *GroupUserJoinEvent) Call() {
	logger.Info("[群成员加入] [%s](%d) -> %s(%d)",this.Group.Name,this.Group.Id,this.QQ.GetDisplayName(),this.QQ.Uin)
}