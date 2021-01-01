package app

import (
	"errors"
	"fmt"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//机器人的配置文件处理
type RobotConfig struct {
	//QQ账户用户信息
	User struct{
		QQ int64
		Password string
	}
	Group struct{ //群组设置
		Allow []int64 //支持接收消息的群组
		Deny []int64 //拒绝接收消息的群组
	}
	Cache struct{ //缓存模块
		//支持的缓存模块配置
		//redis
		Redis struct{
			Host string
			Port int
			Password string
			Index int
		}


	}
}
//载入机器人配置
func LoadRobotConfig() (*RobotConfig,error)  {
	file := "./config.yml"
	_byte,err := ioutil.ReadFile(file)
	if err != nil {
		return nil,err
	}
	conf := RobotConfig{}
	if err := yaml.Unmarshal(_byte,&conf);err != nil {
		return nil,errors.New(fmt.Sprintf("解析配置文件失败: %s",err.Error()))
	}
	config = &conf
	return &conf,nil
}
var config *RobotConfig
func GetRobotConfig() *RobotConfig  {
	if config == nil {
		conf,err := LoadRobotConfig()
		if err != nil {
			logger.Fatal("[配置文件] 加载失败: %s",err.Error())
			return nil
		}
		return conf
	}else{
		return config
	}
}