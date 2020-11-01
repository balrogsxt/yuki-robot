# 简介
基于 [MiraiGo](https://github.com/Mrs4s/MiraiGo) 实现的协议开发的QQ机器人
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
## 调试&编译
```
//运行
go run main.go 
//编译
go build main.go
```
