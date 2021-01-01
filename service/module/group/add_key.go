package group

import "fmt"

//回复关键词添加
type AddKey struct {

}

func (AddKey) Command() string {
	return "^add"
}
func (AddKey) Call() bool {
	fmt.Println("触发命令")

	return true
}

