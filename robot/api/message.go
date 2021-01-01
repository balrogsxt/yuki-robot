package api

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/imroc/req"
	"io/ioutil"
	"regexp"
)

var qqClient *client.QQClient = nil

func SetLoginQQClient(client *client.QQClient)  {
	qqClient = client
}

//创建普通文字消息
func NewText(content string) IMsg {
	return Text{
		Content: content,
	}
}
//At全体
func NewAtAll() IMsg {
	return At{
		QQ: 0,
		Display: "@全体成员",
	}
}
//At用户
func NewAt(qq int64,display ...string) IMsg {
	show := fmt.Sprintf("@%d",qq)
	if qq == 0{
		show = "@全体成员"
	}
	if len(display) >= 1 {
		show = fmt.Sprintf("@%s",display[0])
	}

	return At{
		QQ: qq,
		Display: show,
	}
}

//发送指定ID图片
func NewImageId(imageId string) IMsg {
	return Image{
		Id: imageId,
	}
}

//发送本地图片或指定ID图片
func NewImage(groupId int64,id_file_url string) IMsg {
	flag, _ := regexp.Match("^\\{[a-zA-Z0-9-]+\\}\\.(PNG|JPG|JPEG|GIF|BMP|WEBP)$",[]byte(id_file_url))
	if flag {
		return Image{
			Id: id_file_url,
		}
	}

	flag,_ = regexp.Match("^http(s?)",[]byte(id_file_url))
	var _fileByte []byte = nil
	var er error
	if flag {
		res,err := req.Get(id_file_url)
		if err != nil {
			logger.Warning("[群组图片] 远程请求失败: %s",err.Error())
		}
		_fileByte,er = res.ToBytes()
		if er != nil {
			logger.Warning("[群组图片] 读取失败: %s",er.Error())
			return nil
		}

	}else{
		_fileByte,er = ioutil.ReadFile(id_file_url)
		if er != nil {
			logger.Warning("[群组图片] 读取失败: %s",er.Error())
			return nil
		}
	}
	if _fileByte == nil {
		logger.Warning("[群组图片] 获取图片数据失败")
		return nil
	}
	img,err := qqClient.UploadGroupImage(groupId,_fileByte)
	if err != nil {
		logger.Warning("[群组图片] 上传失败: %s",err.Error())
		return nil
	}
	return Image{
		Id: img.ImageId,
	}
}

///////////主动API

//发送群组消息
func SendGroupMessage(groupId int64,msg []IMsg) GroupMsgId {

	parseMsg := ParseToOldElement(msg)
	result := message.SendingMessage{Elements: parseMsg}

	m := qqClient.SendGroupMessage(groupId,&result)
	return GroupMsgId{
		MsgId:MsgId{
			Id: m.Id,
			InternalId: m.InternalId,
		},
		Group:Group{
			Id:m.GroupCode,
			Name:m.GroupName,
		},
	}
}

//撤回群组消息
func RecallGroupMessage(groupId int64,msgId MsgId)  {
	qqClient.RecallGroupMessage(groupId,msgId.Id,msgId.InternalId)
}