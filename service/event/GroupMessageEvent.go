package event

import (
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/app/db"
	"github.com/balrogsxt/xtbot-go/app/script"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/service/module/group"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"regexp"
	"strings"
)

var groupModules []api.GroupMessageModule

func init() {
	//初始化已封装的群组模块
	groupModules = make([]api.GroupMessageModule, 0)

	groupModules = append(groupModules, new(group.AddKey)) //添加关键字回复
	groupModules = append(groupModules, new(group.DelKey)) //删除关键字回复

}

//群组消息接收事件
type GroupMessageEvent struct {
	api.GroupMessageEventHandle
}

func (this *GroupMessageEvent) Call() {
	simpleMsg := api.ToString(this.MessageList)
	if len(simpleMsg) == 0 {
		return
	}
	go this.saveMsg()

	//是否允许接收
	isAllow := false
	for _, val := range api.GetGroupAllow() {
		if this.Group.Id == val || this.Group.Id == 0 {
			isAllow = true
			break
		}
	}
	if !isAllow {
		return
	}

	//是否被屏蔽
	isDeny := false
	for _, val := range api.GetGroupDeny() {
		if this.Group.Id == val {
			isDeny = true
			break
		}
	}
	if isDeny {
		return
	}
	logger.Info("[群组消息] [%s](%d): %s", this.Group.Name, this.Group.Id, simpleMsg)

	//处理群组模块能力

	//处理go内置机器人功能
	var isCall bool = false
	for _, m := range groupModules {
		flag, _ := regexp.MatchString(m.Command(), simpleMsg)
		if flag {

			reg, _ := regexp.Compile(m.Command())
			value := reg.ReplaceAllString(simpleMsg, "")
			value = strings.Trim(value, " ")
			if m.Call(value, this.GroupMessageEventHandle) {
				isCall = true
				break //执行成功并且命中目标
			}
		}
	}
	if isCall {
		return
	}

	//提供给js虚拟机处理能力
	gmeFile := "./plugins/js/GroupMessageEvent.js"
	if util.IsFile(gmeFile) {
		js := script.NewJs()
		js.SetVars("event", map[string]interface{}{
			"msg_id":    this.MsgId.MsgId.Id,
			"group":     this.Group.Id,
			"qq":        this.QQ.Uin,
			"send_time": this.SendTime,
			"msg_text":  api.ToString(this.MessageList),
			"msg_json":  api.ToJson(this.MessageList),
		})
		if err := js.RunFile("./plugins/js/GroupMessageEvent.js"); err != nil {
			logger.Warning("[群组模块] Js虚拟机运行失败: %s", err.Error())
		}
	}

}

//保存消息到数据库
func (this *GroupMessageEvent) saveMsg() {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("[群组消息保存] 保存消息发生错误: %s", err)
		}
	}()
	msg := new(db.GroupMsg)
	msg.MsgId = this.MsgId.MsgId.Id
	msg.InternalId = this.MsgId.InternalId
	msg.Group = this.Group.Id
	msg.QQ = this.QQ.Uin
	msg.SendTime = this.SendTime
	msg.MsgText = api.ToString(this.MessageList)
	msg.MsgJson = api.ToJson(this.MessageList)

	if _, err := app.GetDb().InsertOne(msg); err != nil {
		logger.Warning("[群组消息保存] 保存消息发生错误: %s", err)
	}

}
