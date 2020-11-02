package util

import "regexp"

//判断是否是ipv4格式
func IsIpv4(str string) bool {
	flag, err := regexp.Match("^\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}$", []byte(str))
	if err == nil && flag {
		return true
	}
	return false
}
