package script

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/balrogsxt/xtbot-go/util"
	"io"
	"io/ioutil"
	"os/exec"
)
type vars map[interface{}]interface{}

type NodeJs struct {
	bindVars vars
}

func NewNodeJs() *NodeJs {
	env := new(NodeJs)
	env.bindVars = make(vars,0)
	return env
}


func (this *NodeJs) BindVar(name interface{},value interface{}) {
	this.bindVars[name] = value
}
func (this *NodeJs) SetVars(vars map[interface{}]interface{}) {
	this.bindVars = vars
}
func (this *NodeJs) RunFile(file string) error {
	byte,err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	vars := util.JsonEncode(this.bindVars)
	cmd := exec.Command("node", "./plugins/env/node.js",vars,string(byte))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start();err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	if err := cmd.Wait();err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
