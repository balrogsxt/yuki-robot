package msg

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
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
}

//快速发送群组消息
func (this *GroupHandle) SendMessage(msg *MsgList) int32 {
	sm := new(message.SendingMessage)
	for _, item := range msg.List {
		if item != nil {
			sm.Append(item)
		}
	}
	m := this.Handle.SendGroupMessage(this.MsgBuild.Event.GroupCode, sm)
	return m.Id
}
func (this *GroupHandle) BuildMsg() *MsgList {
	return new(MsgList)
}
