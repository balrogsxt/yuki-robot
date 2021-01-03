package robot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/service/event"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/imroc/req"
	"io/ioutil"
	"os"
	"time"
)

//////////////群组事件相关

//群组接收消息事件
func OnGroupMessageEvent(qqClient *client.QQClient, msgEvent *message.GroupMessage) {
	defer SetRecover()
	//解析数据内图片、语音消息到缓存
	go RunCacheMsg(msgEvent.Elements)
	elList := api.ParseElement(msgEvent.Elements)
	//群数据、发送者数据、消息数据列表
	_group := api.Group{
		Id:   msgEvent.GroupCode,
		Name: msgEvent.GroupName,
	}
	buildEvent := event.GroupMessageEvent{
		GroupMessageEventHandle: api.GroupMessageEventHandle{
			MsgId: api.GroupMsgId{
				MsgId: api.MsgId{
					Id:         msgEvent.Id,
					InternalId: msgEvent.InternalId,
				},
				Group: _group,
			},
			Group:       _group,
			MessageList: elList,
			QQ: api.GroupUser{
				QQ: api.QQ{
					Uin:  msgEvent.Sender.Uin,
					Name: msgEvent.Sender.Nickname,
				},
				CardName: msgEvent.Sender.CardName,
			},
			SendTime: msgEvent.Time,
		},
	}
	buildEvent.Call()
}

//群消息撤回事件
func OnGroupMessageRecallEvent(qqClient *client.QQClient, msgEvent *client.GroupMessageRecalledEvent) {
	defer SetRecover()

	buildEvent := event.GroupMessageRecallEvent{
		GroupMessageRecallEventHandle: api.GroupMessageRecallEventHandle{
			GroupId:   msgEvent.GroupCode,
			MsgId:     msgEvent.MessageId,
			MsgUin:    msgEvent.AuthorUin,
			ActionUin: msgEvent.OperatorUin,
			Time:      msgEvent.Time,
		},
	}
	buildEvent.Call()
}

//群成员加入事件
func OnGroupUserJoinEvent(qqClient *client.QQClient, msgEvent *client.MemberJoinGroupEvent) {
	defer SetRecover()

	buildEvent := event.GroupUserJoinEvent{
		GroupUserJoinEventHandle: api.GroupUserJoinEventHandle{
			Group: api.Group{
				Id:   msgEvent.Group.Code,
				Name: msgEvent.Group.Name,
			},
			QQ: api.GroupUser{
				QQ: api.QQ{
					Uin:  msgEvent.Member.Uin,
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
func OnGroupUserQuitEvent(qqClient *client.QQClient, msgEvent *client.MemberLeaveGroupEvent) {
	defer SetRecover()
	buildEvent := event.GroupUserQuitEvent{
		GroupUserQuitEventHandle: api.GroupUserQuitEventHandle{
			Group: api.Group{
				Id:   msgEvent.Group.Code,
				Name: msgEvent.Group.Name,
			},
			QQ: api.GroupUser{
				QQ: api.QQ{
					Uin:  msgEvent.Member.Uin,
					Name: msgEvent.Member.Nickname,
				},
				CardName: msgEvent.Member.CardName,
			},
			ActionQQ: api.GroupUser{
				QQ: api.QQ{
					Uin:  msgEvent.Operator.Uin,
					Name: msgEvent.Operator.Nickname,
				},
				CardName: msgEvent.Operator.CardName,
			},
			Time: int32(time.Now().Unix()),
		},
	}
	buildEvent.Call()
}

func RunCacheMsg(elements []message.IMessageElement) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("缓存数据失败: %s", err)
		}
	}()

	imageSave := func(imageId string) {
		//保存
		//if !util.IsDir(app.Images){
		//	os.MkdirAll(app.Images,0644)
		//}
		//name := fmt.Sprintf("%s%s.img",app.Images,util.Md5String(imageId))
		//_ = ioutil.WriteFile(name,[]byte(imageId),0644)
		//fmt.Printf("缓存")
	}
	audioSave := func(imageId string, url string) {
		if !util.IsDir(app.Audios) {
			os.MkdirAll(app.Audios, 0644)
		}
		fmt.Printf(imageId)
		name := fmt.Sprintf("%s%s.amr", app.Audios, util.Md5String(imageId))
		res, err := req.Get(url)
		if err == nil {
			if b, err := res.ToBytes(); err == nil {
				_ = ioutil.WriteFile(name, b, 0644)
			}
		}

	}

	for _, item := range elements {
		switch e := item.(type) {
		case *message.GroupImageElement: //群组图片
			imageSave(e.ImageId)
			break
		case *message.ImageElement:
			imageSave(e.Filename)
			break
		case *message.VoiceElement:
			audioSave(e.Name, e.Url)
			break
		}
	}

}
