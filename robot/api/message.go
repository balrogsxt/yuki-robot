package api

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/app/db"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/imroc/req"
	"io/ioutil"
	"regexp"
	"strconv"
)

var qqClient *client.QQClient = nil

func SetLoginQQClient(client *client.QQClient) {
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
		QQ:      0,
		Display: "@全体成员",
	}
}

//At用户
func NewAt(qq int64, display ...string) IMsg {
	show := fmt.Sprintf("@%d", qq)
	if qq == 0 {
		show = "@全体成员"
	}
	if len(display) >= 1 {
		show = fmt.Sprintf("@%s", display[0])
	}

	return At{
		QQ:      qq,
		Display: show,
	}
}
func AtCode(qq int64) string {
	return fmt.Sprintf("[type=at,value=%d]", qq)
}

//发送指定ID图片
func NewImageId(imageId string) IMsg {
	return Image{
		Id: imageId,
	}
}
func ImageCode(id_url_file string) string {
	return fmt.Sprintf("[type=image,value=%s]", id_url_file)
}
func NewAudio(groupId int64, id_file_url string) IMsg {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("异常: ", err)
		}
	}()
	var _fileByte []byte = nil
	var er error

	flag, _ := regexp.Match("^\\{[a-zA-Z0-9-]+\\}\\.(amr)$", []byte(id_file_url))
	flag2, _ := regexp.Match("^http(s?)", []byte(id_file_url))
	if flag {
		//读取本地缓存语音
		cacheFile := fmt.Sprintf("%s%s.amr", app.Audios, util.Md5String(id_file_url))
		if util.IsFile(cacheFile) {
			_fileByte, er = ioutil.ReadFile(cacheFile)
			if er != nil {
				logger.Warning("[群组语音] 读取本地缓存失败: %s", er.Error())
				return nil
			}
		} else {
			logger.Warning("[群组语音] 缓存文件不存在: %s", cacheFile)
			return nil
		}
	} else if flag2 {
		res, err := req.Get(id_file_url)
		if err != nil {
			logger.Warning("[群组语音] 远程请求失败: %s", err.Error())
			return nil
		}
		_fileByte, er = res.ToBytes()
		if er != nil {
			logger.Warning("[群组语音] 读取失败: %s", er.Error())
			return nil
		}
	} else {
		if util.IsFile(id_file_url) {
			_fileByte, er = ioutil.ReadFile(id_file_url)
			if er != nil {
				logger.Warning("[群组语音] 读取失败: %s", er.Error())
				return nil
			}
		} else {
			logger.Warning("[群组语音] 找不到语音文件: %s", id_file_url)
		}
	}

	if _fileByte == nil {
		logger.Warning("[群组图片] 获取图片数据失败")
		return nil
	}

	a, err := qqClient.UploadGroupPtt(groupId, _fileByte)
	if err != nil {
		logger.Warning("[群组语音] 上传失败: %s", err.Error())
		return nil
	}
	return Audio{
		Data: a.Data,
		Ptt:  a.Ptt,
	}
}

//发送本地图片或指定ID图片
func NewImage(groupId int64, id_file_url string) IMsg {
	flag, _ := regexp.Match("^\\{[a-zA-Z0-9-]+\\}\\.(PNG|JPG|JPEG|GIF|BMP|WEBP)$", []byte(id_file_url))
	if flag {
		return Image{
			Id: id_file_url,
		}
	}

	flag, _ = regexp.Match("^http(s?)", []byte(id_file_url))
	var _fileByte []byte = nil
	var er error
	if flag {
		res, err := req.Get(id_file_url)
		if err != nil {
			logger.Warning("[群组图片] 远程请求失败: %s", err.Error())
			return nil
		}
		_fileByte, er = res.ToBytes()
		if er != nil {
			logger.Warning("[群组图片] 读取失败: %s", er.Error())
			return nil
		}

	} else {
		_fileByte, er = ioutil.ReadFile(id_file_url)
		if er != nil {
			logger.Warning("[群组图片] 读取失败: %s", er.Error())
			return nil
		}
	}
	if _fileByte == nil {
		logger.Warning("[群组图片] 获取图片数据失败")
		return nil
	}

	img, err := qqClient.UploadGroupImage(groupId, _fileByte)
	if err != nil {
		logger.Warning("[群组图片] 上传失败: %s", err.Error())
		return nil
	}
	return Image{
		Id: img.ImageId,
	}
}

