package group

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/app/db"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"strings"
	"time"
)

//回复关键词添加
type AddKey struct {
}

func (AddKey) Command() string {
	return "^#add"
}
func (AddKey) Call(value string, event api.GroupMessageEventHandle) bool {

	//空格分离
	s := strings.Split(value, " ")
	if len(s) < 2 {
		event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "\n自定义回复格式不正确~")
		return true
	}
	key := strings.Trim(value[0:len(s[0])], " ")
	reply := strings.Trim(value[strings.Index(value, s[1]):], " ")
	if len(key) == 0 || len(reply) == 0 {
		event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "\n关键词回复不能为空!")
		return true
	}
	md5 := util.Md5String(fmt.Sprintf("%s_%s", key, reply))

	has, _ := app.GetDb().Where("md5 = ?", md5).Exist(&db.GroupReply{})
	if has {
		event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "\n词条回复已存在")
		return true
	}

	//添加到数据库
	r := new(db.GroupReply)
	r.Key = key
	r.Reply = reply
	r.QQ = event.QQ.Uin
	r.Group = event.Group.Id
	r.Time = int32(time.Now().Unix())
	r.Global = 1
	r.Md5 = md5

	if _, err := app.GetDb().InsertOne(r); err != nil {
		logger.Warning("[添加回复] 添加自定义回复失败: %s", err.Error())
		event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "\n添加自定义回复失败: 数据库错误~")
		return true
	}

	event.Group.SendGroupMessageText(api.AtCode(event.QQ.Uin) + "\n添加自定义回复成功!")
	return true
}
