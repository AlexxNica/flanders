package log

import (
	"errors"
	"fmt"
	"io"
	"log/syslog"
	"runtime"
	"strings"
	"time"
)

var logger io.Writer
var logformat string = "%s: %s - %s\n" // Variables are level, time, msg
var loglevel syslog.Priority = syslog.LOG_WARNING

var syslogger *syslog.Writer
var useSyslog bool = true

var syslogLevelMap = map[syslog.Priority]string{
	syslog.LOG_EMERG:   "EMERGENCY",
	syslog.LOG_ALERT:   "ALERT",
	syslog.LOG_CRIT:    "CRITICAL",
	syslog.LOG_ERR:     "ERROR",
	syslog.LOG_WARNING: "WARNING",
	syslog.LOG_NOTICE:  "NOTICE",
	syslog.LOG_INFO:    "INFO",
	syslog.LOG_DEBUG:   "DEBUG",
}

var stringLevelMap = map[string]syslog.Priority{
	"EMERG":   syslog.LOG_EMERG,
	"ALERT":   syslog.LOG_ALERT,
	"CRIT":    syslog.LOG_CRIT,
	"ERR":     syslog.LOG_ERR,
	"WARNING": syslog.LOG_WARNING,
	"NOTICE":  syslog.LOG_NOTICE,
	"INFO":    syslog.LOG_INFO,
	"DEBUG":   syslog.LOG_DEBUG,
}

func init() {
	switch runtime.GOOS {
	case "linux", "darwin":
		err := setupDefaultLogger()
		if err != nil {
			fmt.Println("Cannot setup the default logger to syslog")
		}
	}
}

func log(level syslog.Priority, msg string) (int, error) {
	var err error
	if level <= loglevel {
		if logger != nil {
			return write(level, msg)
		}
		if useSyslog {
			switch level {
			case syslog.LOG_EMERG:
				err = syslogger.Emerg(msg)
			case syslog.LOG_ALERT:
				err = syslogger.Alert(msg)
			case syslog.LOG_CRIT:
				err = syslogger.Crit(msg)
			case syslog.LOG_ERR:
				err = syslogger.Err(msg)
			case syslog.LOG_WARNING:
				err = syslogger.Warning(msg)
			case syslog.LOG_NOTICE:
				err = syslogger.Notice(msg)
			case syslog.LOG_INFO:
				err = syslogger.Info(msg)
			case syslog.LOG_DEBUG:
				err = syslogger.Debug(msg)
			default:
				_, err = syslogger.Write([]byte(msg))
			}
			return 0, err
		}
	}

	return 0, nil
}

func Emerg(msg string) error {
	log(syslog.LOG_EMERG, msg)
	return nil
}
func Alert(msg string) error {
	log(syslog.LOG_ALERT, msg)
	return nil
}
func Crit(msg string) error {
	log(syslog.LOG_CRIT, msg)
	return nil
}
func Err(msg string) error {
	log(syslog.LOG_ERR, msg)
	return nil
}
func Warning(msg string) error {
	log(syslog.LOG_WARNING, msg)
	return nil
}
func Notice(msg string) error {
	log(syslog.LOG_NOTICE, msg)
	return nil
}
func Info(msg string) error {
	log(syslog.LOG_INFO, msg)
	return nil
}
func Debug(msg string) error {
	log(syslog.LOG_DEBUG, msg)
	return nil
}

func SetLogger(newLogger io.Writer) {
	logger = newLogger
}

// Accepts: EMERG|ALERT|CRIT|ERR|WARNING|NOICE|INFO|DEBUG
func SetLogLevel(level string) error {
	level = strings.ToUpper(level)
	syslogLevel, ok := stringLevelMap[level]
	if ok == false {
		return errors.New("Log level \"" + level + "\" not supported.")
	}
	loglevel = syslogLevel
	return nil
}

func SetSysLog(val bool) error {
	useSyslog = val
	return nil
}

func setupDefaultLogger() error {
	var err error
	syslogger, err = syslog.New(syslog.LOG_INFO, "SipRegistrar")
	return err
}

func write(p syslog.Priority, msg string) (int, error) {
	timestamp := time.Now().Format(time.RFC3339)
	return fmt.Fprintf(logger, logformat, syslogLevelMap[p], timestamp, msg)
}
