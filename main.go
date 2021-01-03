package main

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/robot"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
	_ "github.com/balrogsxt/xtbot-go/util/logger"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatal("运行发生异常: %#v", err)
		}
	}()
	logger.Info("正在尝试启动机器人...")
	config := app.GetRobotConfig()

	//设定缓存模块
	cache := new(app.RedisCache)
	if err := cache.Init(config); err != nil {
		panic(fmt.Sprintf("初始化缓存模块失败: %s", err.Error()))
	} else {
		logger.Info("[缓存模块] 初始化成功")
	}
	api.InitCache(cache)

	if err := app.InitDatabase(config); err != nil {
		panic(fmt.Sprintf("初始化Mysql失败: %s", err.Error()))
	}

	_robot, err := robot.NewRobot()
	if err != nil {
		panic(err.Error())
	}

	_robot.OnRobotStart(func(r *robot.Robot) {
		//绑定任务计划
		//r.AddTask(new(task.TestTask))
	})
	_robot.Run()

}
