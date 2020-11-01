package event

//群聊事件

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/balrogsxt/xtbot-go/util/msg"
)

//收到群成员发送消息事件
func OnGroupMessageEvent(event *msg.GroupHandle) {

	//m := handle.BuildMsg()
	//m.Add(handle.MsgBuild.LocalImage("./test/1.jpg"))                            //发送本地图片
	//m.Add(handle.MsgBuild.ImageId("{D67D3BFA-F98E-9D32-85DF-FFEBFC1ABE18}.jpg")) //发送图片ID
	//m.Add(handle.MsgBuild.At(event.Sender.Uin))
	//at发送者
	//handle.SendMsg(m)
	text := fmt.Sprintf("[type=at,value=%d]", event.Event.Sender.Uin)
	//img := "./test/1.jpg"
	//text += fmt.Sprintf("[type=image,value=%s]",img)
	//text += fmt.Sprintf("[type=image,value=%s]","{E4AD8A49-C2E8-1287-E5B3-559F7E5376AF}.PNG")
	event.SendMessage(text)

}

//收到群聊消息撤回事件
func OnGroupMessageRecallEvent(qqClient *client.QQClient, event *client.GroupMessageRecalledEvent) {

}
