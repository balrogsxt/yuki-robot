package robot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/service/event"
	"time"
)

//////////////群组事件相关

//群组接收消息事件
func OnGroupMessageEvent(qqClient *client.QQClient, msgEvent *message.GroupMessage)  {
	defer SetRecover()

	elList := api.ParseElement(msgEvent.Elements)
	//群数据、发送者数据、消息数据列表
	_group := api.Group{
		Id: msgEvent.GroupCode,
		Name: msgEvent.GroupName,
	}
	buildEvent := event.GroupMessageEvent{
		GroupMessageEventHandle:api.GroupMessageEventHandle{
			MsgId: api.GroupMsgId{
				MsgId:api.MsgId{
					Id:msgEvent.Id,
					InternalId:msgEvent.InternalId,
				},
				Group:_group,
			},
			Group:_group,
			MessageList: elList,
			QQ: api.GroupUser{
				QQ:api.QQ{
					Uin: msgEvent.Sender.Uin,
					Name:msgEvent.Sender.Nickname,
				},
				CardName:msgEvent.Sender.CardName,
			},
			SendTime:msgEvent.Time,
		},
	}
	buildEvent.Call()
}
//群消息撤回事件
func OnGroupMessageRecallEvent(qqClient *client.QQClient, msgEvent *client.GroupMessageRecalledEvent)  {
	defer SetRecover()

	buildEvent := event.GroupMessageRecallEvent{
		GroupMessageRecallEventHandle:api.GroupMessageRecallEventHandle{
			GroupId: msgEvent.GroupCode,
			MsgId:msgEvent.MessageId,
			MsgUin: msgEvent.AuthorUin,
			ActionUin:msgEvent.OperatorUin,
			Time:msgEvent.Time,
		},
	}
	buildEvent.Call()
}
//群成员加入事件
func OnGroupUserJoinEvent(qqClient *client.QQClient, msgEvent *client.MemberJoinGroupEvent)  {
	defer SetRecover()

	buildEvent := event.GroupUserJoinEvent{
		GroupUserJoinEventHandle:api.GroupUserJoinEventHandle{
			Group:api.Group{
				Id: msgEvent.Group.Code,
				Name: msgEvent.Group.Name,
			},
			QQ: api.GroupUser{
				QQ:api.QQ{
					Uin:msgEvent.Member.Uin,
					Name: msgEvent.Member.Nickname,
				},
				CardName: msgEvent.Member.CardName,
			},
			Time: int32(time.Now().Unix()),
		},
	}
	buildEvent.Call()
}
//群成员退出事件
func OnGroupUserQuitEvent(qqClient *client.QQClient, msgEvent *client.MemberLeaveGroupEvent)  {
	defer SetRecover()
	buildEvent := event.GroupUserQuitEvent{
		GroupUserQuitEventHandle:api.GroupUserQuitEventHandle{
			Group:api.Group{
				Id: msgEvent.Group.Code,
				Name: msgEvent.Group.Name,
			},
			QQ: api.GroupUser{
				QQ:api.QQ{
					Uin:msgEvent.Member.Uin,
					Name: msgEvent.Member.Nickname,
				},
				CardName: msgEvent.Member.CardName,
			},
			ActionQQ: api.GroupUser{
				QQ:api.QQ{
					Uin:msgEvent.Operator.Uin,
					Name: msgEvent.Operator.Nickname,
				},
				CardName: msgEvent.Operator.CardName,
			},
			Time: int32(time.Now().Unix()),
		},
	}
	buildEvent.Call()
}

