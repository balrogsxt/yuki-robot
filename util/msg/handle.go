package msg

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"regexp"
	"strconv"
)

type MsgList struct {
	List []message.IMessageElement
}

//添加消息元素
func (this *MsgList) Add(element message.IMessageElement) {
	this.List = append(this.List, element)
}

//提供一些群组快捷操作方法
type GroupHandle struct {
	Handle   *client.QQClient
	MsgBuild *GroupMessageBuilder
	Event    *message.GroupMessage
}

//快速发送群组消息
func (this *GroupHandle) SendMsg(msg *MsgList) *message.GroupMessage {
	sm := new(message.SendingMessage)
	for _, item := range msg.List {
		if item != nil {
			sm.Append(item)
		}
	}
	return this.Handle.SendGroupMessage(this.MsgBuild.Event.GroupCode, sm)
}

//快速发送解析消息
func (this *GroupHandle) SendMessage(text string) *message.GroupMessage {
	/**
	普通文字普通文字普通文字
	[type=at,value=2289453456] = at指定
	[type=at,value=0] = at全体
	[type=image,value={E4AD8A49-C2E8-1287-E5B3-559F7E5376AF}.PNG] = 发送图片ID
	[type=image,value=./test/1.jpg] = 发送本地图片
	*/
	//text = "这是一段话[type=at,value=123456][type=at,value=123456]这是后面一段话"

	reg, _ := regexp.Compile("\\[type=.*?,value=.*?\\]|\\D")
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
	return this.Handle.SendGroupMessage(this.MsgBuild.Event.GroupCode, els)
}
func (this *GroupHandle) BuildMsg() *MsgList {
	return new(MsgList)
}
