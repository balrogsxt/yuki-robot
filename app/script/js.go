package script

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/balrogsxt/xtbot-go/util/logger"
	"github.com/imroc/req"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"time"
)

type Javascript struct {
	vm *otto.Otto
}

func NewJs() *Javascript {
	js := new(Javascript)

	//初始化绑定支持的内置方法
	js.bindInternalFunctions()
	return js
}
func (this *Javascript) bindInternalFunctions() {
	this.vm = otto.New()
	jsfn := NewJsFn(this.vm)

	//发送群消息
	this.vm.Set("sendGroupMessage", jsfn.fn_sendGroupMessage)
	//发送群消息,字符串形式
	this.vm.Set("sendGroupMessageText", jsfn.fn_sendGroupMessageText)
	//撤回群组消息
	this.vm.Set("recallGroupMessage", jsfn.fn_recallGroupMessage)
	//延迟函数,由于setTimeout无法再虚拟机中使用
	this.vm.Set("sleep", jsfn.fn_sleep)
	//提供http get方法
	this.vm.Set("httpGet", jsfn.fn_http_get)
	//提供http post formdata、json方法
	this.vm.Set("httpPost", jsfn.fn_http_post_formdata)
	this.vm.Set("httpPostJson", jsfn.fn_http_post_json)

}
func (this *Javascript) SetVars(name string, vars map[string]interface{}) {
	this.vm.Set(name, vars)
}
func (this *Javascript) RunFile(file string) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("[Js虚拟机] 运行Js脚本失败: %s", err)
		}
	}()
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	code := string(b)
	if _, err := this.vm.Run(code); err != nil {
		return err
	}
	return nil
}

type JsFn struct {
	vm *otto.Otto
}

func NewJsFn(vm *otto.Otto) *JsFn {
	v := new(JsFn)
	v.vm = vm
	return v
}

//////////////系统部分API

//延迟执行函数
func (this *JsFn) fn_sleep(call otto.FunctionCall) otto.Value {
	ms, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	time.Sleep(time.Millisecond * time.Duration(ms))
	return otto.Value{}
}

//http请求
func (this *JsFn) fn_http_get(call otto.FunctionCall) otto.Value {
	url := call.Argument(0).String()
	header, err := call.Argument(1).Export()
	_h := make(req.Header, 0)
	if err == nil {
		o, flag := header.(map[string]interface{})
		if flag {
			for k, v := range o {
				_h[k] = fmt.Sprintf("%#v", v)
			}
		}
	}
	res, err := req.Get(url, _h)
	if err != nil {
		return otto.Value{}
	}
	str, err := res.ToString()
	if err != nil {
		return otto.Value{}
	}
	result, _ := this.vm.ToValue(str)
	return result
}

//http post formdata
func (this *JsFn) fn_http_post_formdata(call otto.FunctionCall) otto.Value {
	return this.httpPost("formdata", call)
}

//http post json
func (this *JsFn) fn_http_post_json(call otto.FunctionCall) otto.Value {
	return this.httpPost("json", call)
}
func (this *JsFn) httpPost(t string, call otto.FunctionCall) otto.Value {
	url := call.Argument(0).String()
	_h := make(req.Header, 0)
	//header set
	header, err := call.Argument(2).Export()
	if err == nil {
		o, flag := header.(map[string]interface{})
		if flag {
			for k, v := range o {
				_h[k] = fmt.Sprintf("%#v", v)
			}
		}
	}

	var res *req.Resp
	if t == "formdata" {
		_d := make(req.Param, 0)
		data, err := call.Argument(1).Export()
		if err == nil {
			o, flag := data.(map[string]interface{})
			if flag {
				for k, v := range o {
					_d[k] = fmt.Sprintf("%#v", v)
				}
			}
		}
		res, err = req.Post(url, _d, _h)
	} else if t == "json" {
		//data set
		data := call.Argument(1).String()
		_d := req.BodyJSON(data)
		res, err = req.Post(url, _d, _h)
	}
	if err != nil {
		return otto.Value{}
	}
	str, err := res.ToString()
	if err != nil {
		return otto.Value{}
	}
	result, _ := this.vm.ToValue(str)
	return result
}

////////////////////聊天相关API

//撤回群组消息
func (this *JsFn) fn_recallGroupMessage(call otto.FunctionCall) otto.Value {
	groupId, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	msgid, err := call.Argument(1).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	api.RecallGroupMessage(groupId, int32(msgid))
	return otto.Value{}
}

//发送群组消息,自定义字符串码
func (this *JsFn) fn_sendGroupMessageText(call otto.FunctionCall) otto.Value {
	groupId, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	text := call.Argument(1).String()
	m := api.SendGroupMessageText(groupId, text)
	result, err := this.vm.ToValue(m.MsgId.Id)
	if err != nil {
		return otto.Value{}
	}
	return result
}

//发送群组消息
func (this *JsFn) fn_sendGroupMessage(call otto.FunctionCall) otto.Value {
	groupId, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	obj, err := call.Argument(1).Export()
	if err != nil {
		return otto.Value{}
	}
	o, flag := obj.([]map[string]interface{})
	if !flag {
		return otto.Value{}
	}
	list := make([]api.IMsg, 0)
	//解析消息
	for _, item := range o {
		_type, h := item["type"]
		if !h {
			continue
		}
		switch _type {
		case "at":
			if value, has := item["target"]; has {
				if qq, _h := value.(int64); _h {
					list = append(list, api.NewAt(qq))
				}
			}
			break
		case "text":
			if value, has := item["content"]; has {
				list = append(list, api.NewText(fmt.Sprintf("%s", value)))
			}
			break
		case "image":
			if value, has := item["value"]; has {
				list = append(list, api.NewImage(groupId, fmt.Sprintf("%s", value)))
			}
			break
		}
	}
	m := api.SendGroupMessage(groupId, list)
	result, _ := this.vm.ToValue(m.MsgId.Id)
	return result
}
