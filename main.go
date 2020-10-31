package main

import (
	"bufio"
	"fmt"
	"github.com/balrogsxt/xtbot-go/app"
	_ "github.com/balrogsxt/xtbot-go/util/logger"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("运行发生异常: 【%s】请输入任意字符退出...\n", err)
			bufio.NewScanner(os.Stdin).Scan()
		}
	}()
	app.AppLinkStart()
}
