package script

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/robot/api"
	"github.com/robertkrimen/otto"
	"io/ioutil"
)

type Javascript struct {
	vm *otto.Otto
}

func NewJs() *Javascript {
	js := new(Javascript)
	js.vm = otto.New()
	//初始化绑定支持的内置方法
	js.bindInternalFunctions()
	return js
}
func (this *Javascript) bindInternalFunctions() {

	//发送群消息
	this.vm.Set("sendGroupMessage", fn_sendGroupMessage)

}

func (this *Javascript) RunFile(file string) error {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	code := string(b)
	rt, err := this.vm.Run(code)
	if err != nil {
		return err
	}
	fmt.Println("返回值:", rt.String())

	return nil
}

func fn_sendGroupMessage(call otto.FunctionCall) otto.Value {
	groupId, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.Value{}
	}
	s := call.Argument(1).Object()

	obj, err := s.Value().Export()
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
	fmt.Printf("%#v \n", m)
	return otto.Value{}
}
