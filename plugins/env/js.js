//发送群组消息结构列表,返回msgId
function sendGroupMessage(groupId,elementList){}
//发送群组消息自定义字符码,返回msgId
function sendGroupMessageText(groupId,text){}
//撤回群组消息
function recallGroupMessage(groupId,msgId){}
//用于延迟执行的函数,参数毫秒,替换setTimeout
function sleep(ms){}
//发送httpGet请求
function httpGet(url,header){}
//发送httpPost formdata请求
function httpPost(url,formdata,header){}
//发送httpPost json数据请求
function httpPostJson(url,json,header){}
//获取事件数据
function getEvent(){}
