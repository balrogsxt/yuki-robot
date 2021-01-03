package robot

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/robfig/cron"
	"os"
	"time"
)

type Robot struct {
	Config *app.RobotConfig
	Handle *client.QQClient

	cronTask *cron.Cron //任务计划

	robotStartEvents func(*Robot) //机器人登录成功事件

}

//创建新的机器人
func NewRobot() (*Robot, error) {
	robot := new(Robot)
	conf, err := app.LoadRobotConfig()
	if err != nil {
		return nil, err
	}
	robot.Config = conf

	return robot, nil
}
func (this *Robot) Run() {
	if err := this.login(); err != nil {
		//首次登录失败,则直接异常结束
		logger.Fatal("[登录失败] 登录QQ账户失败: %s", err.Error())
	}
	this.robotCommand()

}

//登录机器人账户
func (this *Robot) login() error {
	config := this.Config
	pwd := config.User.Password
	_client := client.NewClient(config.User.QQ, pwd)
	res, err := _client.Login()
	if err != nil {
		return errors.New(fmt.Sprintf("初始化登录客户端失败:%s", err.Error()))
	}
	if !res.Success {
		switch res.Error {
		case client.OtherLoginError:
			return errors.New(fmt.Sprintf("登录错误:%s【请检查账户密码是否正确】", res.ErrorMessage))
		case client.SMSNeededError, client.NeedCaptcha, client.SMSOrVerifyNeededError:
			fmt.Printf("\n%s\n\n", res.VerifyUrl)
			return errors.New(fmt.Sprintf("请在浏览器打开验证链接,处理完成后重新启动: %s", res.ErrorMessage))
		default:
			return errors.New(fmt.Sprintf("未处理的异常: %s", res.ErrorMessage))
		}
	}
	this.Handle = _client

	_client.OnDisconnected(func(qqClient *client.QQClient, event *client.ClientDisconnectedEvent) {
		api.SetLoginQQClient(nil) //设定登录失效
		this.cronTask.Stop()      //停止任务计划
		logger.Error("[账户离线] %s", event.Message)
		//重新连接
		for {
			if err := this.login(); err != nil {
				logger.Fatal("[账户重连] 重新连接QQ失败: %s -> 预计60秒后尝试重新连接", err.Error())
				//重连失败,60秒后重试
				time.Sleep(time.Second * 60)
			} else {
				logger.Info("[重新登录] %s(%d) 已重新登录成功!", _client.Nickname, _client.Uin)
				break
			}
		}
	})

	api.SetLoginQQClient(_client)

	logger.Info("[登录成功] %s(%d) 已登录成功!", _client.Nickname, _client.Uin)
	this.registerEvent() //注册客户端相关事件

	return nil
}

//注册消息事件
func (this *Robot) registerEvent() {
	h := this.Handle
	//注册群组相关事件
	h.OnGroupMessage(OnGroupMessageEvent)               //群消息接收事件
	h.OnGroupMessageRecalled(OnGroupMessageRecallEvent) //群消息撤回事件
	h.OnGroupMemberJoined(OnGroupUserJoinEvent)
	h.OnGroupMemberLeaved(OnGroupUserQuitEvent)

	//初始化任务计划
	this.cronTask = cron.New()
	//触发启动事件
	if this.robotStartEvents != nil {
		this.robotStartEvents(this)
	}
	this.cronTask.Start()
}

//机器人启动成功后触发
func (this *Robot) OnRobotStart(e func(*Robot)) {
	this.robotStartEvents = e
}

//添加任务计划
func (this *Robot) AddTask(task api.Task) {
	if err := this.cronTask.AddFunc(task.GetCron(), task.Call); err != nil {
		logger.Warning("[任务计划] 添加失败: [%s] -> %s", task.GetCron(), err.Error())
	} else {
		logger.Info("[任务计划] 添加成功: [%s]", task.GetCron())
	}
}

func (this *Robot) robotCommand() {
	logger.Info("[命令模式] 命令模式已启动")
	fmt.Println()
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		text := input.Text()

		switch text {
		case "stop":
			logger.Info("正在退出...")
			this.Handle.Disconnect()
			logger.Info("已退出QQ")
			return
		default:
			logger.Info("还不支持的命令: %s", text)
		}
	}
}
