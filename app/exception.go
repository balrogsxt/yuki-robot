package app

import "fmt"

//抛出异常
func ThrowException(format string, v ...interface{}) {
	// todo 先抛出简单异常把
	panic(fmt.Sprintf(format, v...))
}
