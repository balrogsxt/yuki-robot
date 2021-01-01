package task

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/robot/api"
)

type TestTask struct {
	api.Task
}

func (receiver TestTask) GetCron() string {
	return "0 */1 * * * *" //每1分钟触发一次
}
func (receiver TestTask) Call()  {
	fmt.Println("命令触发成功")
}