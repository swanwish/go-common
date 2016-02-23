package logs

import (
	"fmt"
	"log"
)

const (
	EMPTY_KEY        = ""
	LOG_CALLER_LEVEL = 5
)

type LogWriter interface {
	Log(level int64, key string, v ...interface{})
	Logf(level int64, key, format string, v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Debugf(formatString string, v ...interface{})
	Infof(formatString string, v ...interface{})
	Errorf(formatString string, v ...interface{})
	Fatalf(formatString string, v ...interface{})
}

type DefaultLogWriter struct {
}

func (w DefaultLogWriter) Log(level int64, key string, v ...interface{}) {
	Log(level, key, v...)
}

func (w DefaultLogWriter) Logf(level int64, key, format string, v ...interface{}) {
	Logf(level, key, format, v...)
}

func (w DefaultLogWriter) Debug(v ...interface{}) {
	Log(LOG_LEVEL_DEBUG, EMPTY_KEY, v...)
}

func (w DefaultLogWriter) Info(v ...interface{}) {
	Log(LOG_LEVEL_INFO, EMPTY_KEY, v...)
}

func (w DefaultLogWriter) Error(v ...interface{}) {
	Log(LOG_LEVEL_ERROR, EMPTY_KEY, v...)
}

func (w DefaultLogWriter) Fatal(v ...interface{}) {
	Log(LOG_LEVEL_FATAL, EMPTY_KEY, v...)
}

func (w DefaultLogWriter) Debugf(formatString string, v ...interface{}) {
	Logf(LOG_LEVEL_DEBUG, EMPTY_KEY, formatString, v...)
}

func (w DefaultLogWriter) Infof(formatString string, v ...interface{}) {
	Logf(LOG_LEVEL_INFO, EMPTY_KEY, formatString, v...)
}

func (w DefaultLogWriter) Errorf(formatString string, v ...interface{}) {
	Logf(LOG_LEVEL_ERROR, EMPTY_KEY, formatString, v...)
}

func (w DefaultLogWriter) Fatalf(formatString string, v ...interface{}) {
	Logf(LOG_LEVEL_FATAL, EMPTY_KEY, formatString, v...)
}

func Log(logLevel int64, key string, v ...interface{}) {
	if logLevel < MaxLogLevel {
		return
	}
	levelText := GetLogLevelText(logLevel)
	logMessage := LogMessage(v...)
	if key == EMPTY_KEY {
		log.Output(LOG_CALLER_LEVEL, fmt.Sprintln(levelText, "-", logMessage))
	} else {
		log.Output(LOG_CALLER_LEVEL, fmt.Sprintln(levelText, "-", key, "-", logMessage))
	}
}

func Logf(logLevel int64, key, formatString string, v ...interface{}) {
	if logLevel < MaxLogLevel {
		return
	}
	logMessage := LogMessagef(formatString, v...)
	levelText := GetLogLevelText(logLevel)
	if key == EMPTY_KEY {
		log.Output(LOG_CALLER_LEVEL, fmt.Sprintln(levelText, "-", logMessage))
	} else {
		log.Output(LOG_CALLER_LEVEL, fmt.Sprintln(levelText, " - ", key, " - ", logMessage))
	}
}
