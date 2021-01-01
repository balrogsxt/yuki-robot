package api

//机器人任务计划接口
type Task interface {
	GetCron() string //任务计划触发时间格式
	Call() //任务计划触发
}