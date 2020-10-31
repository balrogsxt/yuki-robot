package event

//群聊事件

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/util/msg"
)

//收到群成员发送消息事件
func OnGroupMessageEvent(handle *msg.GroupHandle, event *message.GroupMessage) {
	m := handle.BuildMsg()
	m.Add(handle.MsgBuild.LocalImage("./test/1.jpg"))                            //发送本地图片
	m.Add(handle.MsgBuild.ImageId("{D67D3BFA-F98E-9D32-85DF-FFEBFC1ABE18}.jpg")) //发送图片ID
	m.Add(handle.MsgBuild.At(event.Sender.Uin))                                  //at发送者
	handle.SendMessage(m)
}

//收到群聊消息撤回事件
func OnGroupMessageRecallEvent(qqClient *client.QQClient, event *client.GroupMessageRecalledEvent) {

}
