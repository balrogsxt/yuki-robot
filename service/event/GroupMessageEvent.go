package event

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/app/db"
	"github.com/balrogsxt/xtbot-go/app/script"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/robot/cq"
	"github.com/balrogsxt/xtbot-go/service/module/group"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"regexp"
	"strings"
	"time"
)

var groupModules []api.GroupMessageModule
var groupJsModules []api.GroupJsMessageModule

func init() {
	//初始化已封装的群组模块
	groupModules = make([]api.GroupMessageModule, 0)

	groupModules = append(groupModules, new(group.AddKey)) //添加关键字回复
	groupModules = append(groupModules, new(group.DelKey)) //删除关键字回复
	groupModules = append(groupModules, new(group.JsVm))   //运行js虚拟机

	//js模块
	groupJsModules = make([]api.GroupJsMessageModule, 0)
	groupJsModules = append(groupJsModules, api.GroupJsMessageModule{
		Command: "^test",
		File:    "test.js",
	})

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
	logger.Info("[群组消息] [%s](%d): %s", this.Group.Name, this.Group.Id, strings.ReplaceAll(simpleMsg, "\n", "<br/>"))
	//处理群组模块能力

	//处理go内置机器人功能
	var isCall bool = false
	for _, m := range groupModules {
		flag, _ := regexp.MatchString(m.Command(), simpleMsg)
		if flag {
			reg, _ := regexp.Compile(m.Command())
			value := reg.ReplaceAllString(simpleMsg, "")
			value = strings.Trim(value, " ")
			isCall = func() bool {
				//错误捕捉
				defer func() {
					if err := recover(); err != nil {
						exception, has := err.(api.GroupMessageModuleException)
						exmsg := cq.At(this.QQ.Uin) + "模块执行失败~"
						if has {
							exmsg = exception.Message
						}
						m := api.SendGroupMessageText(this.Group.Id, exmsg)
						go func() {
							time.Sleep(time.Millisecond * 5000)
							api.RecallGroupMessage(this.Group.Id, m.MsgId.Id)
						}()
					}
				}()

				if m.Call(value, this.GroupMessageEventHandle) {
					return true
				}
				return false
			}()
			if isCall {
				break
			}

		}
	}
	if isCall {
		return
	}

	//js处理能力
	jsDir := "./plugins/js/"
	for _, m := range groupJsModules {
		flag, _ := regexp.MatchString(m.Command, simpleMsg)
		if flag {
			reg, _ := regexp.Compile(m.Command)
			value := reg.ReplaceAllString(simpleMsg, "")
			value = strings.Trim(value, " ")
			file := fmt.Sprintf("%s%s", jsDir, m.File)
			if util.IsFile(file) {
				js := script.NewJs()
				js.SetVars("event", map[string]interface{}{
					"value":     value,
					"msg_id":    this.MsgId.MsgId.Id,
					"group":     this.Group.Id,
					"qq":        this.QQ.Uin,
					"send_time": this.SendTime,
					"msg_text":  api.ToString(this.MessageList),
					"msg_json":  api.ToJson(this.MessageList),
				})
				ret, err := js.RunFile(file)
				if err != nil {
					logger.Warning("[群组模块] Js虚拟机运行失败: %s", err.Error())
				}
				//脚本模块运行结果在控制台输出
				fmt.Println("运行结果:" + ret)
				isCall = true
				break //执行成功并且命中目标
			}
		}
	}
	if isCall {
		return
	}

	//处理群组自定义关键字回复功能等
	this.last(simpleMsg)
}

func (this *GroupMessageEvent) last(text string) {

	//处理关键字回复

	result := db.GroupReply{}

	//随机指定一条回复,反向正向like查询关键词
	has, err := app.GetDb().Where("`key` like ? or ? like CONCAT('%',`key`,'%')", "%"+text+"%", text).Cols("reply").OrderBy("rand()").Get(&result)
	if err != nil {
		logger.Warning("[群组回复] 查询群组关键词回复失败: %s", err.Error())
	}
	if !has {
		return //没有回复词条
	}
	this.Group.SendGroupMessageText(result.Reply)

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
