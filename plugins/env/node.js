const vm = require('vm');
try{
    //反馈一个特殊字符给服务端检测并触发
    var call = (fn,vars)=>{
        //后期做成http或socket之类的
        console.log(`[§fn_call§]${JSON.stringify({
            fn:fn,
            vars:vars
        })}`);
    }
    const vars = JSON.parse(process.argv[2])
    const code = process.argv[3]
    global.getEvent = ()=>{return vars}
    //发送群组消息
    global.sendGroupMessage = (groupId,text)=>{
        call("sendGroupMessage",{
            groupId:groupId,
            text:text
        })
    }
    vm.runInThisContext(code);
}catch(e){
    console.log("run error: ",e.message)
    throw new Error(e.message)
}

