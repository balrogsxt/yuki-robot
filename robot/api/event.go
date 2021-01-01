package api
//群消息接收事件数据
type GroupMessageEventHandle struct {
	MsgId GroupMsgId //群组消息ID
	Group Group //消息发送的群数据
	QQ GroupUser //消息发送的群成员数据
	MessageList []IMsg //消息结构列表
	SendTime int32 //消息发送时间
}
//消息撤回事件数据
type GroupMessageRecallEventHandle struct {
	GroupId int64 //撤回的所属群组Id
	MsgId int32 //撤回的消息ID
	MsgUin int64 //消息归属QQ号
	ActionUin int64 //撤回操作QQ号
	Time int32 //撤回时间
}
//群成员加入事件数据
type GroupUserJoinEventHandle struct {
	Group Group //加入的群组
	QQ GroupUser //加入的用户
	Time int32 //时间
}
//群成员退出事件数据
type GroupUserQuitEventHandle struct {
	Group Group //退出的群组信息
	QQ GroupUser //退群的用户
	ActionQQ GroupUser //操作的用户
	Time int32 //时间
}