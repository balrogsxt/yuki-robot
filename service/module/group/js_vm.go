package group

import (
	"github.com/balrogsxt/xtbot-go/app/script"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"time"
)

//回复关键词添加
type JsVm struct {
}

func (JsVm) Command() string {
	return "^#js"
}
func (this *JsVm) Call(value string, event api.GroupMessageEventHandle) bool {
	js := script.NewJs()
	ret, err := js.RunCode(value, 30000)
	if err != nil {
		this.SendErrorMessage(event.Group.Id, "运行失败: "+err.Error())
		return true
	}
	event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "运行结果:" + ret)
	return true
}

//运行错误的发送消息,延迟撤回
func (this *JsVm) SendErrorMessage(groupId int64, text string) {
	m := api.SendGroupMessageText(groupId, text)
	go func() {
		time.Sleep(time.Millisecond * 3000)
		api.RecallGroupMessage(groupId, m.MsgId.Id)
	}()
}
