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
	"strings"
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
func NewFace(id int32) IMsg {
	return Face{
		Id: id,
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

//发送指定ID图片
func NewImageId(imageId string) IMsg {
	return Image{
		Id: imageId,
	}
}

func NewAudio(groupId int64, id_file_url string) IMsg {
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
	CQ码消息拆分实现,当前已支持的格式,看着好sb

	--------已实现的功能

	[CQ:at,qq=2289453456] = at指定 0代表全体群员
	[CQ:image,file=xxx.png] 发送图片,支持本地、网络、Id发送
	[CQ:face,id=123] QQ表情
	[CQ:record,file=xxxx.amr] 发送语音文件,如果是别人发的首先本地需要有缓存

	--------还未实现的功能,后面的功能看情况实现吧...
	[CQ:shake]  抖动
	[CQ:dice]   骰子
	[CQ:poke]   戳一戳
	[CQ:share,url=http://xxxxx] 分享链接
	[CQ:contact,type=qq,id=100000] 推荐好友
	[CQ:contact,type=group,id=100000] 推荐群
	[CQ:location,lat=39,lon=39]  位置
	[CQ:music,type=163,id=10000] 音乐分享
	[CQ:reply,id=1000]  回复消息
	[CQ:forward,id=1000] 合并转发
	[CQ:node,id=123456]  合并转发节点
	[CQ:xml,data=<?xml ...]  xml消息
	[CQ:json,data={"app": ...] json消息
	*/
	reg, _ := regexp.Compile("\\[CQ:([a-z]+).*?\\]|\\D|\\d")
	list := reg.FindAllString(text, -1)

	els := make([]IMsg, 0)

	tmpText := ""
	for _, txt := range list {
		regex, _ := regexp.Compile("\\[CQ:([a-z]+)(.*)?\\]")
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
			//解析参数值
			_tmpValue := strings.Split(item[2], ",")
			_value := make(map[string]string, 0)
			for _, p := range _tmpValue {
				ks := strings.Split(p, "=")
				if len(ks) >= 2 {
					_value[ks[0]] = ks[1]
				}
			}

			switch _type {
			case "at":
				//At用户 可能存在的参数:qq
				elem = func(v map[string]string) IMsg {
					if qq, has := v["qq"]; has {
						if qq, e := strconv.ParseInt(qq, 10, 64); e == nil {
							return NewAt(qq)
						}
					}
					return nil
				}(_value)
				break
			case "image":
				//发送图片 可能存在的参数:file 支持id、文件路径、网络url地址
				elem = func(v map[string]string) IMsg {
					if file, has := v["file"]; has {
						return NewImage(groupId, file)
					}
					return nil
				}(_value)
				break

			case "audio":
				//发送语音 如果参数是字符串id则必须在本地接收过该语音文件缓存,否则无法发送
				//可能存在的参数:file 支持id、文件路径、网络url地址
				elem = func(v map[string]string) IMsg {
					if file, has := v["file"]; has {
						return NewAudio(groupId, file)
					}
					return nil
				}(_value)
				break

			case "face": //发送表情 可能存在的参数:id
				elem = func(v map[string]string) IMsg {
					if id, has := v["id"]; has {
						if _id, err := strconv.ParseInt(id, 10, 32); err == nil {
							return NewFace(int32(_id))
						}
					}
					return nil
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
