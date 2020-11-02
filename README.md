# 简介
基于 [MiraiGo](https://github.com/Mrs4s/MiraiGo) 实现的协议开发的QQ机器人(可能不适用于所有人)
> 项目将核心协议库合并整合了一下方便开发

## 安装
克隆项目到本地

## 配置
> 创建config.json文件放置到根目录,配置如下(不含注释)
```
{
    "qq":0, //需要登录的QQ账户
    "password":"", //登录账户的QQ密码
    "allowGroup":[123,456], //允许接收群里消息的群组
    "redis": {                //redis缓存配置
        "host": "127.0.0.1",  //地址
        "port": 6379,         //端口
        "password": "",       //密码
        "index": 1            //数据库
    }
}
```
## 启动
```
//运行
go run main.go 
//编译
go build main.go
```
## 增加群组命令模块
在modules目录下创建结构体并继承`event.GroupModuleHandle`然后实现`GetCommand` `GetName` `Handle` 方法

## 机器人功能
> 以下功能均为本地未开放API接口
- [x] 查询IP归属地
- [ ] 查询手机号归属地
- [ ] 端口检测是否开启
- [ ] SSl证书信息查询
- [ ] 网站Favicon图片获取
- [ ] 二维码生成
- [ ] 更多功能想法....
