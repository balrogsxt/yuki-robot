package group

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/app/db"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/robot/cq"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"strings"
)

//删除关键字回复
type DelKey struct {
}

func (this DelKey) Command() string {
	return "^#del"
}
func (this *DelKey) Call(value string, event api.GroupMessageEventHandle) bool {
	//空格分离
	s := strings.Split(value, " ")
	key := strings.Trim(value[0:len(s[0])], " ")
	if len(s) >= 2 {
		//指定关键词回复删除
		reply := strings.Trim(value[strings.Index(value, s[1]):], " ")
		if len(key) == 0 || len(reply) == 0 {
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n需要删除的关键词回复不能为空!")
			return true
		}
		md5 := util.Md5String(fmt.Sprintf("%s_%s", key, reply))
		size, err := app.GetDb().Where("md5 = ?", md5).Delete(&db.GroupReply{})
		if err != nil {
			logger.Warning("[删除回复] 删除回复指定回复词条失败: %s", err.Error())
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n删除指定回复词条失败!")
			return true
		}
		if size == 0 {
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n没有可删除的词条指定回复!")
			return true
		}
		event.Group.SendGroupMessageText(cq.At(event.QQ.Uin) + "\n删除指定词条回复成功!")
	} else {
		if len(key) == 0 {
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n需要删除的关键词回复不能为空!")
			return true
		}
		//删除指定词条的全部回复
		size, err := app.GetDb().Where("`key` = ?", key).Delete(&db.GroupReply{})
		if err != nil {
			logger.Warning("[删除回复] 删除关键词回复失败: %s", err.Error())
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n删除关键词回复失败!")
			return true
		}
		if size == 0 {
			api.NewGroupException(cq.At(event.QQ.Uin) + "\n这个关键词当前暂无其他回复词条!")
			return true
		}
		event.Group.SendGroupMessageText(fmt.Sprintf("%s\n删除关键词回复成功,累计删除: %#v 个回复", cq.At(event.QQ.Uin), size))
	}
	return true

}
