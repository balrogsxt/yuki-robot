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

//at消息
func (this *MessageBuilder) At(qq int64) *message.AtElement {
	return message.NewAt(qq)
}

//at全体
func (this *MessageBuilder) AtAll() *message.AtElement {
	return this.At(0)
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
