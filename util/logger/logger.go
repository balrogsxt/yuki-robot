package logger

import (
	"fmt"
	"github.com/balrogsxt/xtbot-go/util"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type LoggerFormat struct {
}

func (s *LoggerFormat) Format(entry *log.Entry) ([]byte, error) {
	date := time.Now().Local().Format("2006-01-02 15:04:05")
	level := "NULL"
	color := "f"
	switch entry.Level {
	case log.InfoLevel:
		color = "a"
		level = "INFO "
		break
	case log.WarnLevel:
		color = "e"
		level = "WARN "
		break
	case log.ErrorLevel:
		color = "c"
		level = "ERROR"
		break
	case log.FatalLevel:
		color = "c"
		level = "FATAL"
		break
	case log.DebugLevel:
		color = "7"
		level = "DEBUG"
		break
	}
	text := entry.Message
	msg := util.PrintlnColor("§%s[%s] [%s]:§f %s\n", color, date, level, text)
	return []byte(msg), nil
}
func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(new(LoggerFormat))
}

//普通信息
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

//调试日志
func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

//错误日志
func Error(format string, args ...interface{}) {
	log.Errorf(fmt.Sprintf(format, args...))
}

//致命错误日志
func Fatal(format string, args ...interface{}) {
	log.Fatalf(fmt.Sprintf(format, args...))
}

//警告日志
func Warning(format string, args ...interface{}) {
	log.Warningf(fmt.Sprintf(format, args...))
}
