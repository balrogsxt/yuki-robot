package msg

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/cache"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"io/ioutil"
)

//公共消息基类
type MessageBuilder struct {
	Handle *client.QQClient
	Event  *message.GroupMessage
	Cache  cache.XtCache
}

//单一文本消息
func (this *MessageBuilder) Text(format string, args ...interface{}) *message.TextElement {
	return &message.TextElement{
		Content: fmt.Sprintf(format, args...),
	}
}

//通过图片ID进行发送
func (this *MessageBuilder) ImageId(imgId string) *message.ImageElement {
	return &message.ImageElement{
		Filename: imgId,
	}
}

//群组消息
type GroupMessageBuilder struct {
	MessageBuilder
}

//at消息,群里才有at功能
func (this *GroupMessageBuilder) At(qq int64) *message.AtElement {
	//这里如果不给用户名称的话,手机端会显示为空,不会自动匹配昵称
	key := fmt.Sprintf("cache:group:qq:%d", this.Event.GroupCode)
	field := fmt.Sprintf("%d", qq)

	nickName := ""
	res, err := this.Cache.GetMap(key, field)
	if err == nil {
		sender := new(message.Sender)
		err := util.JsonDecode(res, sender)
		if err == nil {
			nickName = fmt.Sprintf("@%s", sender.Nickname)
		}
	}
	if len(nickName) == 0 {
		//如果没有咋办呢?
		nickName = fmt.Sprintf("@%d", qq)
		// todo 可以使用api远程获取qq昵称
	}

	return &message.AtElement{
		Target:  qq,
		Display: nickName,
	}
}

//at全体
func (this *GroupMessageBuilder) AtAll() *message.AtElement {
	return this.At(0)
}

//上传&发送本地资源图片
func (this *GroupMessageBuilder) LocalImage(fileName string) *message.GroupImageElement {
	md5 := util.Md5File(fileName)
	cacheField := fmt.Sprintf("cache:group:image:%s", md5)
	if cache, err := this.Cache.Get(cacheField); err == nil {
		img := new(message.GroupImageElement)
		if err := util.JsonDecode(cache, &img); err == nil {
			return img
		}
	}
	//读取图片数据
	byte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Warning("读取图片失败: %s", err.Error())
		return nil
	}
	img, err := this.Handle.UploadGroupImage(this.Event.GroupCode, byte)
	if err != nil {
		logger.Warning("上传图片失败: %s", err.Error())
		return nil
	}
	this.Cache.Set(cacheField, util.JsonEncode(&img)) //永久缓存
	return img
}

//type PrivateMessageBuilder struct {
//	MessageBuilder
//}
//func (this *PrivateMessageBuilder) Image(fileName string) *message.FriendImageElement  {
//	md5 := util.Md5File(fileName)
//	cacheField := fmt.Sprintf("cache:private:image:%s",md5)
//	if cache,err := this.Cache.Get(cacheField);err == nil {
//		img := new(message.FriendImageElement)
//		if err := util.JsonDecode(cache,&img);err == nil {
//			return img
//		}
//	}
//	//读取图片数据
//	byte,err := ioutil.ReadFile(fileName)
//	if err != nil {
//		logger.Warning("读取图片失败: %s",err.Error())
//		return nil
//	}
//	img,err := this.Handle.UploadPrivateImage(this.Event.GroupCode,byte)
//	if err != nil {
//		logger.Warning("上传图片失败: %s",err.Error())
//		return nil
//	}
//	this.Cache.Set(cacheField,util.JsonEncode(&img)) //永久缓存
//	return img
//}
