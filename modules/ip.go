package modules

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/event"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/msg"
	"github.com/imroc/req"
)

type Ip struct {
	event.GroupModuleHandle //继承
}

//模块触发命令
func (this *Ip) GetCommand() string {
	return "#ip"
}

//模块名称
func (this *Ip) GetName() string {
	return "IP查询"
}

//模块触发调用
func (this *Ip) Handle() {
	if len(this.Args) == 0 || !util.IsIpv4(this.Args[0]) {
		this.SendMessage("%s阁下输入的格式不正确!", msg.At(this.Event.Sender.Uin))
		return
	}
	//调用七空幻音3.0API获取（开发中）
	res, err := req.Get(fmt.Sprintf("http://127.0.0.1:10260/v1/http/ip?ip=%s", this.Args[0]))
	if err != nil {
		this.SendMessage("%s获取IP归属地失败", msg.At(this.Event.Sender.Uin))
		return
	}
	// todo 之后需要整合一下返回结构体,现在暂时用一下匿名结构
	data := &struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			Ip   string `json:"ip"`
			Addr string `json:"addr"`
		} `json:"data"`
	}{}

	if err := res.ToJSON(data); err != nil {
		this.SendMessage("%s解析IP归属地失败", msg.At(this.Event.Sender.Uin))
		return
	}
	if data.Status != 0 {
		this.SendMessage("%s查询失败: %s", msg.At(this.Event.Sender.Uin), data.Msg)
		return
	}
	this.SendLineString([]string{
		msg.At(this.Event.Sender.Uin) + "查询成功",
		fmt.Sprintf("IP地址: %s", data.Data.Ip),
		fmt.Sprintf("归属地: %s", data.Data.Addr),
	})
}
