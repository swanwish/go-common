package logs

import (
	"errors"
	"fmt"
	"net/smtp"

	"github.com/swanwish/go-common/libsmtp"
)

type EmailConfig struct {
	SSL           bool
	Username      string
	Password      string
	Host          string
	Port          int64
	From          string
	To            []string
	SubjectPrefix string
}

const (
	DEFAULT_SUBJECT_PREFIX = "EMAIL LOG"
)

var (
	EMailLogger             LogWriter = EMailLogWriter{}
	ErrConfigurationMissing           = errors.New("The configuration is missing")
)

type EMailLogWriter struct {
	EmailConfig
}

func (w EMailLogWriter) Log(level int64, key string, v ...interface{}) {
	message := LogMessage(v...)
	levelText := GetLogLevelText(level)
	if err := w.sendMail(w.getMailSubject(levelText, key), message); err != nil {
		Errorf("Failed to send mail, the error is %v", err)
		Log(level, key, v...)
	}
}

func (w EMailLogWriter) Logf(level int64, key, format string, v ...interface{}) {
	message := LogMessagef(format, v...)
	levelText := GetLogLevelText(level)
	w.sendMail(w.getMailSubject(levelText, key), message)
}

func (w EMailLogWriter) Debug(v ...interface{}) {
	w.Log(LOG_LEVEL_DEBUG, EMPTY_KEY, v...)
}

func (w EMailLogWriter) Info(v ...interface{}) {
	w.Log(LOG_LEVEL_INFO, EMPTY_KEY, v...)
}

func (w EMailLogWriter) Warn(v ...interface{}) {
	w.Log(LOG_LEVEL_WARN, EMPTY_KEY, v...)
}

func (w EMailLogWriter) Error(v ...interface{}) {
	w.Log(LOG_LEVEL_ERROR, EMPTY_KEY, v...)
}

func (w EMailLogWriter) Fatal(v ...interface{}) {
	w.Log(LOG_LEVEL_FATAL, EMPTY_KEY, v...)
}

func (w EMailLogWriter) Debugf(formatString string, v ...interface{}) {
	w.Logf(LOG_LEVEL_DEBUG, EMPTY_KEY, formatString, v...)
}

func (w EMailLogWriter) Infof(formatString string, v ...interface{}) {
	w.Logf(LOG_LEVEL_INFO, EMPTY_KEY, formatString, v...)
}

func (w EMailLogWriter) Warnf(formatString string, v ...interface{}) {
	w.Logf(LOG_LEVEL_WARN, EMPTY_KEY, formatString, v...)
}

func (w EMailLogWriter) Errorf(formatString string, v ...interface{}) {
	w.Logf(LOG_LEVEL_ERROR, EMPTY_KEY, formatString, v...)
}

func (w EMailLogWriter) Fatalf(formatString string, v ...interface{}) {
	w.Logf(LOG_LEVEL_FATAL, EMPTY_KEY, formatString, v...)
}

func (l EMailLogWriter) sendMail(title, message string) error {
	if !l.Valid() {
		Info("The email log is not enabled")
		return ErrConfigurationMissing
	}
	auth := smtp.PlainAuth("", l.Username, l.Password, l.Host)
	err := libsmtp.SendMailWithAttachments(l.SSL, fmt.Sprintf("%s:%d", l.Host, l.Port), &auth, l.From, title, l.To, []byte(message), nil)
	if err != nil {
		Errorf("Failed to send mail, the error is %v", err)
	}
	return err
}

func (w EMailLogWriter) getMailSubject(level, key string) string {
	subject := ""
	if w.SubjectPrefix == "" {
		subject = DEFAULT_SUBJECT_PREFIX
	} else {
		subject = w.SubjectPrefix
	}

	if level != "" {
		level = fmt.Sprintf("[%s]", level)
	}

	subject = fmt.Sprintf("%s %s", subject, level)

	if key != "" {
		subject += " - " + key
	}

	return subject
}

func (c EmailConfig) Valid() bool {
	return c.Host != "" && c.Username != "" && c.Password != "" && len(c.To) > 0 && c.From != ""
}
