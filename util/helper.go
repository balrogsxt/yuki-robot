package util

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
	"strings"
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

//输出mc颜色代码
func PrintlnColor(format string, args ...interface{}) string {
	text := fmt.Sprintf(format, args...)
	list := strings.Split(text, "§")

	result := ""
	for _, item := range list {
		if 0 >= len(item) {
			continue
		}
		c := item[0:1]
		val := item[1:]
		switch c {
		case "0":
			result += color.HEX("#000000").Sprintf(val)
			break
		case "1":
			result += color.HEX("#0000AA").Sprintf(val)
			break
		case "2":
			result += color.HEX("#00AA00").Sprintf(val)
			break
		case "3":
			result += color.HEX("#00AAAA").Sprintf(val)
			break
		case "4":
			result += color.HEX("#AA0000").Sprintf(val)
			break
		case "5":
			result += color.HEX("#AA00AA").Sprintf(val)
			break
		case "6":
			result += color.HEX("#FFAA00").Sprintf(val)
			break
		case "7":
			result += color.HEX("#AAAAAA").Sprintf(val)
			break
		case "8":
			result += color.HEX("#555555").Sprintf(val)
			break
		case "9":
			result += color.HEX("#5555FF").Sprintf(val)
			break

		case "a":
			result += color.HEX("#55FF55").Sprintf(val)
			break
		case "b":
			result += color.HEX("#55FFFF").Sprintf(val)
			break
		case "c":
			result += color.HEX("#FF5555").Sprintf(val)
			break
		case "d":
			result += color.HEX("#FF55FF").Sprintf(val)
			break
		case "e":
			result += color.HEX("#FFFF55").Sprintf(val)
			break
		case "f":
			result += color.HEX("#FFFFFF").Sprintf(val)
			break
		case "g":
			result += color.HEX("#DDD605").Sprintf(val)
			break
		default:
			result += color.HEX("#FFFFFF").Sprintf(val)
			break
		}
	}
	return result
}
