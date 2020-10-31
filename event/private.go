package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

//私聊消息事件

func OnPrivateMessageEvent(client *client.QQClient, event *message.PrivateMessage) {
	a := message.SendingMessage{}
	a.Append(&message.TextElement{Content: "收到私聊消息:" + event.ToString()})
	client.SendPrivateMessage(event.Sender.Uin, &a)
}
