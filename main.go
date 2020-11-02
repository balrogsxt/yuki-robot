package main

import (
	"bufio"
	"github.com/balrogsxt/xtbot-go/app"
	"github.com/balrogsxt/xtbot-go/modules"
	"github.com/balrogsxt/xtbot-go/util/logger"
	_ "github.com/balrogsxt/xtbot-go/util/logger"
	"os"
)

//使用方法，直接调用即可输出带颜色的文本
func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatal("运行发生异常: 【%s】请输入任意字符退出...\n", err)
			bufio.NewScanner(os.Stdin).Scan()
		}
	}()

	app.AppLinkStart(func(bot *app.QQBot) {
		//注册自定义群聊模块
		bot.RegisterGroupModule(new(modules.Ip))
	})
}
