package entity

//实体

//登录账户
type UserConfig struct {
	QQ         int64       `json:"qq"`         //QQ号码
	Password   string      `json:"password"`   //密码
	AllowGroup []int64     `json:"allowGroup"` //允许接收群消息的群组
	Redis      RedisConfig `json:"redis"`      //redis配置
}
type RedisConfig struct {
	Host     string `json:"host"`     //redis地址
	Port     uint16 `json:"port"`     //redis端口
	Password string `json:"password"` //redis密码
	Index    int    `json:"index"`    //数据库选择
}
