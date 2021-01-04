package group

import (
	"github.com/balrogsxt/xtbot-go/app/script"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/robot/cq"
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
		api.NewGroupException("运行失败: " + err.Error())
		return true
	}
	event.Group.SendGroupMessageText(cq.At(event.QQ.Uin) + "运行结果:" + ret)
	return true
}
