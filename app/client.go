package app

import (
	"errors"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/balrogsxt/xtbot-go/event"
	"github.com/balrogsxt/xtbot-go/util"
	"github.com/balrogsxt/xtbot-go/util/cache"
	"github.com/balrogsxt/xtbot-go/util/entity"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/balrogsxt/xtbot-go/util/msg"
	"io/ioutil"
	"os"
)

//机器人结构
type QQBot struct {
	Handle *client.QQClient  //核心QQ客户端协议模块
	config entity.UserConfig //登录的用户数据
	Cache  cache.XtCache     //缓存模块
}

//启动
func AppLinkStart() {
	config, err := ParseUserConfig()
	if err != nil {
		ThrowException("登录配置文件处理失败:%s", err.Error())
	}
	bot := new(QQBot)
	bot.Login(config)
}

//解析需要登录的机器人账户
func ParseUserConfig() (entity.UserConfig, error) {
	var la entity.UserConfig
	file, err := os.Open("./config.json")
	if err != nil {
		return la, errors.New("打开配置文件失败")
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return la, errors.New("读取配置文件失败")
	}
	if err := util.JsonDecode(string(b), &la); err != nil {
		return la, errors.New("解析配置文件失败" + err.Error())
	}
	return la, nil
}
func (this *QQBot) Login(config entity.UserConfig) {
	logger.Info("正在准备尝试登录QQ:[%d]...", config.QQ)
	//创建一个新的QQ客户端
	this.Handle = client.NewClient(config.QQ, config.Password)
	//请求登录
	res, err := this.Handle.Login()
	if err != nil {
		ThrowException("初始化登录客户端失败:%s", err.Error())
	}
	//判断登录是否需要验证处理
	if !res.Success {
		switch res.Error {
		case client.OtherLoginError:
			ThrowException("登录错误:%s【请检查账户密码是否正确】", res.ErrorMessage)
			break
		case client.SMSNeededError, client.NeedCaptcha, client.SMSOrVerifyNeededError:
			fmt.Printf("\n%s\n\n", res.VerifyUrl)
			ThrowException("请在浏览器打开验证链接,处理完成后重新启动: %s", res.ErrorMessage)
			break
		default:
			ThrowException("未处理的异常: %s", res.ErrorMessage)
			break
		}
		return
	}
	this.config = config
	logger.Info("QQ:%d 已经登录成功", config.QQ)
	logger.Info("已允许接收群组消息列表: %#v", config.AllowGroup)
	//初始化缓存模块
	this.registerCache()
	//开始监听各项数据
	this.registerEvent() //调用注册事件
	//开启命令行输入进程顺便阻止退出
	StartCommand()
}

//注册缓存模块
func (this *QQBot) registerCache() {
	//默认使用redis缓存
	this.Cache = new(cache.RedisCache)
	err := this.Cache.Init(this.config)
	if err != nil { //缓存模块初始化失败
		ThrowException(err.Error())
	}
	logger.Info("缓存模块初始化成功")
}

//注册QQ事件
func (this *QQBot) registerEvent() {
	//注册群聊消息事件
	this.Handle.OnGroupMessage(func(qqClient *client.QQClient, ev *message.GroupMessage) {
		isAllow := false
		for _, item := range this.config.AllowGroup {
			if item == ev.GroupCode {
				isAllow = true
				break
			}
		}
		if !isAllow {
			return
		}
		this.saveGroupQQ(ev.GroupCode, ev.Sender)
		logger.Info("[群聊消息-> %d -> %s] %s", ev.GroupCode, ev.GroupName, ev.ToString())
		handle := &msg.GroupHandle{
			Handle: qqClient,
			Event:  ev,
			MsgBuild: &msg.GroupMessageBuilder{
				MessageBuilder: msg.MessageBuilder{
					Handle: qqClient,
					Event:  ev,
					Cache:  this.Cache,
				},
			},
		}
		event.OnGroupMessageEvent(handle)
	})
	//注册群聊消息撤回事件
	this.Handle.OnGroupMessageRecalled(func(qqClient *client.QQClient, msg *client.GroupMessageRecalledEvent) {
		logger.Info("[群聊撤回 -> %d] %d", msg.GroupCode, msg.MessageId)
		event.OnGroupMessageRecallEvent(qqClient, msg)
	})
	//注册私聊消息事件
	this.Handle.OnPrivateMessage(func(qqClient *client.QQClient, msg *message.PrivateMessage) {
		logger.Info("[私聊消息 -> %d] %s", msg.Sender.Uin, msg.ToString())
		event.OnPrivateMessageEvent(qqClient, msg)
	})
	//断开连接
	this.Handle.OnDisconnected(func(qqClient *client.QQClient, disconnectedEvent *client.ClientDisconnectedEvent) {
		logger.Info("[断开连接] %s", disconnectedEvent.Message)
	})
	//更多事件,按需求写...
}

//缓存群内成员数据
func (this *QQBot) saveGroupQQ(groupId int64, sender *message.Sender) {
	// todo 暂时收到群成员信息就写入缓存吧...
	key := fmt.Sprintf("cache:group:qq:%d", groupId)
	field := fmt.Sprintf("%d", sender.Uin)
	//if flag := this.Cache.ExistsMap(key,field); !flag{
	//写入缓存
	this.Cache.SetMap(key, field, util.JsonEncode(&sender))
	//}else{
	//存在缓存,这个需要永久缓存,但是需要定期更新数据

	//}
}
