package logs

import (
	"fmt"
	"log"
	"strings"
)

const (
	LOG_LEVEL_DEBUG int64 = 1 + iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
)

var MaxLogLevel = LOG_LEVEL_DEBUG

var Writer LogWriter = DefaultLogWriter{}
var logLevelText = make(map[int64]string, 0)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logLevelText[LOG_LEVEL_DEBUG] = "debug"
	logLevelText[LOG_LEVEL_INFO] = "info"
	logLevelText[LOG_LEVEL_WARN] = "warn"
	logLevelText[LOG_LEVEL_ERROR] = "error"
	logLevelText[LOG_LEVEL_FATAL] = "fatal"
}

func RegisterLogLevelText(level int64, text string) {
	logLevelText[level] = text
}

func SetLogLevel(logLevel string) {
	Debugf("The loglevel is %s", logLevel)
	switch strings.ToLower(logLevel) {
	case "debug":
		MaxLogLevel = LOG_LEVEL_DEBUG
	case "info":
		MaxLogLevel = LOG_LEVEL_INFO
	case "error":
		MaxLogLevel = LOG_LEVEL_ERROR
	}
}

func GetLogLevelText(logLevel int64) string {
	if text, ok := logLevelText[logLevel]; ok {
		return text
	}
	return ""
}

func Debug(v ...interface{}) {
	if Writer != nil {
		Writer.Debug(v...)
	}
}

func Info(v ...interface{}) {
	if Writer != nil {
		Writer.Info(v...)
	}
}

func Warn(v ...interface{}) {
	if Writer != nil {
		Writer.Warn(v...)
	}
}

func Error(v ...interface{}) {
	if Writer != nil {
		Writer.Error(v...)
	}
}

func Fatal(v ...interface{}) {
	if Writer != nil {
		Writer.Fatal(v...)
	}
}

func Debugf(formatString string, v ...interface{}) {
	if Writer != nil {
		Writer.Debugf(formatString, v...)
	}
}

func Infof(formatString string, v ...interface{}) {
	if Writer != nil {
		Writer.Infof(formatString, v...)
	}
}

func Warnf(formatString string, v ...interface{}) {
	if Writer != nil {
		Writer.Warnf(formatString, v...)
	}
}

func Errorf(formatString string, v ...interface{}) {
	if Writer != nil {
		Writer.Errorf(formatString, v...)
	}
}

func Fatalf(formatString string, v ...interface{}) {
	if Writer != nil {
		Writer.Fatalf(formatString, v...)
	}
}

func LogMessage(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func LogMessagef(formatString string, v ...interface{}) string {
	return fmt.Sprintf(formatString, v...)
}
