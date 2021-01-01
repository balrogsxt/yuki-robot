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
)


type Robot struct {
	Config *app.RobotConfig
	Handle *client.QQClient

	cronTask *cron.Cron //任务计划

	robotStartEvents func(*Robot) //机器人登录成功事件

}


//创建新的机器人
func NewRobot() (*Robot,error) {
	robot := new(Robot)
	conf,err := app.LoadRobotConfig()
	if err != nil {
		return nil,err
	}
	robot.Config = conf

	return robot,nil
}

//登录机器人账户
func (this *Robot) Login() error {
	config := this.Config
	pwd := config.User.Password
	_client := client.NewClient(config.User.QQ,pwd)
	res,err := _client.Login()
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

	api.SetLoginQQClient(_client)

	logger.Info("%s(%d) 已登录成功!",_client.Nickname,_client.Uin)

	//加载群组
	groupList,err := _client.GetGroupList()
	if err == nil {
		for _,group := range groupList {
			logger.Info("[群组加载] %s(%d)",group.Name,group.Uin)
		}
	}else{
		logger.Warning("[群组加载失败] %s",err.Error())
	}

	this.registerEvent()

	this.robotCommand(this)
	return nil
}

//注册消息事件
func (this *Robot) registerEvent()  {
	h := this.Handle

	h.OnDisconnected(func(qqClient *client.QQClient, event *client.ClientDisconnectedEvent) {
		api.SetLoginQQClient(nil) //设定登录失效
		this.cronTask.Stop() //停止任务计划
		logger.Fatal("[账户离线] %s",event.Message)
	})
	//注册群组相关事件
	h.OnGroupMessage(OnGroupMessageEvent) //群消息接收事件
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
	StartHttpApi()
}
//机器人启动成功后触发
func (this *Robot) OnRobotStart(e func(*Robot)) {
	this.robotStartEvents = e
}

//添加任务计划
func (this *Robot) AddTask(task api.Task) {
	if err := this.cronTask.AddFunc(task.GetCron(),task.Call); err != nil {
		logger.Warning("[任务计划] 添加失败: [%s] -> %s",task.GetCron(),err.Error())
	}else{
		logger.Info("[任务计划] 添加成功: [%s]",task.GetCron())
	}
}
//添加群组模块
func (this *Robot) AddGroupModule(module api.GroupMessageModule)  {

}


func (this *Robot) robotCommand(robot *Robot) {
	logger.Info("现在可以输入指令来控制啦~")
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		text := input.Text()

		switch text {
		case "stop":
			logger.Info("正在退出...")
			robot.Handle.Disconnect()
			logger.Info("已退出QQ")
			return
		default:
			logger.Info("还不支持的命令: %s", text)
		}

	}
}
