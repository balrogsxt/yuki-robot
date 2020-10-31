package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type LoggerFormat struct {
}

func (s *LoggerFormat) Format(entry *log.Entry) ([]byte, error) {
	date := time.Now().Local().Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	text := entry.Message
	msg := fmt.Sprintf("[%s] [%s]: %s\n", date, level, text)
	return []byte(msg), nil
}
func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(new(LoggerFormat))
}

//普通信息
func Info(format string, args ...interface{}) {
	log.Infof(fmt.Sprintf(format, args...))
}

//调试日志
func Debug(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
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
