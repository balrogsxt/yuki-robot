package db

//数据库实体

type GroupMsg struct {
	Id         int
	MsgId      int32  `xorm:"msg_id"`      //消息id
	InternalId int32  `xorm:"internal_id"` //内部id
	Group      int64  //发送群
	QQ         int64  `xorm:"qq"`        //发送者
	SendTime   int32  `xorm:"send_time"` //发送时间
	MsgText    string `xorm:"msg_text"`  //字符串格式的消息
	MsgJson    string `xorm:"msg_json"`  //json字符串格式的消息
}
