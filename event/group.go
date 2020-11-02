package event

//群聊事件

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/util/cache"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/balrogsxt/xtbot-go/util/msg"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type GroupMessageEvent struct {
	Handle   *client.QQClient
	MsgBuild *msg.GroupMessageBuilder
	Event    *message.GroupMessage
}

//收到群成员发送消息事件
func OnGroupMessageEvent(moduleList []GroupModule, event *GroupMessageEvent) {
	text := strings.Trim(event.Event.ToString(), " ")
	if len(text) == 0 {
		return
	}
	for _, m := range moduleList {
		if 0 >= len(m.GetCommand()) {
			continue
		}
		//命令规则  ^command arg arg
		if flag, err := regexp.Match(fmt.Sprintf("^%s", m.GetCommand()), []byte(text)); !flag || err != nil {
			continue
		}
		//闭包调用
		go func(m GroupModule) {
			defer func() {
				err := recover()
				if err != nil {
					logger.Error("【%s】模块内部发生错误: %s", m.GetName(), err)
				}
			}()
			m.InitHandle(m.GetCommand(), event) //初始化模块
			m.Handle()
		}(m)
		break
	}

}

//收到群聊消息撤回事件
func OnGroupMessageRecallEvent(qqClient *client.QQClient, event *client.GroupMessageRecalledEvent) {

}

//群组聊天模块基类
type GroupModule interface {
	GetName() string                       //模块名称
	GetCommand() string                    //触发命令
	InitHandle(string, *GroupMessageEvent) //命中后初始化模块
	Handle()                               //触发方法
}
type GroupModuleHandle struct {
	Event    *message.GroupMessage    //群组事件
	Client   *client.QQClient         //核心QQ客户端
	Cache    cache.XtCache            //缓存模块
	MsgBuild *msg.GroupMessageBuilder //群组消息构造器
	Args     []string                 //命令解析后的参数
}

//实现必要的方法
func (this *GroupModuleHandle) GetName() string {
	return "未定义模块名称"
}
func (this *GroupModuleHandle) GetCommand() string {
	return "~未定义的模块"
}
func (this *GroupModuleHandle) Handle() {
}

//默认的初始化模块
func (this *GroupModuleHandle) InitHandle(command string, event *GroupMessageEvent) {
	this.Event = event.Event
	this.Client = event.Handle
	this.Cache = event.MsgBuild.Cache
	this.MsgBuild = event.MsgBuild
	//解析命令参数
	msg := event.Event.ToString()
	start := utf8.RuneCountInString(command)
	this.Args = strings.Split(strings.Trim(msg[start:], " "), " ")
}

//快速发送解析消息
func (this *GroupModuleHandle) SendLineString(text []string) *message.GroupMessage {
	msg := ""
	for _, item := range text {
		msg += fmt.Sprintf("%s\n", item)
	}
	return this.SendMessage(msg)
}
func (this *GroupModuleHandle) SendMessage(format string, args ...interface{}) *message.GroupMessage {
	text := fmt.Sprintf(format, args...)
	/**
	普通文字普通文字普通文字
	[type=at,value=2289453456] = at指定
	[type=at,value=0] = at全体
	[type=image,value={E4AD8A49-C2E8-1287-E5B3-559F7E5376AF}.PNG] = 发送图片ID
	[type=image,value=./test/1.jpg] = 发送本地图片
	*/
	//text = "这是一段话[type=at,value=123456][type=at,value=123456]这是后面一段话"

	reg, _ := regexp.Compile("\\[type=.*?,value=.*?\\]|(.|\n){1}")
	list := reg.FindAllString(text, -1)
	els := new(message.SendingMessage)

	tmpText := ""
	for _, txt := range list {
		regex, _ := regexp.Compile("\\[type=(.*)?,value=(.*)?\\]")
		item := regex.FindStringSubmatch(txt)
		if len(item) >= 3 {
			//如果匹配成功到了,结算文字
			if len(tmpText) > 0 {
				els.Append(this.MsgBuild.Text(tmpText))
				tmpText = ""
			}
			//判断类型
			_type := item[1]
			_value := item[2]

			var elem message.IMessageElement
			switch _type {
			case "at": //at对方
				elem = func(v string) message.IMessageElement {
					res, err := strconv.ParseInt(v, 10, 64)
					if err == nil { //为nil时才添加,防止出现为0触发at全体
						return this.MsgBuild.At(res)
					}
					return nil
				}(_value)
				break
			case "image":
				elem = func(v string) message.IMessageElement {
					//判断是否是图片ID格式
					flag, _ := regexp.Match("^\\{[a-zA-Z0-9-]+\\}\\.(PNG|JPG|JPEG|GIF|BMP|WEBP)$", []byte(v))
					if flag {
						return this.MsgBuild.ImageId(v) //载入id图片
					} else {
						return this.MsgBuild.LocalImage(v) //上传本地图片
					}
				}(_value)
				break
			}
			if elem != nil {
				els.Append(elem)
			}
		} else {
			tmpText += txt
		}
	}
	//再次判断是否需要结算文字
	if len(tmpText) > 0 {
		els.Append(this.MsgBuild.Text(tmpText))
		tmpText = ""
	}
	logger.Info("[发送消息] %s", strings.ReplaceAll(this.toString(els), "\n", " "))
	return this.Client.SendGroupMessage(this.MsgBuild.Event.GroupCode, els)
}
func (this *GroupModuleHandle) BuildMsg() *MsgList {
	return new(MsgList)
}

//快速发送群组消息
func (this *GroupModuleHandle) SendMsg(msg *MsgList) *message.GroupMessage {
	sm := new(message.SendingMessage)
	for _, item := range msg.List {
		if item != nil {
			sm.Append(item)
		}
	}
	logger.Info("[发送消息] %s", strings.ReplaceAll(this.toString(sm), "\n", " "))
	return this.Client.SendGroupMessage(this.MsgBuild.Event.GroupCode, sm)
}
func (this *GroupModuleHandle) toString(list *message.SendingMessage) (res string) {
	for _, elem := range list.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			res += e.Content
		case *message.ImageElement:
			res += "[Image:" + e.Filename + "]"
		case *message.FaceElement:
			res += "[" + e.Name + "]"
		case *message.GroupImageElement:
			res += "[Image: " + e.ImageId + "]"
		case *message.AtElement:
			res += e.Display
		case *message.RedBagElement:
			res += "[RedBag:" + e.Title + "]"
		case *message.ReplyElement:
			res += "[Reply:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		}
	}
	return
}

type MsgList struct {
	List []message.IMessageElement
}

//添加消息元素
func (this *MsgList) Add(element message.IMessageElement) {
	this.List = append(this.List, element)
}
