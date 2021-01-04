package cq

import "fmt"

//快速发送CQ码

func At(qq interface{}) string {
	return fmt.Sprintf("[CQ:at,qq=%v]", qq)
}

func Image(file interface{}) string {
	return fmt.Sprintf(`[CQ:image,file=%v]`, file)
}
func Face(id interface{}) string {
	return fmt.Sprintf(`[CQ:face,id=%v]`, id)
}
func Record(file interface{}) string {
	return fmt.Sprintf(`[CQ:record,file=%v]`, file)
}
