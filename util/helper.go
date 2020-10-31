package util

import (
	"bufio"
	"crypto/md5"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//解析JSON
func JsonDecode(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), &v)
}

//转为json字符串
func JsonEncode(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func Md5File(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