///////////主动API
//发送群消息[自定义CQ码方式]
func SendGroupMessageText(groupId int64, text string) GroupMsgId {
	/**
	普通文字普通文字普通文字
	[type=at,value=2289453456] = at指定
	[type=at,value=0] = at全体
	[type=face,id=100]
	[type=image,value={E4AD8A49-C2E8-1287-E5B3-559F7E5376AF}.PNG] = 发送图片ID
	[type=image,value=./test/1.jpg] = 发送本地图片
	*/
	reg, _ := regexp.Compile("\\[type=.*?,value=.*?\\]|\\D|\\d")
	list := reg.FindAllString(text, -1)

	els := make([]IMsg, 0)

	tmpText := ""
	for _, txt := range list {
		regex, _ := regexp.Compile("\\[type=(.*)?,value=(.*)?\\]")
		item := regex.FindStringSubmatch(txt)
		if len(item) >= 3 {
			//如果匹配成功到了,结算文字
			if len(tmpText) > 0 {
				els = append(els, NewText(CQCodeUnescapeText(tmpText)))
				tmpText = ""
			}

			var elem IMsg

			//判断类型
			_type := item[1]
			_value := item[2]

			switch _type {
			case "at":
				elem = func(v string) IMsg {
					if qq, e := strconv.ParseInt(v, 10, 64); e == nil {
						return NewAt(qq)
					} else {
						return nil
					}
				}(_value)
				break
			case "image":
				elem = func(v string) IMsg {
					return NewImage(groupId, v)
				}(_value)
				break
			case "audio":
				elem = func(v string) IMsg {
					return NewAudio(groupId, v)
				}(_value)
				break
			}
			if elem != nil {
				els = append(els, elem)
			}
		} else {
			tmpText += txt
		}
	}
	//再次判断是否需要结算文字
	if len(tmpText) > 0 {
		els = append(els, NewText(CQCodeUnescapeText(tmpText)))
		tmpText = ""
	}
	return SendGroupMessage(groupId, els)
}

//发送群组消息[结构方式]
func SendGroupMessage(groupId int64, msg []IMsg) GroupMsgId {
	parseMsg := ParseToOldElement(msg)
	result := message.SendingMessage{Elements: parseMsg}

	m := qqClient.SendGroupMessage(groupId, &result)
	go saveGroupMsg(groupId, msg, m)
	return GroupMsgId{
		MsgId: MsgId{
			Id:         m.Id,
			InternalId: m.InternalId,
		},
		Group: Group{
			Id:   m.GroupCode,
			Name: m.GroupName,
		},
	}
}

//撤回群组消息
func RecallGroupMessage(groupId int64, msgId int32) {
	//查询数据库获取内部id
	msg := new(db.GroupMsg)
	has, err := app.GetDb().Where("msg_id = ?", msgId).Cols("internal_id").Get(msg)
	if err != nil {
		logger.Warning("[撤回消息] 获取内部消息Id失败: %s", err.Error())
	}
	if has {
		qqClient.RecallGroupMessage(groupId, msgId, msg.InternalId)
	}
}

//保存群组发送的消息
func saveGroupMsg(groupId int64, list []IMsg, m *message.GroupMessage) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("[群组消息保存] 保存消息发生错误: %s", err)
		}
	}()
	msg := new(db.GroupMsg)
	msg.MsgId = m.Id
	msg.InternalId = m.InternalId
	msg.Group = groupId
	msg.QQ = qqClient.Uin //当前登录qq
	msg.SendTime = m.Time
	msg.MsgText = ToString(list)
	msg.MsgJson = ToJson(list)

	if _, err := app.GetDb().InsertOne(msg); err != nil {
		logger.Warning("[群组消息保存] 保存消息发生错误: %s", err)
	}
}
