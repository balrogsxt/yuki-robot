package app

import (
	"bufio"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"os"
)

func StartCommand() {
	logger.Info("现在可以输入指令来控制啦~")
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		text := input.Text()
		logger.Info("还不支持的命令: %s", text)
	}
}
